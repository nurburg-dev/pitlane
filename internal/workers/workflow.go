package workers

import "context"

type WorkflowWorker interface {
	Execute(ctx context.Context)
}
