package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/x3a-tech/logit-go"
	"i-o-bouns-tasks-api/internal/models/request"
	"i-o-bouns-tasks-api/internal/models/response"
	"i-o-bouns-tasks-api/internal/models/status"
	"i-o-bouns-tasks-api/internal/models/tasks"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Tasks struct {
	logger    logit.Logger
	tasks     map[int32]*tasks.Task
	mutex     sync.RWMutex
	idCounter int32
}

type TasksParams struct {
	Logger logit.Logger
}

func NewTasks(params *TasksParams) Repository {
	return &Tasks{
		logger:    params.Logger,
		tasks:     make(map[int32]*tasks.Task),
		idCounter: 0,
	}
}

func (t *Tasks) CreateTask(w http.ResponseWriter, r *http.Request) {
	const op = "app.Init.api.Init.Tasks.CreateTask"
	ctx := t.logger.NewOpCtx(r.Context(), op)

	t.logger.Info(ctx, "Создание новой задачи...")

	var taskReq request.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		t.logger.Error(ctx, fmt.Errorf("ошибка при разборе тела ответа: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := atomic.AddInt32(&t.idCounter, 1)
	t.logger.Info(ctx, fmt.Sprintf("Id задачи: %v", id))

	task := &tasks.Task{
		Id:          id,
		Status:      status.StatusPending,
		CreatedAt:   time.Now(),
		Description: taskReq.Description,
	}

	t.mutex.Lock()
	t.tasks[id] = task
	t.mutex.Unlock()

	task.Status = status.StatusProcessing

	time.Sleep(1 * time.Minute)

	task.Status = status.StatusCompleted

	resp := &response.TaskResponse{
		Id:        task.Id,
		Status:    task.Status,
		CreatedAt: task.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		t.logger.Error(ctx, fmt.Errorf("ошибка при кодировании ответа: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t.logger.Info(ctx, fmt.Sprintf("Задача успешно создана. Id: %v", task.Id))
}

func (t *Tasks) GetTask(w http.ResponseWriter, r *http.Request) {
	const op = "app.Init.api.Init.Tasks.GetTask"
	ctx := t.logger.NewOpCtx(r.Context(), op)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		t.logger.Error(ctx, fmt.Errorf("ошибка при получении ID из запроса: %v", err))
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	t.mutex.RLock()
	task, exists := t.tasks[int32(id)]
	t.mutex.RUnlock()

	if !exists {
		t.logger.Warn(ctx, fmt.Sprintf("Задача с ID %d не найдена", id))
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	executionTime := time.Since(task.CreatedAt).Seconds()

	resp := &response.TaskResponse{
		Id:            task.Id,
		Status:        task.Status,
		CreatedAt:     task.CreatedAt,
		ExecutionTime: executionTime,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		t.logger.Error(ctx, fmt.Errorf("ошибка при кодировании ответа: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t.logger.Info(ctx, fmt.Sprintf("Информация о задаче отправлена. Id: %v, Время выполнения: %.2f секунд", task.Id, executionTime))

}

func (t *Tasks) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "app.Init.api.Init.Tasks.DeleteTask"
	ctx := t.logger.NewOpCtx(r.Context(), op)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		t.logger.Error(ctx, fmt.Errorf("ошибка при получении ID из запроса: %v", err))
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	task, exists := t.tasks[int32(id)]
	if !exists {
		t.logger.Warn(ctx, fmt.Sprintf("Задача с Id %d не найдена", id))
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if task.Status != status.StatusCompleted {
		t.logger.Warn(ctx, fmt.Sprintf("Попытка удалить незавершенную задачу с ID %d", id))
		http.Error(w, "Cannot delete uncompleted task", http.StatusBadRequest)
		return
	}

	delete(t.tasks, int32(id))

	resp := &response.TaskResponse{
		Message: fmt.Sprintf("Задача Id %v удалена", id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		t.logger.Error(ctx, fmt.Errorf("ошибка при кодировании ответа: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	t.logger.Info(ctx, fmt.Sprintf("Задача успешно удалена. Id: %v", id))

}
