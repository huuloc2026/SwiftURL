package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/huuloc2026/SwiftURL/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMockRepo(t *testing.T) (*userRepo, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return &userRepo{db: sqlxDB}, mock, func() { db.Close() }
}

func TestCreate(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	user := &entity.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUser)).
		WithArgs(user.Username, user.Email, user.Password, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestFindByID(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	user := entity.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlFindUserByID)).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password))

	got, err := repo.FindByID(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
}

func TestFindByEmail(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	user := entity.User{
		ID:       2,
		Username: "testuser2",
		Email:    "test2@example.com",
		Password: "hashedpassword2",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlFindUserByEmail)).
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password))

	got, err := repo.FindByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
}

func TestFindByUsername(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	user := entity.User{
		ID:       3,
		Username: "testuser3",
		Email:    "test3@example.com",
		Password: "hashedpassword3",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlFindUserByUsername)).
		WithArgs(user.Username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password))

	got, err := repo.FindByUsername(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
}

func TestUpdateByID(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	id := int64(4)
	username := "updateduser"
	email := "updated@example.com"
	password := "updatedpassword"

	user := entity.User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlUpdateUserByID)).
		WithArgs(&username, &email, &password, id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password))

	got, err := repo.UpdateByID(context.Background(), id, &username, &email, &password)
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
}

func TestDeleteByID(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	id := int64(5)
	user := entity.User{
		ID:       id,
		Username: "deluser",
		Email:    "del@example.com",
		Password: "delpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlDeleteUserByID)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password))

	got, err := repo.DeleteByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, &user, got)
}

func TestList(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	users := []*entity.User{
		{ID: 1, Username: "user1", Email: "user1@example.com", Password: "pass1"},
		{ID: 2, Username: "user2", Email: "user2@example.com", Password: "pass2"},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password"})
	for _, u := range users {
		rows.AddRow(u.ID, u.Username, u.Email, u.Password)
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlListUsers)).
		WithArgs(2, 0).
		WillReturnRows(rows)

	got, err := repo.List(context.Background(), 2, 0)
	assert.NoError(t, err)
	assert.Equal(t, users, got)
}

func TestFindByID_NotFound(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	mock.ExpectQuery(regexp.QuoteMeta(sqlFindUserByID)).
		WithArgs(int64(999)).
		WillReturnError(errors.New("not found"))

	got, err := repo.FindByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestCreate_Error(t *testing.T) {
	repo, mock, close := newMockRepo(t)
	defer close()

	user := &entity.User{
		Username: "failuser",
		Email:    "fail@example.com",
		Password: "failpassword",
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlInsertUser)).
		WithArgs(user.Username, user.Email, user.Password, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	id, err := repo.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
}
