package profile_service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"profile-service/infra/models/domain"
	"testing"
)

func setup(t *testing.T) (*MockRepository, *ProfileService) {
	r := NewMockRepository(gomock.NewController(t))
	l := zerolog.Nop()
	return r, New(r, &l)
}
func TestProfileService(t *testing.T) {
	ctx := context.Background()
	t.Run("get", func(t *testing.T) {
		r, s := setup(t)
		p := &domain.Profile{UserID: "u"}
		r.EXPECT().GetProfile(ctx, "u").Return(p, nil)
		got, e := s.GetProfile(ctx, "u")
		require.NoError(t, e)
		assert.Same(t, p, got)
		r.EXPECT().GetProfile(ctx, "x").Return(nil, errors.New("db"))
		_, e = s.GetProfile(ctx, "x")
		assert.EqualError(t, e, "db")
	})
	t.Run("upsert", func(t *testing.T) {
		r, s := setup(t)
		req := &domain.UpsertProfileRequest{UserID: "u", FullName: "N", Phone: "P", Apartment: "A", HouseID: "H"}
		want := &domain.Profile{UserID: "u", FullName: "N", Phone: "P", Apartment: "A", HouseID: "H"}
		r.EXPECT().UpsertProfile(ctx, want).Return(want, nil)
		got, e := s.UpsertProfile(ctx, req)
		require.NoError(t, e)
		assert.Same(t, want, got)
		r.EXPECT().UpsertProfile(ctx, gomock.Any()).Return(nil, errors.New("db"))
		_, e = s.UpsertProfile(ctx, req)
		assert.EqualError(t, e, "upsert profile: db")
	})
	t.Run("complete", func(t *testing.T) {
		r, s := setup(t)
		r.EXPECT().IsProfileComplete(ctx, "u").Return(true, nil)
		v, e := s.IsProfileComplete(ctx, "u")
		require.NoError(t, e)
		assert.True(t, v)
		r.EXPECT().IsProfileComplete(ctx, "x").Return(false, errors.New("db"))
		_, e = s.IsProfileComplete(ctx, "x")
		assert.EqualError(t, e, "db")
	})
	t.Run("companies", func(t *testing.T) {
		r, s := setup(t)
		req := &domain.CreateManagementCompanyRequest{Name: "UK"}
		c := &domain.ManagementCompany{ID: "c", Name: "UK"}
		r.EXPECT().CreateManagementCompany(ctx, &domain.ManagementCompany{Name: "UK"}).Return(c, nil)
		got, e := s.CreateManagementCompany(ctx, req)
		require.NoError(t, e)
		assert.Same(t, c, got)
		r.EXPECT().CreateManagementCompany(ctx, gomock.Any()).Return(nil, errors.New("db"))
		_, e = s.CreateManagementCompany(ctx, req)
		assert.EqualError(t, e, "create management company: db")
		r.EXPECT().ListManagementCompanies(ctx).Return([]*domain.ManagementCompany{c}, nil)
		list, e := s.ListManagementCompanies(ctx)
		require.NoError(t, e)
		assert.Len(t, list, 1)
	})
	t.Run("houses", func(t *testing.T) {
		r, s := setup(t)
		req := &domain.CreateHouseRequest{Name: "H", Address: "A", UKID: "c"}
		h := &domain.House{ID: "h", Name: "H", Address: "A", UKID: "c"}
		r.EXPECT().CreateHouse(ctx, &domain.House{Name: "H", Address: "A", UKID: "c"}).Return(h, nil)
		got, e := s.CreateHouse(ctx, req)
		require.NoError(t, e)
		assert.Same(t, h, got)
		r.EXPECT().CreateHouse(ctx, gomock.Any()).Return(nil, errors.New("db"))
		_, e = s.CreateHouse(ctx, req)
		assert.EqualError(t, e, "create house: db")
		r.EXPECT().ListHouses(ctx, "c").Return([]*domain.House{h}, nil)
		list, e := s.ListHouses(ctx, &domain.ListHousesRequest{UKID: "c"})
		require.NoError(t, e)
		assert.Len(t, list, 1)
	})
}
