package repository

import (
	"context"
	"time"

	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateByID(ctx context.Context, id int64, username, email, password *string) (*entity.User, error)
	DeleteByID(ctx context.Context, id int64) (*entity.User, error)
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(ctx context.Context, user *entity.User) (int64, error) {
	var id int64
	err := r.db.QueryRowxContext(ctx, sqlInsertUser,
		user.Username, user.Email, user.Password, time.Now(),
	).Scan(&id)
	return id, err
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, sqlFindUserByID, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, sqlFindUserByEmail, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, sqlFindUserByUsername, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdateByID(ctx context.Context, id int64, username, email, password *string) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, sqlUpdateUserByID, username, email, password, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, sqlDeleteUserByID, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.SelectContext(ctx, &users, sqlListUsers, limit, offset)
	return users, err
}
