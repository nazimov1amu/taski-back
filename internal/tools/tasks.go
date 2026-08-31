package tools

import (
	"context"
	"encoding/json"
	"taski_backend/internal/models"
	"taski_backend/internal/service"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

type TasksTools struct {
	TasksService *service.TasksService
}

type CreateTaskArgs struct {
	Title       string `json:"title" jsonschema:"required,Task title"`
	Description string `json:"description,omitempty" jsonschema:"Task description"`
	ProjectID   string `json:"project_id,omitempty" jsonschema:"Project id"`
	PlannedAt   string `json:"planned_at" jsonschema:"required,Date YYYY-MM-DD"`
	StartTime   string `json:"start_time" jsonschema:"required,RFC3339"`
	EndTime     string `json:"end_time" jsonschema:"required,RFC3339"`
}

type UpdateTaskArgs struct {
	ID          string `json:"id" jsonschema:"required,Task id"`
	Title       string `json:"title" jsonschema:"required,Task title"`
	Description string `json:"description,omitempty" jsonschema:"Task description"`
	ProjectID   string `json:"project_id,omitempty" jsonschema:"Project id"`
	PlannedAt   string `json:"planned_at" jsonschema:"required,Date YYYY-MM-DD"`
	StartTime   string `json:"start_time" jsonschema:"required,RFC3339"`
	EndTime     string `json:"end_time" jsonschema:"required,RFC3339"`
}

func NewTasksTools(tasksService *service.TasksService) *TasksTools {
	return &TasksTools{TasksService: tasksService}
}

func (t *TasksTools) CreateTaskHandler(ctx context.Context, req mcp.CallToolRequest, args CreateTaskArgs) (*mcp.CallToolResult, error) {
	startTime, err := time.Parse(time.RFC3339, args.StartTime)
	if err != nil {
		return mcp.NewToolResultError("invalid start_time: " + err.Error()), nil
	}

	endTime, err := time.Parse(time.RFC3339, args.EndTime)
	if err != nil {
		return mcp.NewToolResultError("invalid end_time: " + err.Error()), nil
	}

	_, err = t.TasksService.Create(ctx, models.CreateTaskRequest{
		Title:       args.Title,
		Description: args.Description,
		ProjectID:   args.ProjectID,
		PlannedAt:   args.PlannedAt,
		StartTime:   startTime,
		EndTime:     endTime,
	})
	if err != nil {
		return mcp.NewToolResultError("failed to create task: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Task created successfully"), nil
}

func (t *TasksTools) SearchTasksHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("title")
	if err != nil {
		return nil, err
	}
	tasks, err := t.TasksService.Search(ctx, query)
	if err != nil {
		return mcp.NewToolResultError("failed to search tasks: " + err.Error()), nil
	}

	result, err := json.Marshal(tasks)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal tasks: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func (t *TasksTools) UpdateTaskHandler(ctx context.Context, req mcp.CallToolRequest, args UpdateTaskArgs) (*mcp.CallToolResult, error) {
	startTime, err := time.Parse(time.RFC3339, args.StartTime)
	if err != nil {
		return mcp.NewToolResultError("invalid start_time: " + err.Error()), nil
	}

	endTime, err := time.Parse(time.RFC3339, args.EndTime)
	if err != nil {
		return mcp.NewToolResultError("invalid end_time: " + err.Error()), nil
	}

	task, err := t.TasksService.Update(ctx, models.UpdateTaskRequest{
		ID:          args.ID,
		Title:       args.Title,
		Description: args.Description,
		ProjectID:   args.ProjectID,
		PlannedAt:   args.PlannedAt,
		StartTime:   startTime,
		EndTime:     endTime,
	})
	
	if err != nil {
		return mcp.NewToolResultError("failed to update task: " + err.Error()), nil
	}
	return mcp.NewToolResultStructured(task, "Task updated successfully"), nil
}

func (t *TasksTools) DeleteTaskHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("invalid id: " + err.Error()), nil
	}
	err = t.TasksService.Delete(ctx, id)
	if err != nil {
		return mcp.NewToolResultError("failed to delete task: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Task deleted successfully"), nil
}

func (t *TasksTools) NewCreateTool(tasksService *service.TasksService) mcp.Tool {
	return mcp.NewTool(
		"create_task",
		mcp.WithDescription("Create a new task"),
		mcp.WithInputSchema[CreateTaskArgs](),
	)
}

func (t *TasksTools) NewSearchTool(tasksService *service.TasksService) mcp.Tool {
	return mcp.NewTool(
		"search_tasks",
		mcp.WithDescription("Search for tasks"),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the task")),
	)
}

func (t *TasksTools) NewUpdateTool(tasksService *service.TasksService) mcp.Tool {
	return mcp.NewTool(
		"update_task",
		mcp.WithDescription("Update a task"),
		mcp.WithInputSchema[UpdateTaskArgs](),
	)
}

func (t *TasksTools) NewDeleteTool(tasksService *service.TasksService) mcp.Tool {
	return mcp.NewTool(
		"delete_task",
		mcp.WithDescription("Delete a task"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id of the task")),
	)
}

