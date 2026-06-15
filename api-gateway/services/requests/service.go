package requests_service

import (
	"context"

	"github.com/rs/zerolog"
)

type RequestsService struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *RequestsService {
	return &RequestsService{logger: logger}
}

func (s *RequestsService) Create(ctx context.Context) error {
	return nil
}

func (s *RequestsService) List(ctx context.Context) error {
	return nil
}

func (s *RequestsService) Get(ctx context.Context) error {
	return nil
}

func (s *RequestsService) UpdateStatus(ctx context.Context) error {
	return nil
}

func (s *RequestsService) AddComment(ctx context.Context) error {
	return nil
}
