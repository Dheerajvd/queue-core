package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Dheerajvd/queue-core/core"
	"github.com/Dheerajvd/queue-core/types"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func PushJob(ctx context.Context, queueName, jobType, clientID string, payload any, opts ...types.JobOption) (string, error) {
	job := types.NewJob(uuid.NewString(), jobType, clientID, payload)
	for _, opt := range opts {
		opt(&job)
	}

	data, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	// If unique key specified, check existence
	if job.UniqueKey != "" {
		exists, err := core.RedisClient.Exists(ctx, job.UniqueKey).Result()
		if err != nil {
			return "", err
		}
		if exists > 0 {
			return "", errors.New("duplicate job with same unique key")
		}
		err = core.RedisClient.Set(ctx, job.UniqueKey, job.ID, 24*time.Hour).Err()
		if err != nil {
			return "", err
		}
	}

	// If delayed, add to sorted set with score = scheduled timestamp
	if !job.Scheduled.IsZero() && job.Scheduled.After(time.Now()) {
		zsetKey := "delayed:" + queueName
		score := float64(job.Scheduled.Unix())
		err := core.RedisClient.ZAdd(ctx, zsetKey, redis.Z{
			Score:  score,
			Member: data,
		}).Err()
		return job.ID, err
	}

	// Push to main queue list
	zsetKey := "processing:" + queueName
	return job.ID, core.RedisClient.RPush(ctx, zsetKey, data).Err()
}
