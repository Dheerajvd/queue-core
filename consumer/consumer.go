package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Dheerajvd/queue-core/core"
	"github.com/Dheerajvd/queue-core/types"
	"github.com/redis/go-redis/v9"
)

// ConsumeJobs consumes jobs reliably from queueName.
// processingQueueName is usually "processing:" + queueName
// out channel will receive ready jobs
func ConsumeJobs(ctx context.Context, queueName, processingQueueName string, out chan<- types.Job) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Atomically move job from main queue to processing queue
				result, err := core.RedisClient.BRPopLPush(ctx, queueName, processingQueueName, 5*time.Second).Result()
				if err != nil {
					continue
				}

				var job types.Job
				err = json.Unmarshal([]byte(result), &job)
				if err != nil {
					log.Printf("consumer: failed to unmarshal job: %v", err)
					continue
				}

				// Check scheduled time
				if !job.Scheduled.IsZero() && job.Scheduled.After(time.Now()) {
					// Re-queue with delay
					// Push back to delayed set for scheduling
					go func(j types.Job) {
						// If scheduled for future, push to delayed set
						zsetKey := "delayed:" + queueName
						data, _ := json.Marshal(j)
						score := float64(j.Scheduled.Unix())
						core.RedisClient.ZAdd(ctx, zsetKey, redis.Z{
							Score:  score,
							Member: data,
						})
						// Remove from processing queue since not ready
						core.RedisClient.LRem(ctx, processingQueueName, 1, result)
					}(job)
					continue
				}

				out <- job
			}
		}
	}()
}
