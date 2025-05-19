package worker

import (
    "context"
    "time"
    "shared-queue/types"
)

type Worker struct {
    ClientID string
    Handler  func(context.Context, types.Job) error
    Options  types.WorkerOptions
}

func NewWorker(clientID string, handler func(context.Context, types.Job) error, opts types.WorkerOptions) *Worker {
    return &Worker{ClientID: clientID, Handler: handler, Options: opts}
}

func (w *Worker) Start(ctx context.Context, jobs <-chan types.Job) {
    for {
        select {
        case <-ctx.Done():
            return
        case job := <-jobs:
            err := w.Handler(ctx, job)
            if err != nil && job.RetryCount < w.Options.MaxRetries {
                time.Sleep(w.Options.RetryDelay)
                job.RetryCount++
                w.Options.RetryFunc(job)
            } else if err != nil {
                w.Options.DLQFunc(job)
            } else if w.Options.OnSuccess != nil {
                w.Options.OnSuccess(job)
            }
        }
    }
}
