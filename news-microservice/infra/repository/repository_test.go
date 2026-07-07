package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	infraerrors "news-service/infra/errors"
	"news-service/infra/models/domain"
)

func fixture(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	db, m, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	g, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	require.NoError(t, err)
	return New(g), m
}

func TestCreateAndGet(t *testing.T) {
	r, m := fixture(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "news"`).WithArgs("t", "c", "h", sqlmock.AnyArg(), "u").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("n"))
	m.ExpectCommit()
	got, err := r.Create(context.Background(), &domain.News{Title: "t", Content: "c", HouseID: "h", CreatedBy: "u"})
	require.NoError(t, err)
	assert.Equal(t, "n", got.ID)
	assert.False(t, got.CreatedAt.IsZero())

	r, m = fixture(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "news"`).WillReturnError(errors.New("db"))
	m.ExpectRollback()
	_, err = r.Create(context.Background(), &domain.News{})
	assert.EqualError(t, err, "create news: db")

	cols := []string{"id", "title", "content", "house_id", "created_at", "created_by"}
	now := time.Now()
	r, m = fixture(t)
	m.ExpectQuery(`SELECT \* FROM "news" WHERE id = \$1`).WithArgs("n", 1).WillReturnRows(sqlmock.NewRows(cols).AddRow("n", "t", "c", "h", now, "u"))
	got, err = r.GetByID(context.Background(), "n")
	require.NoError(t, err)
	assert.Equal(t, "t", got.Title)
	r, m = fixture(t)
	m.ExpectQuery(`SELECT \* FROM "news"`).WillReturnRows(sqlmock.NewRows(cols))
	_, err = r.GetByID(context.Background(), "x")
	require.ErrorIs(t, err, infraerrors.ErrNewsNotFound)
	r, m = fixture(t)
	m.ExpectQuery(`SELECT \* FROM "news"`).WillReturnError(errors.New("db"))
	_, err = r.GetByID(context.Background(), "x")
	assert.EqualError(t, err, "get news: db")
}

func TestList(t *testing.T) {
	now := time.Now()
	cols := []string{"id", "title", "content", "house_id", "created_at", "created_by"}
	req := &domain.GetNewsRequest{HouseID: "h", DateFrom: now.Add(-time.Hour), DateTo: now, Limit: 10, Offset: 2}
	r, m := fixture(t)
	m.ExpectQuery(`SELECT count\(\*\) FROM "news" WHERE house_id`).WithArgs("h", req.DateFrom, req.DateTo).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	m.ExpectQuery(`SELECT \* FROM "news" WHERE house_id`).WithArgs("h", req.DateFrom, req.DateTo, 10, 2).WillReturnRows(sqlmock.NewRows(cols).AddRow("n", "t", "c", "h", now, "u"))
	items, total, err := r.List(context.Background(), req)
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(1), total)

	r, m = fixture(t)
	m.ExpectQuery(`SELECT count`).WillReturnError(errors.New("db"))
	_, _, err = r.List(context.Background(), &domain.GetNewsRequest{})
	assert.EqualError(t, err, "count news: db")
	r, m = fixture(t)
	m.ExpectQuery(`SELECT count`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	m.ExpectQuery(`SELECT \*`).WithArgs(20).WillReturnError(errors.New("db"))
	_, _, err = r.List(context.Background(), &domain.GetNewsRequest{Limit: 200, Offset: -1})
	assert.EqualError(t, err, "list news: db")
}

func TestHelpers(t *testing.T) {
	now := time.Now()
	row := &dbNews{ID: "n", Title: "t", Content: "c", HouseID: "h", CreatedAt: now, CreatedBy: "u"}
	assert.Equal(t, "news", row.TableName())
	assert.Equal(t, &domain.News{ID: "n", Title: "t", Content: "c", HouseID: "h", CreatedAt: now, CreatedBy: "u"}, newsToDomain(row))
}
