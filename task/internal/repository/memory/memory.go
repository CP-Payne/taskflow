package memory

import (
	"context"
	"sync"

	"github.com/CP-Payne/taskflow/task/internal/model"
	"github.com/CP-Payne/taskflow/task/internal/repository"
	"github.com/google/uuid"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	task map[uuid.UUID]*model.Task
}

func NewInMemory() *MemoryRepository {
	return &MemoryRepository{
		task: make(map[uuid.UUID]*model.Task),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.task[task.ID] = task
	return r.task[task.ID], nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// taskList := make([]model.Task, len(r.task))
	taskList := []model.Task{}
	for _, v := range r.task {
		taskList = append(taskList, *v)
	}

	return taskList, nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if v, ok := r.task[taskID]; ok {
		return v, nil
	}
	return nil, repository.ErrNotFound
}

func (r *MemoryRepository) ListByAssignedUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskList := []model.Task{}
	for _, v := range r.task {
		if v.AssignedTo != nil && *v.AssignedTo == userID {
			taskList = append(taskList, *v)
		}
	}
	return taskList, nil
}

func (r *MemoryRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	taskList := []model.Task{}
	for _, v := range r.task {
		if v.UserID == userID {
			taskList = append(taskList, *v)
		}
	}

	return taskList, nil
}

func (r *MemoryRepository) ListUnassigned(ctx context.Context) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	taskList := []model.Task{}
	for _, v := range r.task {
		if v.AssignedTo == nil {
			taskList = append(taskList, *v)
		}
	}
	return taskList, nil
}
