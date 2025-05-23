package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Dheerajvd/queue-core/core"
	"github.com/Dheerajvd/queue-core/types"
	"github.com/redis/go-redis/v9"
)

func StartDelayedJobScheduler(ctx context.Context, queueName string, interval time.Duration) {
	go func() {
		zsetKey := "delayed:" + queueName
		for {
			select {
			case <-ctx.Done():
				return
			default:
				now := time.Now().Unix()
				// Get all jobs that should be executed by now
				jobs, err := core.RedisClient.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
					Min: "-inf",
					Max: strconv.FormatInt(now, 10),
				}).Result()
				if err != nil {
					log.Printf("scheduler: error fetching delayed jobs: %v", err)
					time.Sleep(interval)
					continue
				}

				for _, jobData := range jobs {
					var job types.Job
					if err := json.Unmarshal([]byte(jobData), &job); err != nil {
						log.Printf("scheduler: failed to unmarshal delayed job: %v", err)
						continue
					}

					// Push to main queue
					if err := core.RedisClient.RPush(ctx, queueName, jobData).Err(); err != nil {
						log.Printf("scheduler: failed to push job to queue: %v", err)
						continue
					}

					// Remove from delayed zset
					core.RedisClient.ZRem(ctx, zsetKey, jobData)
				}

				time.Sleep(interval)
			}
		}
	}()
}
