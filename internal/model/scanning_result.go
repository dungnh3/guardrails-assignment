package model

import "time"

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
