package main

import (
	cRepo "CommerceCore/internal/cart/repository"
	cSvc "CommerceCore/internal/cart/service"

	cHandler "CommerceCore/internal/cart/handler"

	pRepo "CommerceCore/internal/catalog/repository"
	pSvc "CommerceCore/internal/catalog/service"

	pHandler "CommerceCore/internal/catalog/handler"

	"CommerceCore/internal/infra/kafka"

	oRepo "CommerceCore/internal/order/repository"
	oSvc "CommerceCore/internal/order/service"

	oHandler "CommerceCore/internal/order/handler"

	payRepo "CommerceCore/internal/payment/repository"
	paySvc "CommerceCore/internal/payment/service"

	payHandler "CommerceCore/internal/payment/handler"

	uRepo "CommerceCore/internal/users/repository"
	uSvc "CommerceCore/internal/users/service"

	uHandler "CommerceCore/internal/users/handler"

	"CommerceCore/pkg/config"
	"CommerceCore/pkg/logger"
	"CommerceCore/pkg/middlewares"
	"CommerceCore/pkg/postgres"
	"CommerceCore/pkg/transaction"
	"context"
	"log"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	segkafka "github.com/segmentio/kafka-go"
)

func main() {
	slog.SetDefault(logger.NewLogger())
	cfg := config.MustLoad()

	db, err := postgres.Connect(*cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			slog.Error("failed to close database connection", "error", cerr)
		}
	}()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})
	defer func() {
		if cerr := rdb.Close(); cerr != nil {
			slog.Error("failed to close redis connection", "error", cerr)
		}
	}()

	tx := *transaction.NewTransactor(db)

	//Kafka producers один Writer = один топик
	orderWriter := &segkafka.Writer{
		Addr:                   segkafka.TCP(cfg.Kafka.Brokers...),
		Topic:                  cfg.Kafka.Topics.OrderCreated,
		AllowAutoTopicCreation: true,
	}
	defer func() {
		if cerr := orderWriter.Close(); cerr != nil {
			slog.Error("failed to close order kafka writer", "error", cerr)
		}
	}()
	orderEvents := kafka.NewProducer(orderWriter)

	paymentWriter := &segkafka.Writer{
		Addr:                   segkafka.TCP(cfg.Kafka.Brokers...),
		Topic:                  cfg.Kafka.Topics.PaymentSucceeded,
		AllowAutoTopicCreation: true,
	}
	defer func() {
		if cerr := paymentWriter.Close(); cerr != nil {
			slog.Error("failed to close payment kafka writer", "error", cerr)
		}
	}()
	paymentEvents := kafka.NewProducer(paymentWriter)

	//  Репозитории
	userRepo := uRepo.NewUserRepo(db)
	productRepo := pRepo.NewProductRepoImpl(db)
	productCache := pRepo.NewProductCache(rdb)
	cartItemRepo := cRepo.NewCartItemRepo(db)
	cartRepo := cRepo.NewCartRepo(db, cartItemRepo)
	cartRedisRepo := cRepo.NewRedisRepo(rdb)
	orderRepo := oRepo.NewOrderRepo(db)
	paymentRepo := payRepo.NewPaymentRepo(db)
	idempotencyRepo := payRepo.NewIdempotencyRepo(rdb)

	//Сервисы
	productSvc := pSvc.NewProductService(productRepo, productCache)
	cartSvc := cSvc.NewCartService(cartRepo, cartRedisRepo, cartItemRepo, productSvc)
	userSvc := uSvc.NewUserServiceImpl(userRepo, cartSvc, cfg.Jwt.Secret)
	orderSvc := oSvc.NewOrderService(orderRepo, cartRepo, productRepo, tx, orderEvents)
	paymentSvcPlain := paySvc.NewPaymentService(paymentRepo, orderRepo, tx, paymentEvents)
	paymentSvc := paySvc.NewIdempotencyService(paymentSvcPlain, idempotencyRepo, paymentRepo)

	//ендлеры
	userHandler := uHandler.NewUserHandlerImpl(userSvc)
	productHandler := pHandler.NewProductsHandlerImpl(productSvc)
	cartHandler := cHandler.NewCartHandlerImpl(*cartSvc)
	orderHandler := oHandler.NewOrderHandler(orderSvc)
	paymentHandler := payHandler.NewPaymentHandler(paymentSvc)

	//Kafka: consumerts, по одному Reader на топик, каждый в своей горутине
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startConsumer := func(topic string) {
		reader := segkafka.NewReader(segkafka.ReaderConfig{
			Brokers: cfg.Kafka.Brokers,
			Topic:   topic,
			GroupID: cfg.Kafka.GroupId,
		})
		consumer := kafka.NewConsumer(reader, &kafka.EventLogger{})
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("consumer goroutine panicked", "topic", topic, "recover", rec)
				}
			}()
			if err := consumer.Run(ctx); err != nil {
				slog.Error("consumer stopped", "topic", topic, "error", err)
			}
		}()
	}
	startConsumer(cfg.Kafka.Topics.OrderCreated)
	startConsumer(cfg.Kafka.Topics.PaymentSucceeded)

	// HTTP
	r := gin.Default()
	r.Use(middlewares.Recovery())
	r.Use(middlewares.RequestID())
	r.Use(middlewares.Logger())

	r.GET("/connection", postgres.HealthCheck(db))

	//публичные - без JWT
	r.POST("/api/v1/auth/register", userHandler.Register)
	r.POST("/api/v1/auth/login", userHandler.Login)
	r.GET("/api/v1/products", productHandler.GetAllProducts)
	r.GET("/api/v1/products/:id", productHandler.GetProductById)

	//гостевая корзина - по guestId, без JWT
	r.GET("/api/v1/carts/guest", cartHandler.GetGuestCart)
	r.POST("/api/v1/carts/guest/items", cartHandler.AddItemToCart)
	r.PATCH("/api/v1/carts/guest/items", cartHandler.UpdateItemIntoCart)
	r.DELETE("/api/v1/carts/guest/items", cartHandler.DeleteItemIntoCart)

	//защищённые - требуют валидный JWT (Auth кладёт user_id/role в контекст)
	auth := r.Group("/api/v1")
	auth.Use(middlewares.Auth())
	{
		auth.GET("/users/:id", userHandler.GetUserById)
		auth.PATCH("/users/:id", userHandler.UpdateUser)
		auth.PATCH("/users/:id/role", userHandler.UpdateRole)

		auth.POST("/products", productHandler.CreateProduct)
		auth.PATCH("/products/:id", productHandler.UpdateProduct)
		auth.DELETE("/products/:id", productHandler.DeleteProduct)

		auth.GET("/carts", cartHandler.CreateOrGet)
		auth.POST("/carts/items", cartHandler.Create)
		auth.PATCH("/carts/items/:id", cartHandler.Update)
		auth.DELETE("/carts/items/:id", cartHandler.Delete)
		// TODO: /carts/merge - CartHandlerImpl.MergeGuestCart не реализован
		// (сервисный CartServiceImpl.MergeGuestCart есть, HTTP-обёртки нет)
		auth.POST("/carts/checkout", orderHandler.Checkout)

		auth.GET("/orders", orderHandler.ListOrders)
		auth.GET("/orders/:orderId", orderHandler.GetOrder)
		auth.PATCH("/orders/:orderId/status", orderHandler.TransitionStatus)

		auth.POST("/orders/:orderId/payments", paymentHandler.Create)

		auth.GET("/payments/total-amount", paymentHandler.TotalAmountPayments)
		auth.GET("/payments/:id", paymentHandler.GetPaymentById)
		auth.GET("/payments", paymentHandler.ListPayments)
	}

	if err := r.Run(":" + strconv.Itoa(cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
