package types

import "time"

type Job struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	ClientID   string    `json:"client_id"`
	Payload    any       `json:"payload"`
	RetryCount int       `json:"retry_count"`
	Scheduled  time.Time `json:"scheduled"`
	UniqueKey  string    `json:"unique_key"`
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
