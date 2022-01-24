package model

import (
	"time"

	"github.com/dungnh3/guardrails-assignment/internal/apps/rule"
)

type ScanningStatus string

const (
	QueuedStatus     ScanningStatus = "queued"
	InProgressStatus                = "in_progress"
	SuccessStatus                   = "success"
	FailureStatus                   = "failure"
)

type Result struct {
	ID                 uint32
	SourceRepositoryID uint32
	Name               string
	Link               string
	Status             ScanningStatus
	Findings           JSON
	QueuedAt           time.Time
	ScanningAt         *time.Time
	FinishedAt         *time.Time
}

type (
	Begin struct {
		Line int `json:"line"`
	}

	Positions struct {
		Begin Begin `json:"begin"`
	}

	Location struct {
		Path      string    `json:"path"`
		Positions Positions `json:"positions"`
	}

	Finding struct {
		Type     rule.Type     `json:"type"`
		RuleId   string        `json:"rule_id"`
		Metadata rule.Metadata `json:"metadata"`
		Location Location      `json:"location"`
	}
)
