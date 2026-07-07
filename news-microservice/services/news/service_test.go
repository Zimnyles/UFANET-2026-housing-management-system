package news_service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"news-service/infra/models/domain"
)

func TestNewsService(t *testing.T) {
	ctx := context.Background()
	req := &domain.CreateNewsRequest{Title: "title", Content: "content text", HouseID: "h", CreatedBy: "u"}
	t.Run("create and publish", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		pub := NewMockPublisher(ctrl)
		logger := zerolog.Nop()
		svc := New(repo, pub, &logger)
		created := &domain.News{ID: "n", Title: req.Title, Content: req.Content, HouseID: "h", CreatedBy: "u", CreatedAt: time.Now()}
		repo.EXPECT().Create(ctx, &domain.News{Title: req.Title, Content: req.Content, HouseID: "h", CreatedBy: "u"}).Return(created, nil)
		pub.EXPECT().PublishNewsCreated(ctx, domain.NewsCreatedEvent{Type: "news.created", NewsID: "n", Title: "title", HouseID: "h"}).Return(nil)
		got, err := svc.CreateNews(ctx, req)
		require.NoError(t, err)
		assert.Same(t, created, got)
	})
	t.Run("publish error does not fail create", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		pub := NewMockPublisher(ctrl)
		logger := zerolog.Nop()
		svc := New(repo, pub, &logger)
		created := &domain.News{ID: "n"}
		repo.EXPECT().Create(ctx, gomock.Any()).Return(created, nil)
		pub.EXPECT().PublishNewsCreated(ctx, gomock.Any()).Return(errors.New("mq"))
		got, err := svc.CreateNews(ctx, req)
		require.NoError(t, err)
		assert.Same(t, created, got)
	})
	t.Run("nil publisher and repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		logger := zerolog.Nop()
		svc := New(repo, nil, &logger)
		repo.EXPECT().Create(ctx, gomock.Any()).Return(&domain.News{ID: "n"}, nil)
		_, err := svc.CreateNews(ctx, req)
		require.NoError(t, err)
		repo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("db"))
		_, err = svc.CreateNews(ctx, req)
		assert.EqualError(t, err, "create news: db")
	})
	t.Run("get and list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		logger := zerolog.Nop()
		svc := New(repo, nil, &logger)
		n := &domain.News{ID: "n"}
		repo.EXPECT().GetByID(ctx, "n").Return(n, nil)
		got, err := svc.GetNewsItem(ctx, "n")
		require.NoError(t, err)
		assert.Same(t, n, got)
		listReq := &domain.GetNewsRequest{HouseID: "h"}
		repo.EXPECT().List(ctx, listReq).Return([]*domain.News{n}, int64(1), nil)
		items, total, err := svc.GetNews(ctx, listReq)
		require.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, int64(1), total)
	})
}
