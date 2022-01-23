package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/dungnh3/guardrails-assignment/internal/model"
)

func (r *Repository) CreateSourceRepository(ctx context.Context, sr model.SourceRepository) (*model.SourceRepository, error) {
	tx := r.db.WithContext(ctx).Create(&sr)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &sr, nil
}

func (r *Repository) UpdateSourceRepositoryById(ctx context.Context, id uint32, name, link string) (*model.SourceRepository, error) {
	var sr model.SourceRepository
	tx := r.db.WithContext(ctx).Model(&sr).Where("id = ? AND is_active = ?", id, true).Updates(model.SourceRepository{
		Name: name,
		Link: link,
	})
	if err := tx.Error; err != nil {
		return nil, err
	}
	if rows := tx.RowsAffected; rows < 1 {
		return nil, ErrNotAnyRecordAffect
	}
	return &sr, nil
}

func (r *Repository) GetSourceRepositoryById(ctx context.Context, id uint32) (*model.SourceRepository, error) {
	var sr model.SourceRepository
	tx := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&sr)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &sr, nil
}

func (r *Repository) ListSourceRepository(ctx context.Context, nextId uint32, limit int) ([]model.SourceRepository, error) {
	var srs []model.SourceRepository
	tx := r.db.Debug().WithContext(ctx).Limit(limit).Where("id > ? AND is_active = ?", nextId, true).Find(&srs)
	if err := tx.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return srs, nil
}

func (r *Repository) RemoveSourceRepository(ctx context.Context, id uint32) error {
	var sr model.SourceRepository
	tx := r.db.WithContext(ctx).Model(&sr).Where("id = ? AND is_active = ?", id, true).
		Updates(map[string]interface{}{
			"is_active": false,
		})
	if err := tx.Error; err != nil {
		return err
	}
	if rows := tx.RowsAffected; rows < 1 {
		return ErrNotAnyRecordAffect
	}
	return nil
}
