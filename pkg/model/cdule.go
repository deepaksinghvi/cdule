package model

import (
	"gorm.io/gorm"
	"time"
)

// JobStatus for job status
type JobStatus string

const (
	// JobStatusNew status NEW
	JobStatusNew JobStatus = "NEW"
	// JobStatusInProgress status IN_PROGRESS
	JobStatusInProgress JobStatus = "IN_PROGRESS"
	// JobStatusCompleted status COMPLETED
	JobStatusCompleted JobStatus = "COMPLETED"
	// JobStatusFailed status FAILED
	JobStatusFailed JobStatus = "FAILED"
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Job struct {
	Model
	JobName        string `json:"job_name"`
	GroupName      string `json:"group_name"`
	CronExpression string `json:"cron"`
	Expired        bool   `json:"expired"`
	JobData        string `json:"job_data"`
}

type JobHistory struct {
	Model
	JobID       int64          `json:"job_id"`
	ExecutionID int64          `json:"execution_id"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Status      JobStatus      `json:"status"`
	WorkerID    string         `json:"worker_id"`
	RetryCount  int            `json:"retry_count"`
}

// Schedule used by Execution Routine to execute a scheduled job in the evert one minute duration
type Schedule struct {
	ExecutionID int64 `gorm:"primaryKey",json:"execution_id"`
	JobID       int64 `gorm:"primaryKey",json:"job_id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	WorkerID    string         `json:"worker_id"`
	JobData     string         `json:"job_data"`
}

// Worker Node health check via the heartbeat
type Worker struct {
	WorkerID  string `gorm:"primaryKey",json:"worker_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
