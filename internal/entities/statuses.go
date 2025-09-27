package entities

type ActivityRetryStatus struct {
	RetryCount int
}

type ActivityStatus string

const (
	ActivityStatusExecuting ActivityStatus = "executing"
	ActivityStatusFailed    ActivityStatus = "failed"
	ActivityStatusPending   ActivityStatus = "pending"
	ActivityStatusFinished  ActivityStatus = "finished"
)

type WorkflowStatus string

const (
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusExecuting WorkflowStatus = "executing"
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusFinished  WorkflowStatus = "finished"
	WorkflowStatusAborted   WorkflowStatus = "aborted"
)
