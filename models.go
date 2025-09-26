package pitlane

type ActivityRetryStatus struct {
	RetryCount int
}

type ActivityStatus string

const (
	ActivityStatusExecuting ActivityStatus = "executing"
	ActivityStatusFailed    ActivityStatus = "failed"
	ActivityStatusWaiting   ActivityStatus = "waiting"
	ActivityStatusFinished  ActivityStatus = "finished"
)

type WorkflowStatus string

const (
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusExecuting WorkflowStatus = "executing"
	WorkflowStatusFinished  WorkflowStatus = "finished"
	WorkflowStatusAborted   WorkflowStatus = "aborted"
)
