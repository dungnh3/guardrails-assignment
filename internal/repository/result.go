package repository

import (
	"context"

	"github.com/dungnh3/guardrails-assignment/internal/model"
	"gorm.io/gorm/clause"
)

func (r *Repository) TriggerScanRepository(ctx context.Context, result model.Result) (*model.Result, error) {
	tx := r.db.Debug().WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "source_repository_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":      result.Name,
			"link":      result.Link,
			"status":    result.Status,
			"queued_at": result.QueuedAt,
		}),
	}).Create(&result)

	if err := tx.Error; err != nil {
		return nil, err
	}
	return &result, nil
}
