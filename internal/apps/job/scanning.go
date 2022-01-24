package job

import (
	"bufio"
	"context"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dungnh3/guardrails-assignment/internal/apps/rule"

	"github.com/dungnh3/guardrails-assignment/internal/model"

	"github.com/dungnh3/guardrails-assignment/internal/repository"
	"github.com/go-git/go-git/v5"
	"github.com/go-logr/logr"
	"github.com/speedata/gogit"
	"gorm.io/gorm"
)

const interval = 100 * time.Millisecond

type ScanningEngine struct {
	logger logr.Logger
	repo   repository.IRepository
	rules  []rule.IChecked
}

func NewScanning(logger logr.Logger, db *gorm.DB, rules []rule.IChecked) *ScanningEngine {
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
	// TODO implement me
	panic("implement me")
}

func walk(dirname string, te *gogit.TreeEntry) int {
	log.Println(path.Join(dirname, te.Name))
	return 0
}

func (s *ScanningEngine) cloneRepository(link string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	url, err := url.Parse(link)
	if err != nil {
		return err
	}
	filePath := filepath.Join(wd, "tmp", path.Base(url.Path))
	_, err = git.PlainClone(filePath, false, &git.CloneOptions{
		URL:      "link",
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ScanningEngine) process(ctx context.Context) error {
	var err error
	//wd, err := os.Getwd()
	//if err != nil {
	//	return err
	//}
	//
	//url, err := url.Parse("https://github.com/dungnh3/buf-demo")
	//if err != nil {
	//	return err
	//}
	//filePath := filepath.Join(wd, "tmp", path.Base(url.Path))
	//gitRepo, err := git.PlainClone(filePath, false, &git.CloneOptions{
	//	URL:      "https://github.com/dungnh3/buf-demo",
	//	Progress: os.Stdout,
	//})
	//if err != nil {
	//	return err
	//}

	//commits, err := gitRepo.CommitObjects()
	//if err != nil {
	//	s.logger.Error(err, "scan error")
	//	return err
	//}
	//defer commits.Close()
	//
	//latestCommit, err := commits.Next()
	//if err != nil {
	//	return err
	//}
	//fileIter, err := latestCommit.Files()
	//if err != nil {
	//	return err
	//}
	//defer fileIter.Close()

	findings, err := s.scan("./tmp/buf-demo")
	if err != nil {
		return err
	}

	if err = s.repo.UpdateFindingsResultSuccess(ctx, 4, findings); err != nil {
		return err
	}
	return nil
}

func (s *ScanningEngine) scan(pathRoot string) ([]model.Finding, error) {
	var findings []model.Finding
	if err := filepath.Walk(pathRoot, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
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
				ok, err := r.Detect(scanner.Text())
				if err != nil {
					return err
				}
				if ok {
					findings = append(findings, model.Finding{
						Location: model.Location{
							Path: path,
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
