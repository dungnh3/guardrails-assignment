package repository

import (
	"context"

	"github.com/dungnh3/guardrails-assignment/internal/model"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

type IRepository interface {
	Transaction(txFunc func(IRepository) error) error
	CreateSourceRepository(ctx context.Context, sr model.SourceRepository) (*model.SourceRepository, error)
	UpdateSourceRepositoryById(ctx context.Context, id uint32, name, link string) (*model.SourceRepository, error)
	GetSourceRepositoryById(ctx context.Context, id uint32) (*model.SourceRepository, error)
	ListSourceRepository(ctx context.Context, nextId uint32, limit int) ([]model.SourceRepository, error)
	RemoveSourceRepository(ctx context.Context, id uint32) error
	TriggerScanRepository(ctx context.Context, result model.Result) (*model.Result, error)
}

type Repository struct {
	db     *gorm.DB
	logger logr.Logger
}

func New(db *gorm.DB, logger logr.Logger) IRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) WithTx(tx *gorm.DB) *Repository {
	newRepo := *r
	newRepo.db = tx
	return &newRepo
}

func (r *Repository) Transaction(txFunc func(IRepository) error) error {
	tx := r.db.Begin()
	defer func() {
		if rc := recover(); rc != nil {
			r.logger.Error(nil, "rollback now because listening recover: %v \n", rc)
			if execErr := tx.Rollback().Error; execErr != nil {
				r.logger.Error(execErr, "exception error when execute rollback")
			}
			panic(rc)
		}
	}()

	err := txFunc(r.WithTx(tx))
	if err != nil {
		if execErr := tx.Rollback().Error; execErr != nil {
			r.logger.Error(execErr, "exception error when execute rollback")
		}
		return err
	}
	return tx.Commit().Error
}
