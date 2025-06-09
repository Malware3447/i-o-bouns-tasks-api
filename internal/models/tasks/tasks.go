package tasks

import (
	"i-o-bouns-tasks-api/internal/models/status"
	"time"
)

type Task struct {
	Id          int32
	Status      status.TaskStatus
	CreatedAt   time.Time
	Description string
}
