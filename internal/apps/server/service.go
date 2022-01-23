package server

import (
	"context"

	"github.com/dungnh3/guardrails-assignment/config"
	"github.com/dungnh3/guardrails-assignment/internal/repository"
	"github.com/dungnh3/guardrails-assignment/pkg/grpc/health_api"
	"github.com/go-logr/logr"
	"google.golang.org/genproto/googleapis/rpc/code"
	"gorm.io/gorm"
)

type Service struct {
	cfg    *config.Config
	logger logr.Logger
	repo   repository.IRepository
}

func NewService(cfg *config.Config, db *gorm.DB) (*Service, error) {
	logger := cfg.Logger.MustBuildLogR().WithName("guardrails-service")
	repo := repository.New(db, logger)
	return &Service{
		cfg:    cfg,
		logger: logger,
		repo:   repo,
	}, nil
}

func (s *Service) Close() error {
	return nil
}

// Readiness is a health check handler
func (s *Service) Readiness(ctx context.Context, request *health_api.ReadinessRequest) (*health_api.ReadinessResponse, error) {
	return &health_api.ReadinessResponse{
		Code:    code.Code_OK,
		Content: "OK",
	}, nil
}
