package domain

import "context"

// IdempotencyRepo - защита Payment.Create от дублирующих запросов через Redis.
// Реализацию (SETNX/SET/DEL на реальном Redis-клиенте) кладём в payment/repository,
// зеркалит cart.RedisRepo - Redis-доступ всегда идёт через отдельный интерфейс,
// не сырым клиентом внутри сервиса.
type IdempotencyRepo interface {
	// Claim - атомарно занимает ключ (аналог SETNX), вызывается ДО Payment.Create.
	// claimed=true  - ключа раньше не было, можно продолжать обработку.
	// claimed=false - ключ уже занят; cachedPaymentId - id платежа, сохранённый
	//                 предыдущим (успешным) вызовом с этим же ключом - Create
	//                 в этом случае вызывать не нужно, отдаём его результат.
	Claim(ctx context.Context, key string) (cachedPaymentId int, claimed bool, err error)

	// Store - вызывается ПОСЛЕ успешного Create: сохраняет id платежа под ключом,
	// чтобы повторный запрос с тем же ключом получил закешированный результат.
	Store(ctx context.Context, key string, paymentId int) error

	// Release - вызывается, если Create ПОСЛЕ успешного Claim завершился ошибкой:
	// освобождает ключ, чтобы легитимный повтор не был заблокирован до истечения TTL.
	Release(ctx context.Context, key string) error
}
