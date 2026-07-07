package server

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"

	infraerrors "news-service/infra/errors"
	"news-service/infra/models/domain"
)

func TestValidators(t *testing.T) {
	long := strings.Repeat("x", 256)
	for _, tc := range []struct {
		name string
		req  *newspb.CreateNewsRequest
		want error
	}{
		{"title required", &newspb.CreateNewsRequest{}, infraerrors.ErrTitleRequired}, {"title short", &newspb.CreateNewsRequest{Title: "ab"}, infraerrors.ErrTitleTooShort},
		{"title long", &newspb.CreateNewsRequest{Title: long}, infraerrors.ErrTitleTooLong}, {"content required", &newspb.CreateNewsRequest{Title: "abc"}, infraerrors.ErrContentRequired},
		{"content short", &newspb.CreateNewsRequest{Title: "abc", Content: "short"}, infraerrors.ErrContentTooShort}, {"house required", &newspb.CreateNewsRequest{Title: "abc", Content: "long enough"}, infraerrors.ErrHouseIDRequired},
	} {
		t.Run(tc.name, func(t *testing.T) { require.ErrorIs(t, ValidateCreateNewsRequest(tc.req), tc.want) })
	}
	require.NoError(t, ValidateCreateNewsRequest(&newspb.CreateNewsRequest{Title: "abc", Content: "long enough", HouseId: "h"}))
	require.ErrorIs(t, ValidateGetNewsItemRequest(&newspb.GetNewsItemRequest{}), infraerrors.ErrNewsIDRequired)
	require.NoError(t, ValidateGetNewsItemRequest(&newspb.GetNewsItemRequest{Id: "n"}))
}

func TestConverters(t *testing.T) {
	created := protoToCreateNews(&newspb.CreateNewsRequest{Title: "t", Content: "c", HouseId: "h", CreatedBy: "u"})
	assert.Equal(t, &domain.CreateNewsRequest{Title: "t", Content: "c", HouseID: "h", CreatedBy: "u"}, created)
	from := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	to := from.Add(time.Hour)
	got := protoToGetNews(&newspb.GetNewsRequest{HouseId: "h", Limit: 10, Offset: 2, DateFrom: from.Format(time.RFC3339), DateTo: to.Format(time.RFC3339)})
	assert.Equal(t, &domain.GetNewsRequest{HouseID: "h", Limit: 10, Offset: 2, DateFrom: from, DateTo: to}, got)
	invalid := protoToGetNews(&newspb.GetNewsRequest{DateFrom: "bad", DateTo: "bad"})
	assert.True(t, invalid.DateFrom.IsZero())
	assert.True(t, invalid.DateTo.IsZero())
	pb := domainToProtoNews(&domain.News{ID: "n", Title: "t", Content: "c", HouseID: "h", CreatedBy: "u", CreatedAt: from})
	assert.Equal(t, "n", pb.Id)
	assert.Equal(t, from.Format(time.RFC3339), pb.CreatedAt)
}

func TestNewsServer(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	svc := NewMockNewsService(ctrl)
	logger := zerolog.Nop()
	server := NewNewsServer(svc, &logger)
	now := time.Now()
	n := &domain.News{ID: "n", CreatedAt: now}
	svc.EXPECT().CreateNews(ctx, gomock.Any()).Return(n, nil)
	resp, err := server.CreateNews(ctx, &newspb.CreateNewsRequest{})
	require.NoError(t, err)
	assert.Equal(t, "n", resp.News.Id)
	svc.EXPECT().GetNewsItem(ctx, "n").Return(n, nil)
	resp, err = server.GetNewsItem(ctx, &newspb.GetNewsItemRequest{Id: "n"})
	require.NoError(t, err)
	assert.Equal(t, "n", resp.News.Id)
	svc.EXPECT().GetNews(ctx, gomock.Any()).Return([]*domain.News{n}, int64(1), nil)
	list, err := server.GetNews(ctx, &newspb.GetNewsRequest{})
	require.NoError(t, err)
	assert.Len(t, list.News, 1)
	assert.Equal(t, int32(1), list.Total)

	ctrl = gomock.NewController(t)
	svc = NewMockNewsService(ctrl)
	server = NewNewsServer(svc, &logger)
	svc.EXPECT().CreateNews(ctx, gomock.Any()).Return(nil, errors.New("create"))
	_, err = server.CreateNews(ctx, &newspb.CreateNewsRequest{})
	assert.EqualError(t, err, "create")
	svc.EXPECT().GetNewsItem(ctx, "n").Return(nil, errors.New("get"))
	_, err = server.GetNewsItem(ctx, &newspb.GetNewsItemRequest{Id: "n"})
	assert.EqualError(t, err, "get")
	svc.EXPECT().GetNews(ctx, gomock.Any()).Return(nil, int64(0), errors.New("list"))
	_, err = server.GetNews(ctx, &newspb.GetNewsRequest{})
	assert.EqualError(t, err, "list")
}

func TestTotalToInt32(t *testing.T) {
	assert.Equal(t, int32(0), totalToInt32(-1))
	assert.Equal(t, int32(7), totalToInt32(7))
	assert.Equal(t, int32(math.MaxInt32), totalToInt32(math.MaxInt64))
}
