package tools

import (
	"context"
	"encoding/json"
	"log"

	"taski_backend/internal/models"
	"taski_backend/internal/service"

	"github.com/mark3labs/mcp-go/mcp"
)

type ProjectsTools struct {
	ProjectsService *service.ProjectsService
}

type CreateProjectArgs struct {
	Name        string `json:"name" jsonschema:"required,Project name"`
	Description string `json:"description,omitempty" jsonschema:"Project description"`
}

type UpdateProjectArgs struct {
	ID          string `json:"id" jsonschema:"required,Project id"`
	Name        string `json:"name" jsonschema:"required,Project name"`
	Description string `json:"description,omitempty" jsonschema:"Project description"`
}

func NewProjectsTools(projectsService *service.ProjectsService) *ProjectsTools {
	return &ProjectsTools{ProjectsService: projectsService}
}

func (t *ProjectsTools) CreateProjectHandler(ctx context.Context, req mcp.CallToolRequest, args CreateProjectArgs) (*mcp.CallToolResult, error) {
	project, err := t.ProjectsService.Create(ctx, models.CreateProjectRequest{
		Name:        args.Name,
		Description: args.Description,
	})
	log.Println("created project", project)
	if err != nil {
		log.Println("error creating project:", err)
		return mcp.NewToolResultError("failed to create project: " + err.Error()), nil
	}
	return mcp.NewToolResultStructured(project, "Project created successfully"), nil
}

func (t *ProjectsTools) SearchProjectsHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var filters models.ProjectFilters
	if name, err := req.RequireString("name"); err == nil {
		filters.Name = &name
	}

	projects, err := t.ProjectsService.GetBulk(ctx, filters)
	if err != nil {
		return mcp.NewToolResultError("failed to search projects: " + err.Error()), nil
	}
	log.Println("projects", projects)
	result, err := json.Marshal(projects)
	if err != nil {
		return mcp.NewToolResultError("failed to marshal projects: " + err.Error()), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

func (t *ProjectsTools) UpdateProjectHandler(ctx context.Context, req mcp.CallToolRequest, args UpdateProjectArgs) (*mcp.CallToolResult, error) {
	project, err := t.ProjectsService.Update(ctx, args.ID, models.UpdateProjectRequest{
		Name:        args.Name,
		Description: args.Description,
	})
	log.Println("updated project", project)
	if err != nil {
		return mcp.NewToolResultError("failed to update project: " + err.Error()), nil
	}
	return mcp.NewToolResultStructured(project, "Project updated successfully"), nil
}

func (t *ProjectsTools) DeleteProjectHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError("invalid id: " + err.Error()), nil
	}
	log.Println("deleting project", id)
	err = t.ProjectsService.Delete(ctx, id)
	if err != nil {
		return mcp.NewToolResultError("failed to delete project: " + err.Error()), nil
	}
	return mcp.NewToolResultText("Project deleted successfully"), nil
}

func (t *ProjectsTools) NewCreateTool() mcp.Tool {
	return mcp.NewTool(
		"create_project",
		mcp.WithDescription("Create a new project"),
		mcp.WithInputSchema[CreateProjectArgs](),
	)
}

func (t *ProjectsTools) NewSearchTool() mcp.Tool {
	return mcp.NewTool(
		"search_projects",
		mcp.WithDescription("Search for projects by name or description"),
		mcp.WithString("name", mcp.Description("Project name or description fragment to search")),
	)
}

func (t *ProjectsTools) NewUpdateTool() mcp.Tool {
	return mcp.NewTool(
		"update_project",
		mcp.WithDescription("Update a project"),
		mcp.WithInputSchema[UpdateProjectArgs](),
	)
}

func (t *ProjectsTools) NewDeleteTool() mcp.Tool {
	return mcp.NewTool(
		"delete_project",
		mcp.WithDescription("Delete a project"),
		mcp.WithString("id", mcp.Required(), mcp.Description("The id of the project")),
	)
}
