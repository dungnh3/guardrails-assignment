package model

import "time"

type ScanningStatus string

const (
	QueuedStatus     ScanningStatus = "Queued"
	InProgressStatus                = "InProgress"
	SuccessStatus                   = "Success"
	FailureStatus                   = "Failure"
)

type Result struct {
	ID             uint32
	RepositoryID   uint32
	RepositoryName string
	RepositoryUrl  string
	Status         ScanningStatus
	Findings       JSON
	QueuedAt       time.Time
	ScanningAt     *time.Time
	FinishedAt     *time.Time
}
