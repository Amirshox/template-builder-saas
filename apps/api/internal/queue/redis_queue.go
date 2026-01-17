package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type JobPayload struct {
	JobID      uuid.UUID       `json:"jobId"`
	OrgID      uuid.UUID       `json:"orgId"`
	TemplateID uuid.UUID       `json:"templateId"`
	Data       json.RawMessage `json:"data"`
}

type Queue struct {
	client *redis.Client
	stream string
}

func NewQueue(redisAddr string, password string) *Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return &Queue{
		client: rdb,
		stream: "generation_jobs",
	}
}

func (q *Queue) EnqueueJob(ctx context.Context, job JobPayload) error {
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.stream,
		Values: map[string]interface{}{
			"payload": bytes,
		},
	}).Err()
}

// For Worker
func (q *Queue) Consume(ctx context.Context, group string, consumer string, handler func(JobPayload) error) {
	// Create group if not exists
	q.client.XGroupCreateMkStream(ctx, q.stream, group, "0")

	for {
		streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{q.stream, ">"},
			Count:    1,
			Block:    2 * time.Second,
		}).Result()

		if err != nil {
			if err != redis.Nil {
				fmt.Printf("Redis Read Error: %v\n", err)
			}
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				payloadStr := message.Values["payload"].(string)
				var job JobPayload
				if err := json.Unmarshal([]byte(payloadStr), &job); err != nil {
					fmt.Println("Invalid payload", err)
					q.client.XAck(ctx, q.stream, group, message.ID)
					continue
				}

				if err := handler(job); err == nil {
					q.client.XAck(ctx, q.stream, group, message.ID)
				} else {
					fmt.Printf("Job Failed: %v\n", err)
					// Handle DLQ or retry later
				}
			}
		}
	}
}
