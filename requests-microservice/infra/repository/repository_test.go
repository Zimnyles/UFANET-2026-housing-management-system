package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	infraerrors "requests-service/infra/errors"
	"requests-service/infra/models/domain"
	"testing"
	"time"
)

func fix(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	db, m, e := sqlmock.New()
	require.NoError(t, e)
	t.Cleanup(func() { _ = db.Close() })
	g, e := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	require.NoError(t, e)
	return New(g), m
}

var reqCols = []string{"id", "title", "description", "type", "status", "user_id", "created_at", "updated_at"}

func TestRepository(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	r, m := fix(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "maintenance_requests"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("r"))
	m.ExpectCommit()
	got, e := r.Create(ctx, &domain.MaintenanceRequest{Title: "t", Description: "d", Type: "plumber", UserID: "u"})
	require.NoError(t, e)
	assert.Equal(t, "r", got.ID)
	assert.Equal(t, "open", got.Status)
	r, m = fix(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT`).WillReturnError(errors.New("db"))
	m.ExpectRollback()
	_, e = r.Create(ctx, &domain.MaintenanceRequest{})
	assert.EqualError(t, e, "create request: db")
	r, m = fix(t)
	m.ExpectQuery(`SELECT \* FROM "maintenance_requests" WHERE id`).WithArgs("r", 1).WillReturnRows(sqlmock.NewRows(reqCols).AddRow("r", "t", "d", "plumber", "open", "u", now, now))
	got, e = r.GetByID(ctx, "r")
	require.NoError(t, e)
	assert.Equal(t, "t", got.Title)
	r, m = fix(t)
	m.ExpectQuery(`SELECT \*`).WillReturnRows(sqlmock.NewRows(reqCols))
	_, e = r.GetByID(ctx, "x")
	require.ErrorIs(t, e, infraerrors.ErrRequestNotFound)
	r, m = fix(t)
	m.ExpectQuery(`SELECT \*`).WillReturnError(errors.New("db"))
	_, e = r.GetByID(ctx, "x")
	assert.EqualError(t, e, "get request: db")
}
func TestListErrors(t *testing.T) {
	r, m := fix(t)
	m.ExpectQuery(`SELECT count`).WillReturnError(errors.New("db"))
	_, _, e := r.List(context.Background(), &domain.GetRequestsRequest{UserID: "u", Status: "open", Type: "plumber"})
	assert.EqualError(t, e, "count requests: db")
	r, m = fix(t)
	m.ExpectQuery(`SELECT count`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	m.ExpectQuery(`SELECT \*`).WithArgs(20).WillReturnError(errors.New("db"))
	_, _, e = r.List(context.Background(), &domain.GetRequestsRequest{Limit: 200, Offset: -1})
	assert.EqualError(t, e, "list requests: db")
}
func TestUpdateAndCommentsErrors(t *testing.T) {
	r, m := fix(t)
	m.ExpectQuery(`SELECT \*`).WillReturnRows(sqlmock.NewRows(reqCols))
	_, e := r.UpdateStatus(context.Background(), "x", "done")
	require.ErrorIs(t, e, infraerrors.ErrRequestNotFound)
	r, m = fix(t)
	m.ExpectQuery(`SELECT \*`).WillReturnRows(sqlmock.NewRows(reqCols))
	_, e = r.AddComment(context.Background(), &domain.Comment{RequestID: "x"})
	require.ErrorIs(t, e, infraerrors.ErrRequestNotFound)
	r, m = fix(t)
	m.ExpectQuery(`SELECT \* FROM "request_comments"`).WillReturnError(errors.New("db"))
	_, e = r.ListComments(context.Background(), "r")
	assert.EqualError(t, e, "list comments: db")
}
func TestConverters(t *testing.T) {
	now := time.Now()
	assert.Equal(t, "maintenance_requests", (dbMaintenanceRequest{}).TableName())
	assert.Equal(t, "request_comments", (dbComment{}).TableName())
	assert.Equal(t, "r", requestToDomain(&dbMaintenanceRequest{ID: "r", CreatedAt: now}).ID)
	assert.Equal(t, "c", commentToDomain(&dbComment{ID: "c", CreatedAt: now}).ID)
}
