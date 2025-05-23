package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Dheerajvd/queue-core/core"
	"github.com/Dheerajvd/queue-core/types"
)

type Worker struct {
	QueueName           string
	ProcessingQueueName string
	Handler             func(context.Context, types.Job) error
	Options             types.WorkerOptions
}

func NewWorker(queueName string, handler func(context.Context, types.Job) error, opts types.WorkerOptions) *Worker {
	return &Worker{
		QueueName:           queueName,
		ProcessingQueueName: "processing:" + queueName,
		Handler:             handler,
		Options:             opts,
	}
}

func (w *Worker) Start(ctx context.Context, jobs <-chan types.Job) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-jobs:
			err := w.Handler(ctx, job)
			if err != nil {
				log.Printf("worker: job %s failed: %v", job.ID, err)
				if job.RetryCount < w.Options.MaxRetries {
					time.Sleep(w.Options.RetryDelay)
					job.RetryCount++
					if w.Options.RetryFunc != nil {
						w.Options.RetryFunc(job)
					}
					// Re-queue job for retry
					data, _ := json.Marshal(job)
					core.RedisClient.RPush(ctx, w.QueueName, data)
					// Remove from processing queue
					core.RedisClient.LRem(ctx, w.ProcessingQueueName, 1, data)
				} else {
					// Move to DLQ
					if w.Options.DLQFunc != nil {
						w.Options.DLQFunc(job)
					}
					// Remove from processing queue
					data, _ := json.Marshal(job)
					core.RedisClient.LRem(ctx, w.ProcessingQueueName, 1, data)
					// Push to DLQ queue
					dlqName := "dlq:" + w.QueueName
					core.RedisClient.RPush(ctx, dlqName, data)
				}
			} else {
				if w.Options.OnSuccess != nil {
					w.Options.OnSuccess(job)
				}
				// Remove job from processing queue on success
				data, _ := json.Marshal(job)
				core.RedisClient.LRem(ctx, w.ProcessingQueueName, 1, data)
			}
		}
	}
}
