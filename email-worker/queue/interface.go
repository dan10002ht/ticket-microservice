package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"booking-system/email-worker/models"
)

// Queue defines the interface for email job queues
type Queue interface {
	// Publish adds an email job to the queue
	Publish(ctx context.Context, job *models.EmailJob) error

	// Consume retrieves and removes the next job from the queue
	Consume(ctx context.Context) (*models.EmailJob, error)

	// ConsumeBatch retrieves multiple jobs from the queue
	ConsumeBatch(ctx context.Context, batchSize int) ([]*models.EmailJob, error)

	// Size returns the current queue size
	Size(ctx context.Context) (int64, error)

	// Clear removes all jobs from the queue
	Clear(ctx context.Context) error

	// Health checks if the queue is healthy
	Health(ctx context.Context) error

	// Close closes the queue connection
	Close() error

	// PublishScheduled publishes a job for scheduled delivery
	PublishScheduled(ctx context.Context, job *models.EmailJob, scheduledAt time.Time) error

	// ProcessScheduledJobs moves ready scheduled jobs to the main queue
	ProcessScheduledJobs(ctx context.Context) error
}

// QueueConfig holds configuration for queue implementations
type QueueConfig struct {
	Type         string `mapstructure:"type"` // redis, kafka, etc.
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database"`
	QueueName    string `mapstructure:"queue_name"`
	BatchSize    int    `mapstructure:"batch_size"`
	PollInterval string `mapstructure:"poll_interval"`
}

// QueueFactory creates queue instances based on configuration
type QueueFactory struct {
	logger *zap.Logger
}

// NewQueueFactory creates a new queue factory
func NewQueueFactory(logger *zap.Logger) *QueueFactory {
	return &QueueFactory{
		logger: logger,
	}
}

// CreateQueue creates a queue instance based on configuration
func (f *QueueFactory) CreateQueue(config QueueConfig) (Queue, error) {
	switch config.Type {
	case "redis":
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		return NewRedisQueue(addr, config.Password, config.Database, config.QueueName, f.logger), nil
	case "kafka":
		// TODO: Implement Kafka queue
		return nil, fmt.Errorf("kafka queue not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported queue type: %s", config.Type)
	}
}

// RedisQueue implements the Queue interface for Redis
type RedisQueue struct {
	addr      string
	password  string
	database  int
	queueName string
	logger    *zap.Logger
	client    *redis.Client
}

// NewRedisQueue creates a new RedisQueue instance
func NewRedisQueue(addr, password string, database int, queueName string, logger *zap.Logger) *RedisQueue {
	logger.Info("Creating RedisQueue", zap.String("addr", addr), zap.String("password", password), zap.Int("database", database), zap.String("queueName", queueName))
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})
	return &RedisQueue{
		addr:      addr,
		password:  password,
		database:  database,
		queueName: queueName,
		logger:    logger,
		client:    client,
	}
}

// Publish adds an email job to the queue
func (q *RedisQueue) Publish(ctx context.Context, job *models.EmailJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return q.client.LPush(ctx, q.queueName, data).Err()
}

// Consume retrieves and removes the next job from the queue
func (q *RedisQueue) Consume(ctx context.Context) (*models.EmailJob, error) {
	// Implementation of Consume method
	return nil, nil
}

// ConsumeBatch retrieves multiple jobs from the queue
func (q *RedisQueue) ConsumeBatch(ctx context.Context, batchSize int) ([]*models.EmailJob, error) {
	var jobs []*models.EmailJob
	for i := 0; i < batchSize; i++ {
		res, err := q.client.RPop(ctx, q.queueName).Result()
		if err == redis.Nil {
			break // Queue empty
		}
		if err != nil {
			return jobs, err
		}
		var job models.EmailJob
		if err := json.Unmarshal([]byte(res), &job); err != nil {
			q.logger.Error("Failed to unmarshal job from Redis", zap.Error(err))
			continue
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// Size returns the current queue size
func (q *RedisQueue) Size(ctx context.Context) (int64, error) {
	return q.client.LLen(ctx, q.queueName).Result()
}

// Clear removes all jobs from the queue
func (q *RedisQueue) Clear(ctx context.Context) error {
	return q.client.Del(ctx, q.queueName).Err()
}

// Health checks if the queue is healthy
func (q *RedisQueue) Health(ctx context.Context) error {
	return q.client.Ping(ctx).Err()
}

// Close closes the queue connection
func (q *RedisQueue) Close() error {
	return q.client.Close()
}

// PublishScheduled publishes a job for scheduled delivery
func (q *RedisQueue) PublishScheduled(ctx context.Context, job *models.EmailJob, scheduledAt time.Time) error {
	// For simplicity, push to a scheduled list with score = scheduledAt.Unix()
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}
	key := q.queueName + ":scheduled"
	return q.client.ZAdd(ctx, key, redis.Z{Score: float64(scheduledAt.Unix()), Member: data}).Err()
}

// ProcessScheduledJobs moves ready scheduled jobs to the main queue
func (q *RedisQueue) ProcessScheduledJobs(ctx context.Context) error {
	key := q.queueName + ":scheduled"
	now := time.Now().Unix()
	// Get jobs with score <= now
	jobs, err := q.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%d", now),
	}).Result()
	if err != nil {
		return err
	}
	if len(jobs) == 0 {
		return nil
	}
	// Move jobs to main queue and remove from scheduled
	for _, jobData := range jobs {
		if err := q.client.LPush(ctx, q.queueName, jobData).Err(); err != nil {
			q.logger.Error("Failed to move scheduled job to main queue", zap.Error(err))
			continue
		}
		q.client.ZRem(ctx, key, jobData)
	}
	return nil
}

// KafkaQueue implements the Queue interface for Kafka
type KafkaQueue struct {
	// Kafka-specific fields
}

// NewKafkaQueue creates a new KafkaQueue instance
func NewKafkaQueue(brokers []string, topic string, logger *zap.Logger) *KafkaQueue {
	// Kafka-specific initialization
	return &KafkaQueue{}
}

// Publish adds an email job to the queue
func (q *KafkaQueue) Publish(ctx context.Context, job *models.EmailJob) error {
	// Implementation of Publish method
	return nil
}

// Consume retrieves and removes the next job from the queue
func (q *KafkaQueue) Consume(ctx context.Context) (*models.EmailJob, error) {
	// Implementation of Consume method
	return nil, nil
}

// ConsumeBatch retrieves multiple jobs from the queue
func (q *KafkaQueue) ConsumeBatch(ctx context.Context, batchSize int) ([]*models.EmailJob, error) {
	// Implementation of ConsumeBatch method
	return nil, nil
}

// Size returns the current queue size
func (q *KafkaQueue) Size(ctx context.Context) (int64, error) {
	// Implementation of Size method
	return 0, nil
}

// Clear removes all jobs from the queue
func (q *KafkaQueue) Clear(ctx context.Context) error {
	// Implementation of Clear method
	return nil
}

// Health checks if the queue is healthy
func (q *KafkaQueue) Health(ctx context.Context) error {
	// Implementation of Health method
	return nil
}

// Close closes the queue connection
func (q *KafkaQueue) Close() error {
	// Implementation of Close method
	return nil
}

// PublishScheduled publishes a job for scheduled delivery
func (q *KafkaQueue) PublishScheduled(ctx context.Context, job *models.EmailJob, scheduledAt time.Time) error {
	// Implementation of PublishScheduled method
	return nil
}

// ProcessScheduledJobs moves ready scheduled jobs to the main queue
func (q *KafkaQueue) ProcessScheduledJobs(ctx context.Context) error {
	// Implementation of ProcessScheduledJobs method
	return nil
}

// Queue errors
var (
	ErrQueueEmpty = fmt.Errorf("queue is empty")
)
