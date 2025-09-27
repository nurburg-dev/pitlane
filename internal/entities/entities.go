package entities

import (
	"encoding/json"
	"time"
)

type DBWorkflow struct {
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type DBWorkflowRun struct {
	ID           string          `json:"id" db:"id"`
	Input        json.RawMessage `json:"input" db:"input"`
	WorkflowName string          `json:"workflow_name" db:"workflow_name"`
	Status       WorkflowStatus  `json:"status" db:"status"`
	ScheduledAt  time.Time       `json:"scheduled_at" db:"scheduled_at"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

type DBActivityRun struct {
	ID            string           `json:"id" db:"id"`
	ActivityName  string           `json:"activity_name" db:"activity_name"`
	WorkflowRunID string           `json:"workflow_run_id" db:"workflow_run_id"`
	ErrorMessage  *string          `json:"error_message" db:"errorMessage"`
	Input         json.RawMessage  `json:"input" db:"input"`
	Output        *json.RawMessage `json:"output" db:"output"`
	Status        ActivityStatus   `json:"status" db:"status"`
	RetryStatus   *json.RawMessage `json:"retry_status" db:"retry_status"`
	ScheduledAt   time.Time        `json:"scheduled_at" db:"scheduled_at"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
}
