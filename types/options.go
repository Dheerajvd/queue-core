package types

import "time"

type WorkerOptions struct {
	MaxRetries int
	RetryDelay time.Duration
	RetryFunc  func(Job)
	DLQFunc    func(Job)
	OnSuccess  func(Job)
}
