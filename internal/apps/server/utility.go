package server

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/dungnh3/guardrails-assignment/api"

	"github.com/dungnh3/guardrails-assignment/internal/model"
)

const (
	Duration10 = 10 * time.Second
)

func defaultTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, Duration10)
}

func formatSourceRepositoryToResponse(sr *model.SourceRepository) *api.SourceRepository {
	if sr == nil {
		return nil
	}
	return &api.SourceRepository{
		Id:        sr.ID,
		Name:      sr.Name,
		Link:      sr.Link,
		CreatedAt: formatTimestampPbToResponse(&sr.CreatedAt),
		UpdatedAt: formatTimestampPbToResponse(sr.UpdatedAt),
	}
}

func formatTimestampPbToResponse(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
