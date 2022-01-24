package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

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

func (r *Repository) GetTriggerRepository(ctx context.Context) ([]model.Result, error) {
	// TODO implement me
	panic("implement me")
}

func (r *Repository) UpdateFindingsResultSuccess(ctx context.Context, id uint32, findings []model.Finding) error {
	var result model.Result

	data, err := json.Marshal(findings)
	if err != nil {
		return err
	}

	var findingsJson model.JSON
	if err = findingsJson.Scan(data); err != nil {
		return err
	}
	tx := r.db.WithContext(ctx).Model(&result).Where("id = ?", id).
		Updates(map[string]interface{}{
			"finished_at": time.Now(),
			"findings":    findingsJson,
		})
	if err := tx.Error; err != nil {
		return err
	}
	if rows := tx.RowsAffected; rows < 1 {
		return ErrNotAnyRecordAffect
	}
	return nil
}

func (r *Repository) ListResult(ctx context.Context, nextId uint32, limit int) ([]model.Result, error) {
	var results []model.Result
	tx := r.db.Debug().WithContext(ctx).Limit(limit).Where("id > ?", nextId).Find(&results)
	if err := tx.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return results, nil
}
