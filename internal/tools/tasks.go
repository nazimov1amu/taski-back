package tools

import (
	"context"
	"encoding/json"
	"log"
	"taski_backend/internal/models"
	"taski_backend/internal/service"

	"github.com/mark3labs/mcp-go/mcp"
)

type TasksTools struct {
	TasksService *service.TasksService
}

type CreateTaskArgs struct {
	Title       string `json:"title" jsonschema:"required,Task title"`
	Description string `json:"description,omitempty" jsonschema:"Task description"`
	ProjectID   string `json:"project_id,omitempty" jsonschema:"Project id"`
	PlannedAt   string `json:"planned_at,omitempty" jsonschema:"optional,Local date YYYY-MM-DD"`
	StartTime   string `json:"start_time,omitempty" jsonschema:"optional,Local datetime YYYY-MM-DDTHH:MM:SS"`
	EndTime     string `json:"end_time,omitempty" jsonschema:"optional,Local datetime YYYY-MM-DDTHH:MM:SS"`
}

type UpdateTaskArgs struct {
	ID          string `json:"id" jsonschema:"required,Task id"`
	Title       string `json:"title,omitempty" jsonschema:"Task title"`
	Description string `json:"description,omitempty" jsonschema:"Task description"`
	ProjectID   string `json:"project_id,omitempty" jsonschema:"Project id"`
	PlannedAt   string `json:"planned_at,omitempty" jsonschema:"optional,Local date YYYY-MM-DD"`
	StartTime   string `json:"start_time,omitempty" jsonschema:"optional,Local datetime YYYY-MM-DDTHH:MM:SS"`
	EndTime     string `json:"end_time,omitempty" jsonschema:"optional,Local datetime YYYY-MM-DDTHH:MM:SS"`
}

func NewTasksTools(tasksService *service.TasksService) *TasksTools {
	return &TasksTools{TasksService: tasksService}
}

func (t *TasksTools) CreateTaskHandler(ctx context.Context, req mcp.CallToolRequest, args CreateTaskArgs) (*mcp.CallToolResult, error) {
	created, err := t.TasksService.Create(ctx, models.CreateTaskRequest{
		TaskFields: models.TaskFields{
			Title:       args.Title,
			Description: args.Description,
			ProjectID:   args.ProjectID,
			PlannedAt:   args.PlannedAt,
			StartTime:   args.StartTime,
			EndTime:     args.EndTime,
		},
	})
	if err != nil {
		return mcp.NewToolResultError("failed to create task: " + err.Error()), nil
	}
	return mcp.NewToolResultStructured(created, "Task created successfully"), nil
}

func (t *TasksTools) SearchTasksHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var filters models.TaskFilters
	if title, err := req.RequireString("title"); err == nil {
		filters.Title = &title
	}
	log.Println("filters", filters)
	tasks, err := t.TasksService.GetBulk(ctx, filters)
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
	log.Println("args", args)
	task, err := t.TasksService.Update(ctx, args.ID, models.UpdateTaskRequest{
		Title:       args.Title,
		Description: args.Description,
		ProjectID:   args.ProjectID,
		PlannedAt:   args.PlannedAt,
		StartTime:   args.StartTime,
		EndTime:     args.EndTime,
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
	log.Println("deleting task", id)
	err = t.TasksService.Delete(ctx, id)
	if err != nil {
		return mcp.NewToolResultError("failed to delete task: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Task deleted successfully"), nil
}

func (t *TasksTools) CompleteTaskHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("invalid id: " + err.Error()), nil
	}
	log.Println("completing task", id)
	err = t.TasksService.Complete(ctx, id)
	if err != nil {
		return mcp.NewToolResultError("failed to complete task: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Task completed successfully"), nil
}

func (t *TasksTools) NewCreateTool() mcp.Tool {
	return mcp.NewTool(
		"create_task",
		mcp.WithDescription("Create a new task"),
		mcp.WithInputSchema[CreateTaskArgs](),
	)
}

func (t *TasksTools) NewSearchTool() mcp.Tool {
	return mcp.NewTool(
		"search_tasks",
		mcp.WithDescription("Search for tasks"),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the task")),
	)
}

func (t *TasksTools) NewUpdateTool() mcp.Tool {
	return mcp.NewTool(
		"update_task",
		mcp.WithDescription("Update a task"),
		mcp.WithInputSchema[UpdateTaskArgs](),
	)
}

func (t *TasksTools) NewDeleteTool() mcp.Tool {
	return mcp.NewTool(
		"delete_task",
		mcp.WithDescription("Delete a task"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id of the task")),
	)
}


func (t *TasksTools) NewCompleteTool() mcp.Tool {
	return mcp.NewTool(
		"complete_task",
		mcp.WithDescription("Complete a task"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id of the task to complete")),
	)
}