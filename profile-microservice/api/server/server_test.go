package server

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	infraerrors "profile-service/infra/errors"
	"profile-service/infra/models/domain"
	"strings"
	"testing"
)

func TestValidators(t *testing.T) {
	long101 := strings.Repeat("x", 101)
	long21 := strings.Repeat("x", 21)
	for _, tc := range []struct {
		name string
		run  func() error
		want error
	}{{"get user", func() error { return ValidateGetProfileRequest(&profilepb.GetProfileRequest{}) }, infraerrors.ErrUserIDRequired}, {"upsert user", func() error { return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{}) }, infraerrors.ErrUserIDRequired}, {"name", func() error { return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u"}) }, infraerrors.ErrFullNameRequired}, {"short", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "x"})
	}, infraerrors.ErrFullNameTooShort}, {"long", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: long101})
	}, infraerrors.ErrFullNameTooLong}, {"phone", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "xx"})
	}, infraerrors.ErrPhoneRequired}, {"phone long", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "xx", Phone: long21})
	}, infraerrors.ErrPhoneTooLong}, {"apt", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "xx", Phone: "p"})
	}, infraerrors.ErrApartmentRequired}, {"apt long", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "xx", Phone: "p", Apartment: long21})
	}, infraerrors.ErrApartmentTooLong}, {"house", func() error {
		return ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "xx", Phone: "p", Apartment: "a"})
	}, infraerrors.ErrHouseIDRequired}, {"company", func() error {
		return ValidateCreateManagementCompanyRequest(&profilepb.CreateManagementCompanyRequest{})
	}, infraerrors.ErrCompanyNameRequired}, {"company long", func() error {
		return ValidateCreateManagementCompanyRequest(&profilepb.CreateManagementCompanyRequest{Name: strings.Repeat("x", 201)})
	}, infraerrors.ErrCompanyNameTooLong}, {"house name", func() error { return ValidateCreateHouseRequest(&profilepb.CreateHouseRequest{}) }, infraerrors.ErrHouseNameRequired}, {"address", func() error { return ValidateCreateHouseRequest(&profilepb.CreateHouseRequest{Name: "h"}) }, infraerrors.ErrHouseAddressRequired}, {"uk", func() error {
		return ValidateCreateHouseRequest(&profilepb.CreateHouseRequest{Name: "h", Address: "a"})
	}, infraerrors.ErrUKIDRequired}} {
		t.Run(tc.name, func(t *testing.T) { require.ErrorIs(t, tc.run(), tc.want) })
	}
	require.NoError(t, ValidateGetProfileRequest(&profilepb.GetProfileRequest{UserId: "u"}))
	require.NoError(t, ValidateUpsertProfileRequest(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "xx", Phone: "p", Apartment: "a", HouseId: "h"}))
	require.NoError(t, ValidateCreateManagementCompanyRequest(&profilepb.CreateManagementCompanyRequest{Name: "c"}))
	require.NoError(t, ValidateCreateHouseRequest(&profilepb.CreateHouseRequest{Name: "h", Address: "a", UkId: "u"}))
}

func TestConverters(t *testing.T) {
	assert.Equal(t, &domain.UpsertProfileRequest{UserID: "u", FullName: "n", Phone: "p", Apartment: "a", HouseID: "h"}, protoToUpsertProfile(&profilepb.UpsertProfileRequest{UserId: "u", FullName: "n", Phone: "p", Apartment: "a", HouseId: "h"}))
	assert.Equal(t, &domain.CreateManagementCompanyRequest{Name: "c"}, protoToCreateCompany(&profilepb.CreateManagementCompanyRequest{Name: "c"}))
	assert.Equal(t, &domain.CreateHouseRequest{Name: "h", Address: "a", UKID: "u"}, protoToCreateHouse(&profilepb.CreateHouseRequest{Name: "h", Address: "a", UkId: "u"}))
	assert.Equal(t, &domain.ListHousesRequest{UKID: "u"}, protoToListHouses(&profilepb.ListHousesRequest{UkId: "u"}))
	assert.Equal(t, "u", domainToProtoProfile(&domain.Profile{UserID: "u"}).UserId)
	assert.Equal(t, "c", domainToProtoCompany(&domain.ManagementCompany{ID: "c"}).Id)
	assert.Equal(t, "h", domainToProtoHouse(&domain.House{ID: "h"}).Id)
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	svc := NewMockProfileService(ctrl)
	l := zerolog.Nop()
	s := NewProfileServer(svc, &l)
	p := &domain.Profile{UserID: "u"}
	c := &domain.ManagementCompany{ID: "c"}
	h := &domain.House{ID: "h"}
	svc.EXPECT().GetProfile(ctx, "u").Return(p, nil)
	r, e := s.GetProfile(ctx, &profilepb.GetProfileRequest{UserId: "u"})
	require.NoError(t, e)
	assert.Equal(t, "u", r.Profile.UserId)
	svc.EXPECT().UpsertProfile(ctx, gomock.Any()).Return(p, nil)
	_, e = s.UpsertProfile(ctx, &profilepb.UpsertProfileRequest{})
	require.NoError(t, e)
	svc.EXPECT().IsProfileComplete(ctx, "u").Return(true, nil)
	complete, e := s.IsProfileComplete(ctx, &profilepb.IsProfileCompleteRequest{UserId: "u"})
	require.NoError(t, e)
	assert.True(t, complete.Complete)
	svc.EXPECT().CreateManagementCompany(ctx, gomock.Any()).Return(c, nil)
	_, e = s.CreateManagementCompany(ctx, &profilepb.CreateManagementCompanyRequest{})
	require.NoError(t, e)
	svc.EXPECT().ListManagementCompanies(ctx).Return([]*domain.ManagementCompany{c}, nil)
	cs, e := s.ListManagementCompanies(ctx, &profilepb.ListManagementCompaniesRequest{})
	require.NoError(t, e)
	assert.Len(t, cs.Companies, 1)
	svc.EXPECT().CreateHouse(ctx, gomock.Any()).Return(h, nil)
	_, e = s.CreateHouse(ctx, &profilepb.CreateHouseRequest{})
	require.NoError(t, e)
	svc.EXPECT().ListHouses(ctx, gomock.Any()).Return([]*domain.House{h}, nil)
	hs, e := s.ListHouses(ctx, &profilepb.ListHousesRequest{})
	require.NoError(t, e)
	assert.Len(t, hs.Houses, 1)
	ctrl = gomock.NewController(t)
	svc = NewMockProfileService(ctrl)
	s = NewProfileServer(svc, &l)
	svc.EXPECT().GetProfile(ctx, "x").Return(nil, errors.New("get"))
	_, e = s.GetProfile(ctx, &profilepb.GetProfileRequest{UserId: "x"})
	assert.EqualError(t, e, "get")
	svc.EXPECT().UpsertProfile(ctx, gomock.Any()).Return(nil, errors.New("upsert"))
	_, e = s.UpsertProfile(ctx, &profilepb.UpsertProfileRequest{})
	assert.EqualError(t, e, "upsert")
	svc.EXPECT().IsProfileComplete(ctx, "x").Return(false, errors.New("complete"))
	_, e = s.IsProfileComplete(ctx, &profilepb.IsProfileCompleteRequest{UserId: "x"})
	assert.EqualError(t, e, "complete")
	svc.EXPECT().CreateManagementCompany(ctx, gomock.Any()).Return(nil, errors.New("company"))
	_, e = s.CreateManagementCompany(ctx, &profilepb.CreateManagementCompanyRequest{})
	assert.EqualError(t, e, "company")
	svc.EXPECT().ListManagementCompanies(ctx).Return(nil, errors.New("companies"))
	_, e = s.ListManagementCompanies(ctx, nil)
	assert.EqualError(t, e, "companies")
	svc.EXPECT().CreateHouse(ctx, gomock.Any()).Return(nil, errors.New("house"))
	_, e = s.CreateHouse(ctx, &profilepb.CreateHouseRequest{})
	assert.EqualError(t, e, "house")
	svc.EXPECT().ListHouses(ctx, gomock.Any()).Return(nil, errors.New("houses"))
	_, e = s.ListHouses(ctx, &profilepb.ListHousesRequest{})
	assert.EqualError(t, e, "houses")
}
