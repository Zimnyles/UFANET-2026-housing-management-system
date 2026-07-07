package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	infraerrors "profile-service/infra/errors"
	"profile-service/infra/models/domain"
	"testing"
	"time"
)

func repoFixture(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	db, m, e := sqlmock.New()
	require.NoError(t, e)
	t.Cleanup(func() { _ = db.Close() })
	g, e := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	require.NoError(t, e)
	return New(g), m
}

var profileCols = []string{"user_id", "full_name", "phone", "apartment", "house_id", "uk_id", "created_at", "updated_at"}

func TestGetProfile(t *testing.T) {
	now := time.Now()
	r, m := repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "profiles" WHERE user_id = \$1`).WithArgs("u", 1).WillReturnRows(sqlmock.NewRows(profileCols).AddRow("u", "n", "p", "a", "h", "c", now, now))
	got, e := r.GetProfile(context.Background(), "u")
	require.NoError(t, e)
	assert.Equal(t, "n", got.FullName)
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "profiles"`).WillReturnRows(sqlmock.NewRows(profileCols))
	_, e = r.GetProfile(context.Background(), "x")
	require.ErrorIs(t, e, infraerrors.ErrProfileNotFound)
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "profiles"`).WillReturnError(errors.New("db"))
	_, e = r.GetProfile(context.Background(), "x")
	assert.EqualError(t, e, "get profile: db")
}
func TestComplete(t *testing.T) {
	r, m := repoFixture(t)
	m.ExpectQuery(`SELECT count`).WithArgs("u").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	v, e := r.IsProfileComplete(context.Background(), "u")
	require.NoError(t, e)
	assert.True(t, v)
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT count`).WillReturnError(errors.New("db"))
	_, e = r.IsProfileComplete(context.Background(), "u")
	assert.EqualError(t, e, "is profile complete: db")
}
func TestCompanies(t *testing.T) {
	cols := []string{"id", "name", "created_at"}
	r, m := repoFixture(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "management_companies"`).WithArgs("UK", sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("c"))
	m.ExpectCommit()
	got, e := r.CreateManagementCompany(context.Background(), &domain.ManagementCompany{Name: "UK"})
	require.NoError(t, e)
	assert.Equal(t, "c", got.ID)
	r, m = repoFixture(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT`).WillReturnError(errors.New("db"))
	m.ExpectRollback()
	_, e = r.CreateManagementCompany(context.Background(), &domain.ManagementCompany{})
	assert.EqualError(t, e, "create management company: db")
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "management_companies" ORDER BY name`).WillReturnRows(sqlmock.NewRows(cols).AddRow("c", "UK", time.Now()))
	list, e := r.ListManagementCompanies(context.Background())
	require.NoError(t, e)
	assert.Len(t, list, 1)
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \*`).WillReturnError(errors.New("db"))
	_, e = r.ListManagementCompanies(context.Background())
	assert.EqualError(t, e, "list management companies: db")
}
func TestHouses(t *testing.T) {
	cols := []string{"id", "name", "address", "uk_id", "created_at"}
	r, m := repoFixture(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT INTO "houses"`).WithArgs("H", "A", "c", sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("h"))
	m.ExpectCommit()
	got, e := r.CreateHouse(context.Background(), &domain.House{Name: "H", Address: "A", UKID: "c"})
	require.NoError(t, e)
	assert.Equal(t, "h", got.ID)
	r, m = repoFixture(t)
	m.ExpectBegin()
	m.ExpectQuery(`INSERT`).WillReturnError(errors.New("db"))
	m.ExpectRollback()
	_, e = r.CreateHouse(context.Background(), &domain.House{})
	assert.EqualError(t, e, "create house: db")
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "houses" WHERE uk_id`).WithArgs("c").WillReturnRows(sqlmock.NewRows(cols).AddRow("h", "H", "A", "c", time.Now()))
	list, e := r.ListHouses(context.Background(), "c")
	require.NoError(t, e)
	assert.Len(t, list, 1)
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "houses" ORDER BY name`).WillReturnError(errors.New("db"))
	_, e = r.ListHouses(context.Background(), "")
	assert.EqualError(t, e, "list houses: db")
}
func TestUpsertProfile(t *testing.T) {
	houseCols := []string{"id", "name", "address", "uk_id", "created_at"}
	r, m := repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "houses" WHERE id`).WithArgs("h", 1).WillReturnRows(sqlmock.NewRows(houseCols))
	_, e := r.UpsertProfile(context.Background(), &domain.Profile{HouseID: "h"})
	require.ErrorIs(t, e, infraerrors.ErrHouseIDInvalid)
	r, m = repoFixture(t)
	m.ExpectQuery(`SELECT \* FROM "houses"`).WillReturnError(errors.New("db"))
	_, e = r.UpsertProfile(context.Background(), &domain.Profile{HouseID: "h"})
	assert.EqualError(t, e, "get house: db")
}
func TestConverters(t *testing.T) {
	now := time.Now()
	assert.Equal(t, "management_companies", (dbManagementCompany{}).TableName())
	assert.Equal(t, "houses", (dbHouse{}).TableName())
	assert.Equal(t, "profiles", (dbProfile{}).TableName())
	assert.Equal(t, &domain.ManagementCompany{ID: "c", Name: "C"}, companyToDomain(&dbManagementCompany{ID: "c", Name: "C"}))
	assert.Equal(t, &domain.House{ID: "h", Name: "H", Address: "A", UKID: "c"}, houseToDomain(&dbHouse{ID: "h", Name: "H", Address: "A", UKID: "c"}))
	assert.Equal(t, now, profileToDomain(&dbProfile{CreatedAt: now}).CreatedAt)
}
