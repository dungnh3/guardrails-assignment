package job

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/dungnh3/guardrails-assignment/internal/apps/rule"

	"github.com/dungnh3/guardrails-assignment/internal/model"

	"github.com/dungnh3/guardrails-assignment/internal/repository"
	"github.com/go-git/go-git/v5"
	"github.com/go-logr/logr"
	"gorm.io/gorm"
)

const interval = 100 * time.Millisecond

const gitSkipDir = ".git"

const processBatch = 10

type ScanningEngine struct {
	logger logr.Logger
	repo   repository.IRepository
	rules  []rule.Rule
}

func NewScanner(logger logr.Logger, db *gorm.DB, rules ...rule.Rule) *ScanningEngine {
	logger = logger.WithName("guardrails-scanning-jobs")
	repo := repository.New(db, logger)
	return &ScanningEngine{
		logger: logger,
		repo:   repo,
		rules:  rules,
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
	return nil
}

func (s *ScanningEngine) process(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	results, err := s.repo.GetQueuedTriggerRepository(ctx, processBatch)
	if err != nil && err != repository.ErrRecordNotFound {
		return err
	}
	s.logger.Info(fmt.Sprintf("process scan repositories [%v]", len(results)))

	for _, result := range results {
		findings, err := s.scan(wd, result)
		if err != nil {
			return err
		}

		if len(findings) > 0 {
			if err = s.repo.UpdateFindingsResultFailure(ctx, result.ID, findings); err != nil {
				s.logger.Error(err, "UpdateFindingsResultFailure failed", "id", result.ID, "findings", findings)
				return errors.Wrap(err, "UpdateFindingsResultFailure failed")
			}
		} else {
			if err = s.repo.UpdateFindingsResultSuccess(ctx, result.ID); err != nil {
				s.logger.Error(err, "UpdateFindingsResultFailure failed", "id", result.ID)
				return errors.Wrap(err, "UpdateFindingsResultSuccess failed")
			}
		}
	}
	return nil
}

func (s *ScanningEngine) scan(wd string, result model.Result) ([]model.Finding, error) {
	filename := fmt.Sprintf("%v-%v", result.Name, uuid.New())
	filePath := filepath.Join(wd, "tmp", filename)
	if err := s.cloneRepository(filePath, result.Link); err != nil {
		s.logger.Error(err, "clone repository failed", "url", result.Link)
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(filePath); err != nil {
			s.logger.Error(err, "remove file failed", "filepath", filePath)
		}
	}()

	var findings []model.Finding
	pathRoot := filepath.Join("tmp", filename)
	if err := filepath.Walk(pathRoot, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == gitSkipDir {
			return filepath.SkipDir
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		counter := 0
		for scanner.Scan() {
			counter++
			for _, r := range s.rules {
				isErrDetect, err := r.DetectFn(scanner.Text())
				if err != nil {
					return err
				}
				if isErrDetect {
					findings = append(findings, model.Finding{
						Type:     r.Type,
						RuleId:   r.RuleId,
						Metadata: r.Metadata,
						Location: model.Location{
							Path: func() string {
								fp := strings.Replace(path, pathRoot, "", 1)
								return strings.Replace(fp, "/", "", 1)
							}(),
							Positions: model.Positions{
								Begin: model.Begin{
									Line: counter,
								},
							},
						},
					})
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return findings, nil
}

func (s *ScanningEngine) cloneRepository(filePath, url string) error {
	if _, err := git.PlainClone(filePath, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}); err != nil {
		return err
	}
	return nil
}
