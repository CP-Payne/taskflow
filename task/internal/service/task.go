package service

import (
	"context"
	"errors"

	"github.com/CP-Payne/taskflow/task/internal/model"
	"github.com/CP-Payne/taskflow/task/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrInternal = errors.New("internal server error")
)

type TaskService struct {
	repo   repository.TaskRepository
	logger *zap.SugaredLogger
}

func New(repo repository.TaskRepository, logger *zap.SugaredLogger) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	task, err := s.repo.Create(ctx, task)
	if err != nil {
		return &model.Task{}, ErrInternal
	}
	return task, nil
}

func (s *TaskService) ListAll(ctx context.Context) ([]model.Task, error) {
	tasks, err := s.repo.List(ctx)
	if err != nil {
		return []model.Task{}, err
	}

	return tasks, nil
}

func (s *TaskService) ListAllUnassigned(ctx context.Context) ([]model.Task, error) {
	tasks, err := s.repo.ListUnassigned(ctx)
	if err != nil {
		return []model.Task{}, ErrInternal
	}
	return tasks, nil
}

func (s *TaskService) GetByID(ctx context.Context, taskID uuid.UUID) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &model.Task{}, ErrNotFound
		}
		return &model.Task{}, ErrInternal
	}
	return task, nil
}

func (s *TaskService) ListByAssignedUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	userTasks, err := s.repo.ListByAssignedUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return userTasks, nil
}

func (s *TaskService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	userTasks, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return userTasks, nil
}
