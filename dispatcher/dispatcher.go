package dispatcher

import (
    "context"
    "time"

    "github.com/google/uuid"
    "shared-queue/types"
)

type Dispatcher struct {
    EnqueueFunc func(job types.Job) error
}

func NewDispatcher(enqueueFunc func(job types.Job) error) *Dispatcher {
    return &Dispatcher{EnqueueFunc: enqueueFunc}
}

func (d *Dispatcher) Enqueue(ctx context.Context, jobType, clientID string, payload any, opts ...types.JobOption) error {
    job := types.NewJob(uuid.NewString(), jobType, clientID, payload)
    for _, opt := range opts {
        opt(&job)
    }
    return d.EnqueueFunc(job)
}
