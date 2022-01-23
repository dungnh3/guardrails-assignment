package server

import (
	"context"

	"github.com/dungnh3/guardrails-assignment/api"
	"github.com/dungnh3/guardrails-assignment/internal/model"
	"google.golang.org/genproto/googleapis/rpc/code"
)

func (s *Service) CreateRepository(ctx context.Context, request *api.CreateRepositoryRequest) (*api.CreateRepositoryResponse, error) {
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	sr, err := s.repo.CreateSourceRepository(timeoutCtx, model.SourceRepository{
		Name:     request.Name,
		Link:     request.Link,
		IsActive: true,
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateRepositoryResponse{
		Code: code.Code_OK,
		Data: &api.CreateRepositoryResponse_Data{
			SourceRepository: formatSourceRepositoryToResponse(sr),
		},
	}, nil
}

func (s *Service) GetRepositoryById(ctx context.Context, request *api.GetRepositoryByIdRequest) (*api.GetRepositoryByIdResponse, error) {
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	sr, err := s.repo.GetSourceRepositoryById(timeoutCtx, request.Id)
	if err != nil {
		return nil, err
	}
	return &api.GetRepositoryByIdResponse{
		Code: code.Code_OK,
		Data: &api.GetRepositoryByIdResponse_Data{
			SourceRepository: formatSourceRepositoryToResponse(sr),
		},
	}, nil
}

func (s *Service) ListRepository(ctx context.Context, request *api.ListRepositoryRequest) (*api.ListRepositoryResponse, error) {
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	srs, err := s.repo.ListSourceRepository(timeoutCtx, request.NextId, int(request.Limit))
	if err != nil {
		return nil, err
	}

	var repositories []*api.SourceRepository
	for _, sr := range srs {
		repositories = append(repositories, formatSourceRepositoryToResponse(&sr))
	}

	nextId := uint32(0)
	if len(repositories) != 0 {
		nextId = repositories[len(repositories)-1].Id
	}
	return &api.ListRepositoryResponse{
		Code:   code.Code_OK,
		NextId: nextId,
		Data: &api.ListRepositoryResponse_Data{
			SourceRepositories: repositories,
		},
	}, nil
}

func (s *Service) UpdateRepository(ctx context.Context, request *api.UpdateRepositoryRequest) (*api.UpdateRepositoryResponse, error) {
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	sr, err := s.repo.UpdateSourceRepositoryById(timeoutCtx, request.Id, request.Name, request.Link)
	if err != nil {
		return nil, err
	}

	return &api.UpdateRepositoryResponse{
		Code: code.Code_OK,
		Data: &api.UpdateRepositoryResponse_Data{
			SourceRepository: formatSourceRepositoryToResponse(sr),
		},
	}, nil
}

func (s *Service) RemoveRepository(ctx context.Context, request *api.RemoveRepositoryRequest) (*api.RemoveRepositoryResponse, error) {
	timeoutCtx, cancel := defaultTimeoutContext(ctx)
	defer cancel()

	if err := s.repo.RemoveSourceRepository(timeoutCtx, request.Id); err != nil {
		return nil, err
	}
	return &api.RemoveRepositoryResponse{
		Code: code.Code_OK,
	}, nil
}
