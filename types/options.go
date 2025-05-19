package types

type WorkerOptions struct {
    MaxRetries int
    RetryDelay Duration
    RetryFunc  func(Job)
    DLQFunc    func(Job)
    OnSuccess  func(Job)
}
