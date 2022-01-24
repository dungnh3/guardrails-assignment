package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/pkg/errors"

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

func formatResultToResponse(result *model.Result) (*api.Result, error) {
	if result == nil {
		return nil, nil
	}

	findings, err := formatFindingsToResponse(result.Findings)
	if err != nil {
		return nil, err
	}
	return &api.Result{
		Id:                 result.ID,
		SourceRepositoryId: result.SourceRepositoryID,
		Name:               result.Name,
		Link:               result.Link,
		Status:             string(result.Status),
		Findings:           findings,
		QueuedAt:           formatTimestampPbToResponse(&result.QueuedAt),
		ScanningAt:         formatTimestampPbToResponse(result.ScanningAt),
		FinishedAt:         formatTimestampPbToResponse(result.FinishedAt),
	}, nil
}

func formatFindingsToResponse(j model.JSON) ([]*structpb.Struct, error) {
	if j == nil {
		return nil, nil
	}

	var m []*structpb.Struct
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unmarshal json failed, value: %v", string(j)))
	}
	return m, nil
}

func formatTimestampPbToResponse(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
