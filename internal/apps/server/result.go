package server

import (
	"context"
	"time"

	"github.com/dungnh3/guardrails-assignment/api"
	"github.com/dungnh3/guardrails-assignment/internal/model"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (s *Service) TriggerScanRepository(ctx context.Context, request *api.TriggerScanRepositoryRequest) (*api.TriggerScanRepositoryResponse, error) {
	repo, err := s.repo.GetSourceRepositoryById(ctx, request.RepoId)
	if err != nil {
		return nil, err
	}

	r := model.Result{
		SourceRepositoryID: repo.ID,
		Name:               repo.Name,
		Link:               repo.Link,
		Status:             model.QueuedStatus,
		QueuedAt:           time.Now(),
	}
	result, err := s.repo.TriggerScanRepository(ctx, r)
	if err != nil {
		return nil, err
	}
	return &api.TriggerScanRepositoryResponse{
		Code:     code.Code_OK,
		Id:       result.ID,
		QueuedAt: formatTimestampPbToResponse(&result.QueuedAt),
	}, nil
}

func (s *Service) ListResult(ctx context.Context, request *api.ListResultRequest) (*api.ListResultResponse, error) {
	// TODO implement me
	panic("implement me")
}
