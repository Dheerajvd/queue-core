package core

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Dheerajvd/queue-core/types"
)

func FlushQueue(ctx context.Context, queueName string) error {
	return RedisClient.Del(ctx, queueName).Err()
}

func RemoveJobByID(ctx context.Context, queueName string, jobID string) error {
	jobs, err := RedisClient.LRange(ctx, queueName, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, raw := range jobs {
		var job types.Job
		err := json.Unmarshal([]byte(raw), &job)
		if err != nil {
			continue
		}
		if job.ID == jobID {
			return RedisClient.LRem(ctx, queueName, 1, raw).Err()
		}
	}

	return errors.New("job ID not found in queue")
}
