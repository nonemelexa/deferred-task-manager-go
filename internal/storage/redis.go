package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"task-scheduler/internal/domain"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStorage establishes a connection to Redis and returns a RedisStorage struct for managing delayed tasks and failed tasks using Redis data structures.
func NewRedisStorage() *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	return &RedisStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Save adds a task to the Redis sorted set "delayed_tasks" with a score based on the execution time and priority, allowing for efficient retrieval of ready tasks.
func (r *RedisStorage) Save(task domain.Task) error {
	data, _ := json.Marshal(task)

	score := float64(task.ExecuteAt.Unix()) - float64(task.Priority)

	return r.client.ZAdd(r.ctx, "delayed_tasks", redis.Z{
		Score:  score,
		Member: data,
	}).Err()
}

// GetReadyTasks retrieves tasks that are ready to be executed based on the current time by querying the Redis sorted set "delayed_tasks" for members with scores less than or equal to the current time, unmarshaling them from JSON, and returning them as a slice. It also removes the retrieved tasks from the sorted set to prevent them from being processed again.
func (r *RedisStorage) GetReadyTasks(now time.Time) ([]domain.Task, error) {
	maxScore := float64(now.Unix())

	result, err := r.client.ZRangeByScore(r.ctx, "delayed_tasks", &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", maxScore),
	}).Result()

	if err != nil {
		return nil, err
	}

	var tasks []domain.Task

	for _, item := range result {
		var task domain.Task
		json.Unmarshal([]byte(item), &task)
		tasks = append(tasks, task)
	}

	if len(result) > 0 {
		r.client.ZRemRangeByScore(r.ctx, "delayed_tasks", "-inf", fmt.Sprintf("%f", maxScore))
	}

	return tasks, nil
}

// SaveFailed adds a task to the Redis list "failed_tasks" by marshaling it to JSON and pushing it to the list, allowing for later retrieval and analysis of failed tasks.
func (r *RedisStorage) SaveFailed(task domain.Task) error {
	data, _ := json.Marshal(task)
	return r.client.LPush(r.ctx, "failed_tasks", data).Err()
}
