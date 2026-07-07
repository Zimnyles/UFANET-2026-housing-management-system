package requests_service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"requests-service/infra/models/domain"
	"testing"
)

func setup(t *testing.T) (*MockRepository, *MockPublisher, *RequestsService) {
	c := gomock.NewController(t)
	r, p := NewMockRepository(c), NewMockPublisher(c)
	l := zerolog.Nop()
	return r, p, New(r, p, &l)
}
func TestService(t *testing.T) {
	ctx := context.Background()
	r, p, s := setup(t)
	cr := &domain.CreateRequestRequest{Title: "t", Description: "d", Type: "plumber", UserID: "u"}
	item := &domain.MaintenanceRequest{ID: "r", Title: "t", Description: "d", Type: "plumber", UserID: "u", Status: "open"}
	r.EXPECT().Create(ctx, gomock.Any()).Return(item, nil)
	p.EXPECT().PublishRequestCreated(ctx, gomock.Any()).Return(nil)
	got, e := s.CreateRequest(ctx, cr)
	require.NoError(t, e)
	assert.Same(t, item, got)
	r.EXPECT().List(ctx, gomock.Any()).Return([]*domain.MaintenanceRequest{item}, int64(1), nil)
	items, total, e := s.GetRequests(ctx, &domain.GetRequestsRequest{})
	require.NoError(t, e)
	assert.Len(t, items, 1)
	assert.Equal(t, int64(1), total)
	r.EXPECT().GetByID(ctx, "r").Return(item, nil)
	_, e = s.GetRequest(ctx, "r")
	require.NoError(t, e)
	r.EXPECT().UpdateStatus(ctx, "r", "done").Return(item, nil)
	p.EXPECT().PublishRequestStatusUpdated(ctx, gomock.Any()).Return(errors.New("mq"))
	_, e = s.UpdateRequestStatus(ctx, &domain.UpdateStatusRequest{ID: "r", Status: "done"})
	require.NoError(t, e)
	comment := &domain.Comment{ID: "c", RequestID: "r", UserID: "u", Content: "x"}
	r.EXPECT().AddComment(ctx, gomock.Any()).Return(comment, nil)
	p.EXPECT().PublishRequestCommentAdded(ctx, gomock.Any()).Return(nil)
	_, e = s.AddComment(ctx, &domain.AddCommentRequest{RequestID: "r", UserID: "u", Content: "x"})
	require.NoError(t, e)
	r.EXPECT().ListComments(ctx, "r").Return([]*domain.Comment{comment}, nil)
	cs, e := s.GetComments(ctx, "r")
	require.NoError(t, e)
	assert.Len(t, cs, 1)
}
func TestServiceErrorsAndNilPublisher(t *testing.T) {
	ctx := context.Background()
	r, _, s := setup(t)
	r.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("db"))
	_, e := s.CreateRequest(ctx, &domain.CreateRequestRequest{})
	assert.EqualError(t, e, "create request: db")
	r.EXPECT().UpdateStatus(ctx, "r", "x").Return(nil, errors.New("db"))
	_, e = s.UpdateRequestStatus(ctx, &domain.UpdateStatusRequest{ID: "r", Status: "x"})
	assert.EqualError(t, e, "update request status: db")
	r.EXPECT().AddComment(ctx, gomock.Any()).Return(nil, errors.New("db"))
	_, e = s.AddComment(ctx, &domain.AddCommentRequest{})
	assert.EqualError(t, e, "add comment: db")
	c := gomock.NewController(t)
	repo := NewMockRepository(c)
	l := zerolog.Nop()
	s = New(repo, nil, &l)
	item := &domain.MaintenanceRequest{}
	repo.EXPECT().Create(ctx, gomock.Any()).Return(item, nil)
	_, e = s.CreateRequest(ctx, &domain.CreateRequestRequest{})
	require.NoError(t, e)
	repo.EXPECT().UpdateStatus(ctx, "", "").Return(item, nil)
	_, e = s.UpdateRequestStatus(ctx, &domain.UpdateStatusRequest{})
	require.NoError(t, e)
	comment := &domain.Comment{}
	repo.EXPECT().AddComment(ctx, gomock.Any()).Return(comment, nil)
	_, e = s.AddComment(ctx, &domain.AddCommentRequest{})
	require.NoError(t, e)
}
