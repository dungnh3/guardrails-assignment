package job

import (
	"context"
	"log"
	"path"
	"time"

	"github.com/speedata/gogit"

	"github.com/dungnh3/guardrails-assignment/config"
	"gorm.io/gorm"

	"github.com/go-logr/logr"

	"github.com/dungnh3/guardrails-assignment/internal/repository"
)

const interval = 100 * time.Millisecond

type ScanningEngine struct {
	logger logr.Logger
	repo   repository.IRepository
}

func NewScanning(cfg *config.Config, db *gorm.DB) *ScanningEngine {
	logger := cfg.Logger.MustBuildLogR().WithName("guardrails-scanning-jobs")
	repo := repository.New(db, logger)
	return &ScanningEngine{
		logger: logger,
		repo:   repo,
	}
}

func (s *ScanningEngine) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Nanosecond)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ticker.Stop()

			if err := s.process(ctx); err != nil {
				s.logger.Error(err, "scan repository", "time", time.Now().Unix())
			}
			ticker.Reset(interval)
		}
	}
}

func (s *ScanningEngine) Close(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func walk(dirname string, te *gogit.TreeEntry) int {
	log.Println(path.Join(dirname, te.Name))
	return 0
}

func (s *ScanningEngine) process(ctx context.Context) error {
	s.logger.Info("hello")
	return nil
	//wd, err := os.Getwd()
	//if err != nil {
	//	return err
	//}
	//r, err := gogit.OpenRepository(filepath.Join(wd, "src/github.com/speedata/gogit/_testdata/testrepo.git"))
	//if err != nil {
	//	return err
	//}
	//ref, err := r.LookupReference("HEAD")
	//if err != nil {
	//	return err
	//}
	//ci, err := r.LookupCommit(ref.Oid)
	//if err != nil {
	//	return err
	//}
	//return ci.Tree.Walk(walk)
}
