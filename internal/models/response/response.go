package response

import (
	"i-o-bouns-tasks-api/internal/models/status"
	"time"
)

type TaskResponse struct {
	Id            int32             `json:"id"`
	Status        status.TaskStatus `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	ExecutionTime float64           `json:"execution_time"`
	Message       string            `json:"message"`
}
