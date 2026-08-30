package app

import (
	"database/sql"
	"fmt"
	"log"

	"taski_backend/internal/config"
	"taski_backend/internal/db"
	"taski_backend/internal/repository"
	"taski_backend/internal/service"
	"taski_backend/internal/tools"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCP struct {
	config     *config.AppConfig
	mcpServer  *server.MCPServer
	httpServer *server.StreamableHTTPServer
	db         *sql.DB
}

func NewMCP() (*MCP, error) {
	cfg, err := config.NewAppConfig()
	if err != nil {
		return nil, err
	}
	if cfg.Address == "" {
		cfg.Address = ":8081"
	}

	sqlDB, err := db.NewDB(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	tasksService := service.NewTasksService(repository.NewTasksRepository(sqlDB))
	tasksTools := tools.NewTasksTools(tasksService)

	mcpServer := server.NewMCPServer(
		"Taski MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)
	mcpServer.AddTool(tasksTools.NewCreateTool(tasksService), mcp.NewTypedToolHandler(tasksTools.CreateTaskHandler))
	mcpServer.AddTool(tasksTools.NewSearchTool(tasksService), tasksTools.SearchTasksHandler)
	mcpServer.AddTool(tasksTools.NewUpdateTool(tasksService), mcp.NewTypedToolHandler(tasksTools.UpdateTaskHandler))
	mcpServer.AddTool(tasksTools.NewDeleteTool(tasksService), tasksTools.DeleteTaskHandler)

	httpServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithEndpointPath("/mcp"),
	)

	return &MCP{
		config:     cfg,
		mcpServer:  mcpServer,
		httpServer: httpServer,
		db:         sqlDB,
	}, nil
}

func (m *MCP) Addr() string {
	return m.config.Address
}

func (m *MCP) Run() error {
	log.Printf("MCP Streamable HTTP on %s/mcp", m.config.Address)
	if err := m.httpServer.Start(m.config.Address); err != nil {
		return fmt.Errorf("serve mcp http: %w", err)
	}
	return nil
}

func (m *MCP) Close() error {
	if m.db == nil {
		return nil
	}
	return m.db.Close()
}
