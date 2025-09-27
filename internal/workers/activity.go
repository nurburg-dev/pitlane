package workers

import "context"

type ActivityWorker interface {
	Execute(ctx context.Context)
}
