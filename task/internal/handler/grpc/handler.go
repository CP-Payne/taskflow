package grpc

import (
	"context"
	"errors"
	"time"

	api "github.com/CP-Payne/taskflow/pkg/gen/task/v1"
	"github.com/CP-Payne/taskflow/task/internal/model"
	"github.com/CP-Payne/taskflow/task/internal/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	api.UnimplementedTaskServiceServer
	taskService *service.TaskService
	logger      *zap.SugaredLogger
}

func NewTaskHandler(taskService *service.TaskService, logger *zap.SugaredLogger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		logger:      logger,
	}
}

func (h *TaskHandler) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	if req == nil || req.Title == "" || req.UserId == nil {
		h.logger.Warnw("Create validation failed: invalid arguments",
			"userID", req.UserId.GetValue(),
			"title", req.GetTitle(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}

	task := &model.Task{
		ID:          uuid.New(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Status:      model.Pending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	userID, err := uuid.Parse(req.GetUserId().GetValue())
	if err != nil {
		h.logger.Errorw("Create: invalid arguements",
			"userID", req.GetUserId().GetValue(),
			"error", zap.Error(err),
		)

		return nil, status.Errorf(codes.InvalidArgument, "failed to create task: invalid userID")
	}

	task.UserID = userID

	assignedToStr := req.GetAssignedTo().GetValue()
	if assignedToStr != "" {
		assignedToID, err := uuid.Parse(assignedToStr)
		if err != nil {
			h.logger.Warn("Create: Task assignment failed",
				"userID", req.GetUserId().GetValue(),
				"assignedToID", assignedToStr,
				"error", err,
			)
		} else {
			task.AssignedTo = &assignedToID
		}
	}

	task, err = h.taskService.CreateTask(ctx, task)
	if err != nil {
		h.logger.Errorw("Internal error during task creation",
			"userID", req.GetUserId(),
			"TaskID", task.ID,
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal server error")

	}

	h.logger.Infow("Task created successfully",
		"taskID", task.ID,
		"userID", task.UserID,
	)

	return &api.CreateResponse{Task: task.ToProto()}, nil
}

func (h *TaskHandler) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	if req == nil {
		h.logger.Warnw("List request is nil")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}
	tasks, err := h.taskService.ListAll(ctx)
	if err != nil {
		h.logger.Errorw("Internal error during tasks retrieval",
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	h.logger.Infow("Tasks listed successfully", "count", len(tasks))

	return &api.ListResponse{Tasks: model.TaskListToProto(tasks)}, nil
}

func (h *TaskHandler) ListUnassigned(ctx context.Context, req *api.ListUnassignedRequest) (*api.ListUnassignedResponse, error) {
	if req == nil {
		h.logger.Warnw("ListUnassigned request is nil")
		return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}
	tasks, err := h.taskService.ListAllUnassigned(ctx)
	if err != nil {
		h.logger.Errorw("Internal error during unassigned tasks retrieval",
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if tasks == nil {
		tasks = []model.Task{}
	}
	h.logger.Infow("Unassigned tasks listed successfully", "count", len(tasks))

	return &api.ListUnassignedResponse{Tasks: model.TaskListToProto(tasks)}, nil
}

func (h *TaskHandler) GetByID(ctx context.Context, req *api.GetByIDRequest) (*api.GetByIDResponse, error) {
	if req == nil || req.GetTaskId() == nil {
		h.logger.Warnw("GetByID validation failed: invalid arguments",
			"taskID", req.GetTaskId(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}

	taskID, err := uuid.Parse(req.GetTaskId().GetValue())
	if err != nil {
		h.logger.Warnw("GetByID invalid taskID",
			"taskID", req.GetTaskId().GetValue(),
			"error", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid arguments")

	}

	task, err := h.taskService.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.logger.Warnw("GetByID task not found",
				"taskID", taskID.String(),
			)
			return nil, status.Errorf(codes.NotFound, "resource not found")
		}
		h.logger.Errorw("GetByID internal error",
			"taskID", taskID.String(),
			"error", err,
		)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &api.GetByIDResponse{Task: task.ToProto()}, nil
}

func (h *TaskHandler) ListByAssignedUserID(ctx context.Context, req *api.ListByAssignedUserIDRequest) (*api.ListByAssignedUserIDResponse, error) {
	if req == nil || req.GetUserId() == nil {
		h.logger.Warnw("ListByAssignedUserID validation failed: invalid arguments",
			"userID", req.GetUserId(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}

	userID, err := uuid.Parse(req.GetUserId().GetValue())
	if err != nil {
		h.logger.Warnw("ListByAssignedUserID invalid userID",
			"userID", req.GetUserId().GetValue(),
			"error", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid userID")
	}

	tasks, err := h.taskService.ListByAssignedUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.logger.Warnw("ListByAssignedUserID user not found",
				"userID", userID.String(),
			)
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		h.logger.Errorw("ListByAssignedUserID internal error",
			"userID", userID.String(),
			"error", err,
		)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	h.logger.Infow("Tasks assigned to user listed successfully", "count", len(tasks), "userID", userID)

	return &api.ListByAssignedUserIDResponse{Tasks: model.TaskListToProto(tasks)}, nil
}

func (h *TaskHandler) ListByUserID(ctx context.Context, req *api.ListByUserIDRequest) (*api.ListByUserIDResponse, error) {
	if req == nil || req.GetUserId() == nil {
		h.logger.Warnw("ListByUserID validation failed: invalid arguments",
			"userID", req.GetUserId(),
		)
		return nil, status.Errorf(codes.InvalidArgument, "nil request or invalid arguments")
	}

	userID, err := uuid.Parse(req.GetUserId().GetValue())
	if err != nil {
		h.logger.Warnw("ListByUserID invalid userID",
			"userID", req.GetUserId().GetValue(),
			"error", err,
		)
		return nil, status.Errorf(codes.InvalidArgument, "invalid userID")
	}

	tasks, err := h.taskService.ListByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			h.logger.Warnw("ListByUserID user not found",
				"userID", userID.String(),
			)
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		h.logger.Errorw("ListByUserID internal error",
			"userID", userID.String(),
			"error", err,
		)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	h.logger.Infow("Tasks created by user listed successfully", "count", len(tasks), "userID", userID)
	return &api.ListByUserIDResponse{Tasks: model.TaskListToProto(tasks)}, nil
}
