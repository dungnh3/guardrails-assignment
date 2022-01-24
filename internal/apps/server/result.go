package server

import (
	"context"
	"time"

	"github.com/dungnh3/guardrails-assignment/api"
	"github.com/dungnh3/guardrails-assignment/internal/model"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (s *Service) TriggerScanRepository(ctx context.Context, request *api.TriggerScanRepositoryRequest) (*api.TriggerScanRepositoryResponse, error) {
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	repo, err := s.repo.GetSourceRepositoryById(timeoutCtx, request.RepoId)
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
	result, err := s.repo.TriggerScanRepository(timeoutCtx, r)
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
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	rs, err := s.repo.ListResult(timeoutCtx, request.NextId, int(request.Limit))
	if err != nil {
		return nil, err
	}

	var results []*api.Result
	for _, r := range rs {
		result, err := formatResultToResponse(&r)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	nextId := uint32(0)
	if len(results) != 0 {
		nextId = rs[len(rs)-1].ID
	}
	return &api.ListResultResponse{
		Code:   code.Code_OK,
		NextId: nextId,
		Data: &api.ListResultResponse_Data{
			Results: results,
		},
	}, nil
}
