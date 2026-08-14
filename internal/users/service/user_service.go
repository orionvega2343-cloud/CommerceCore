package service

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/internal/users/dto"
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var _ domain.UserService = (*UserServiceImpl)(nil)

type UserServiceImpl struct {
	repo   domain.UserRepo
	secret string
}

func NewUserServiceImpl(repo domain.UserRepo, secret string) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, secret: secret}
}

// toDomainUser - собирает доменную модель пользователя из DTO запроса
func toDomainUser(req dto.UserRequest) domain.User {
	return domain.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
}

// toUserResponse - собирает DTO ответа из доменной модели пользователя
func toUserResponse(u *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:        u.Id,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}

// Register - функция для хэширования пользовательского пароля,
// при помощи криптографии(bcrypt), превращает строку в набор символов,
// подставляет значение к доменной модели и сохраняет в БД
func (o *UserServiceImpl) Register(ctx context.Context, req dto.UserRequest) (*dto.UserResponse, error) {
	user := toDomainUser(req)

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, domain.FailedHashingPassword
	}
	user.Password = string(hashed)

	u, err := o.repo.CreateUser(ctx, &user)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, domain.FailedCreatedUser
	}
	return toUserResponse(u), nil
}

// Login - авторизует пользователя по Email,
// проверяет, что пароль успешно хэширован
// и пропускает пользователя
func (o *UserServiceImpl) Login(ctx context.Context, req dto.UserRequest) (string, error) {
	user, err := o.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("failed to get user by email", "error", err)
		return "", err
	}
	//Проверяем, что хэширование применилось и пароль существует
	comparePass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if comparePass != nil {
		slog.Error("failed to compare password")
		return "", domain.InvalidPassword
	}

	var c domain.Claims
	c.Email = req.Email
	c.UserId = user.Id
	//Устанавливаем срок действия JWT токена
	//TODO: пробросить TTL из конфига
	c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString([]byte(o.secret))
	if err != nil {
		slog.Error("failed to sign token", "error", err)
		return "", err
	}
	return tokenString, nil

}

func (o *UserServiceImpl) GetById(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := o.repo.GetUserByID(ctx, id)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		return nil, domain.FailedToGetUser
	}
	return toUserResponse(user), nil
}

func (o *UserServiceImpl) UpdateUser(ctx context.Context, id string, req dto.UserRequest) error {
	_, err := o.repo.GetUserByID(ctx, id)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		return domain.FailedToGetUser
	}

	user := toDomainUser(req)
	user.Id = id

	err = o.repo.UpdateUser(ctx, user)
	if err != nil {
		slog.Error("failed to update user", "error", err)
		return domain.FailedToUpdateUser
	}
	return nil
}

func (o *UserServiceImpl) UpdateRole(ctx context.Context, role, id string) error {
	user, err := o.repo.GetUserByID(ctx, id)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		return domain.FailedToGetUser
	}

	if role != "admin" && role != "customer" {
		slog.Info("role is not admin", "role", user.Role)
		return domain.InvalidUserRole
	}

	err = o.repo.UpdateUserRole(ctx, id, role)
	if err != nil {
		slog.Error("failed to update user role", "error", err)
		return domain.FailedToUpdateUser
	}
	return nil
}
