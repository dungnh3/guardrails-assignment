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
	tx := r.db.WithContext(ctx).Clauses(clause.OnConflict{
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

const getQueuedTriggerRepositorySQL = `
UPDATE results
SET scanning_at = NOW(), status = 'in_progress'
WHERE id IN (
    SELECT r.id
    FROM results r
    WHERE r.status = 'queued'
    ORDER BY r.queued_at
    LIMIT ?
)
RETURNING results.*;
`

func (r *Repository) GetQueuedTriggerRepository(ctx context.Context, limit int) ([]model.Result, error) {
	var results []model.Result
	tx := r.db.WithContext(ctx).Raw(getQueuedTriggerRepositorySQL, limit).Scan(&results)
	if err := tx.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if rows := tx.RowsAffected; rows < 1 {
		return nil, ErrRecordNotFound
	}
	return results, nil
}

func (r *Repository) UpdateFindingsResultSuccess(ctx context.Context, id uint32) error {
	var result model.Result
	tx := r.db.WithContext(ctx).Model(&result).Where("id = ?", id).
		Updates(map[string]interface{}{
			"finished_at": time.Now(),
			"status":      model.SuccessStatus,
		})
	if err := tx.Error; err != nil {
		return err
	}
	if rows := tx.RowsAffected; rows < 1 {
		return ErrNotAnyRecordAffect
	}
	return nil
}

func (r *Repository) UpdateFindingsResultFailure(ctx context.Context, id uint32, findings []model.Finding) error {
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
			"status":      model.FailureStatus,
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
	tx := r.db.WithContext(ctx).Limit(limit).Where("id > ?", nextId).Find(&results)
	if err := tx.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return results, nil
}
