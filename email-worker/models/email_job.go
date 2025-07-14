package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// StringArray represents a string array for database storage
type StringArray []string

// VariablesMap represents a variables map for database storage
type VariablesMap map[string]any

// JobStatus represents the status of an email job
type JobStatus string

// Job status constants
const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

// JobPriority represents the priority of an email job
type JobPriority int

// Job priority constants
const (
	JobPriorityHigh   JobPriority = 1
	JobPriorityNormal JobPriority = 2
	JobPriorityLow    JobPriority = 3
)

// EmailJob represents an email job in the system
type EmailJob struct {
	ID           int64        `db:"id" json:"-"`                    // Internal ID for performance
	PublicID     uuid.UUID    `db:"public_id" json:"id"`            // Public ID for API
	To           StringArray  `db:"to_emails" json:"to"`
	CC           StringArray  `db:"cc_emails" json:"cc"`
	BCC          StringArray  `db:"bcc_emails" json:"bcc"`
	TemplateName string       `db:"template_name" json:"template_name"`
	Variables    VariablesMap `db:"variables" json:"variables"`
	Status       JobStatus    `db:"status" json:"status"`
	Priority     JobPriority  `db:"priority" json:"priority"`
	RetryCount   int          `db:"retry_count" json:"retry_count"`
	MaxRetries   int          `db:"max_retries" json:"max_retries"`
	ErrorMessage string       `db:"error_message" json:"error_message"`
	ProcessedAt  *time.Time   `db:"processed_at" json:"processed_at"`
	SentAt       *time.Time   `db:"sent_at" json:"sent_at"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`

	// Queue-specific fields
	IsTracked    bool       `json:"is_tracked"`
	QueueID      string     `json:"queue_id"`
	ProcessingAt *time.Time `json:"processing_at"`
	CompletedAt  *time.Time `json:"completed_at"`
}

// Value implements driver.Valuer for StringArray
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	// Convert to PostgreSQL array format: {"item1","item2","item3"}
	if len(s) == 0 {
		return "{}", nil
	}

	// Escape quotes and format as PostgreSQL array
	result := "{"
	for i, item := range s {
		if i > 0 {
			result += ","
		}
		// Escape double quotes by doubling them
		escaped := strings.ReplaceAll(item, `"`, `""`)
		result += `"` + escaped + `"`
	}
	result += "}"

	return result, nil
}

// Scan implements sql.Scanner for StringArray
func (s *StringArray) Scan(value any) error {
	if value == nil {
		*s = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// Handle PostgreSQL array format: {"item1","item2","item3"}
		str := string(v)
		if str == "{}" {
			*s = StringArray{}
			return nil
		}

		// Remove outer braces
		if len(str) < 2 || str[0] != '{' || str[len(str)-1] != '}' {
			return errors.New("invalid PostgreSQL array format")
		}
		str = str[1 : len(str)-1]

		if str == "" {
			*s = StringArray{}
			return nil
		}

		// Parse array elements
		var result []string
		var current strings.Builder
		inQuotes := false
		escapeNext := false

		for i := 0; i < len(str); i++ {
			char := str[i]

			if escapeNext {
				current.WriteByte(char)
				escapeNext = false
				continue
			}

			if char == '\\' {
				escapeNext = true
				continue
			}

			if char == '"' {
				inQuotes = !inQuotes
				continue
			}

			if char == ',' && !inQuotes {
				result = append(result, current.String())
				current.Reset()
				continue
			}

			current.WriteByte(char)
		}

		// Add the last element
		if current.Len() > 0 || len(result) > 0 {
			result = append(result, current.String())
		}

		*s = StringArray(result)
		return nil

	case string:
		// Handle string representation
		return s.Scan([]byte(v))

	default:
		return errors.New("unsupported type for StringArray scan")
	}
}

// Value implements driver.Valuer for VariablesMap
func (m VariablesMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements sql.Scanner for VariablesMap
func (m *VariablesMap) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

// NewEmailJob tạo một email job mới
func NewEmailJob(to, cc, bcc []string, templateName string, variables map[string]any, priority JobPriority) *EmailJob {
	return &EmailJob{
		PublicID:     uuid.New(),
		To:           StringArray(to),
		CC:           StringArray(cc),
		BCC:          StringArray(bcc),
		TemplateName: templateName,
		Variables:    VariablesMap(variables),
		Status:       JobStatusPending,
		Priority:     priority,
		RetryCount:   0,
		MaxRetries:   3,
		ErrorMessage: "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// SetScheduledAt sets the scheduled time for the job
func (j *EmailJob) SetScheduledAt(scheduledAt time.Time) {
	j.ProcessedAt = &scheduledAt
}

// SetPriority sets the priority for the job
func (j *EmailJob) SetPriority(priority JobPriority) {
	j.Priority = priority
}

// SetMaxRetries sets the maximum number of retries
func (j *EmailJob) SetMaxRetries(maxRetries int) {
	j.MaxRetries = maxRetries
}

// SetQueueID sets the queue message ID
func (j *EmailJob) SetQueueID(queueID string) {
	j.QueueID = queueID
}

// CanRetry checks if the job can be retried
func (j *EmailJob) CanRetry() bool {
	return j.RetryCount < j.MaxRetries
}

// IncrementRetry increments the retry count
func (j *EmailJob) IncrementRetry() {
	j.RetryCount++
	j.UpdatedAt = time.Now()
}

// IsReadyToProcess checks if the job is ready to be processed
func (j *EmailJob) IsReadyToProcess() bool {
	if j.Status != JobStatusPending {
		return false
	}
	return true
}

// MarkAsProcessing marks the job as processing
func (j *EmailJob) MarkAsProcessing() {
	now := time.Now()
	j.Status = JobStatusProcessing
	j.ProcessingAt = &now
	j.UpdatedAt = now
}

// MarkAsCompleted marks the job as completed
func (j *EmailJob) MarkAsCompleted() {
	now := time.Now()
	j.Status = JobStatusCompleted
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// MarkAsFailed marks the job as failed
func (j *EmailJob) MarkAsFailed() {
	now := time.Now()
	j.Status = JobStatusFailed
	j.CompletedAt = &now
	j.UpdatedAt = now
}

// MarkAsRetrying marks the job as retrying
func (j *EmailJob) MarkAsRetrying() {
	j.Status = JobStatusProcessing // Use processing for retrying
	j.UpdatedAt = time.Now()
}

// IsCompleted checks if the job is completed (success or failure)
func (j *EmailJob) IsCompleted() bool {
	return j.Status == JobStatusCompleted || j.Status == JobStatusFailed
}

// GetProcessingDuration returns the processing duration if completed
func (j *EmailJob) GetProcessingDuration() *time.Duration {
	if j.ProcessingAt == nil || j.CompletedAt == nil {
		return nil
	}
	duration := j.CompletedAt.Sub(*j.ProcessingAt)
	return &duration
}

// ShouldBeTracked determines if this job should be tracked in database
func (j *EmailJob) ShouldBeTracked() bool {
	// Track high priority jobs
	if j.Priority >= 2 {
		return true
	}
	return j.IsTracked
}
