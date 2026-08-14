package repository

import (
	"CommerceCore/internal/users/domain"
	"CommerceCore/pkg/querier"
	"context"
	"log/slog"
)

var _ domain.UserRepo = (*UserRepoImpl)(nil)

type UserRepoImpl struct {
	q querier.Querier
}

func NewUserRepo(q querier.Querier) *UserRepoImpl {
	return &UserRepoImpl{q: q}
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, m *domain.User) (*domain.User, error) {
	err := r.q.GetContext(ctx, &m, `INSERT INTO users(email, password, role) VALUES($1, $2, $3) RETURNING id`, m.Email, m.Password, m.Role)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *UserRepoImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var m *domain.User
	err := r.q.GetContext(ctx, &m, `SELECT id, email, password, role, created_at FROM users WHERE id = $1`, id)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, err
	}
	return m, nil
}

func (r *UserRepoImpl) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var m domain.User
	err := r.q.GetContext(ctx, &m, `SELECT id, email, password, role, created_at FROM users WHERE email = $1`, email)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return domain.User{}, err
	}
	return m, nil
}

func (r *UserRepoImpl) UpdateUser(ctx context.Context, m domain.User) error {
	_, err := r.q.ExecContext(ctx, `UPDATE users SET email = $1, password = $2 WHERE id = $3`, m.Email, m.Password, m.Id)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	return nil
}

func (r *UserRepoImpl) UpdateUserRole(ctx context.Context, role, id string) error {
	_, err := r.q.ExecContext(ctx, `UPDATE users SET role = $1 WHERE id = $2`, role, id)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	return nil
}
