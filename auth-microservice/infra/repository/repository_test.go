package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	infraerrors "auth-service/infra/errors"
	"auth-service/infra/models/domain"
)

func newRepo(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{TranslateError: true})
	require.NoError(t, err)
	return New(gdb), mock
}

func TestRepositoryIsActiveAdminCode(t *testing.T) {
	r, mock := newRepo(t)
	mock.ExpectQuery(`SELECT count\(\*\) FROM "admin_codes" WHERE code = \$1 AND status = \$2`).WithArgs("code", "active").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	active, err := r.IsActiveAdminCode(context.Background(), "code")
	require.NoError(t, err)
	assert.True(t, active)
	require.NoError(t, mock.ExpectationsWereMet())

	r, mock = newRepo(t)
	mock.ExpectQuery(`SELECT count`).WillReturnError(errors.New("db"))
	active, err = r.IsActiveAdminCode(context.Background(), "code")
	assert.False(t, active)
	assert.EqualError(t, err, "check admin code: db")
}

func TestRepositoryCreateUser(t *testing.T) {
	user := &domain.User{ID: "u", Email: "a@b.ru", PasswordHash: "h", Role: "user"}
	r, mock := newRepo(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).WithArgs(user.Email, "h", "user", sqlmock.AnyArg(), "u").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("u"))
	mock.ExpectCommit()
	got, err := r.CreateUser(context.Background(), user)
	require.NoError(t, err)
	assert.Same(t, user, got)
	require.NoError(t, mock.ExpectationsWereMet())

	r, mock = newRepo(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()
	_, err = r.CreateUser(context.Background(), &domain.User{})
	require.ErrorIs(t, err, infraerrors.ErrEmailAlreadyExists)

	r, mock = newRepo(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).WillReturnError(errors.New("db"))
	mock.ExpectRollback()
	_, err = r.CreateUser(context.Background(), &domain.User{})
	assert.EqualError(t, err, "create user: db")
}

func TestRepositoryGetUserByEmail(t *testing.T) {
	cols := []string{"id", "email", "password", "role", "created_at"}
	r, mock := newRepo(t)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).WithArgs("a@b.ru", 1).WillReturnRows(sqlmock.NewRows(cols).AddRow("u", "a@b.ru", "h", "admin", nil))
	got, err := r.GetUserByEmail(context.Background(), "a@b.ru")
	require.NoError(t, err)
	assert.Equal(t, "u", got.ID)
	assert.Equal(t, "admin", got.Role)

	r, mock = newRepo(t)
	mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(sqlmock.NewRows(cols))
	_, err = r.GetUserByEmail(context.Background(), "none")
	require.ErrorIs(t, err, infraerrors.ErrUserNotFound)

	r, mock = newRepo(t)
	mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(errors.New("db"))
	_, err = r.GetUserByEmail(context.Background(), "x")
	assert.EqualError(t, err, "get user by email: db")
}

type pgError string

func (e pgError) Error() string    { return string(e) }
func (e pgError) SQLState() string { return string(e) }

func TestHelpers(t *testing.T) {
	assert.Equal(t, "admin_codes", (AdminCode{}).TableName())
	assert.Equal(t, "users", (domain.User{}).TableName())
	assert.True(t, containsCode(pgError("23505"), "23505"))
	assert.False(t, containsCode(errors.New("x"), "23505"))
}
