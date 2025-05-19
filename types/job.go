package types

import "time"

type Job struct {
    ID         string
    Type       string
    ClientID   string
    Payload    any
    RetryCount int
    Scheduled  time.Time
    UniqueKey  string
}

type JobOption func(*Job)

func WithDelay(t time.Time) JobOption {
    return func(j *Job) {
        j.Scheduled = t
    }
}

func WithUniqueKey(key string) JobOption {
    return func(j *Job) {
        j.UniqueKey = key
    }
}

func NewJob(id, jobType, clientID string, payload any) Job {
    return Job{
        ID:       id,
        Type:     jobType,
        ClientID: clientID,
        Payload:  payload,
    }
}
