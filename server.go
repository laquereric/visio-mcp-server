package server

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/negokaz/visio-mcp-server/internal/tools"
)

// Server represents the MCP server
type Server struct {
	version string
	mcp     *server.MCPServer
}

// New creates a new server instance
func New(version string) *Server {
	return &Server{
		version: version,
	}
}

// Start starts the MCP server
func (s *Server) Start() error {
	// Create MCP server
	s.mcp = server.NewMCPServer(
		"visio-mcp-server",
		s.version,
		server.WithStdio(),
	)

	// Register tools
	s.registerTools()

	// Start server
	return s.mcp.Serve()
}

// registerTools registers all available tools
func (s *Server) registerTools() {
	// Describe pages tool
	s.mcp.AddTool(mcp.Tool{
		Name:        "visio_describe_pages",
		Description: "List all pages in a Visio file with metadata",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"fileAbsolutePath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the Visio file",
				},
			},
			Required: []string{"fileAbsolutePath"},
		},
	}, tools.DescribePagesHandler)

	// Read page tool
	s.mcp.AddTool(mcp.Tool{
		Name:        "visio_read_page",
		Description: "Read shapes and their properties from a specific page",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"fileAbsolutePath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the Visio file",
				},
				"pageName": map[string]interface{}{
					"type":        "string",
					"description": "Name of the page to read",
				},
				"includeConnections": map[string]interface{}{
					"type":        "boolean",
					"description": "Include connector information",
					"default":     false,
				},
			},
			Required: []string{"fileAbsolutePath", "pageName"},
		},
	}, tools.ReadPageHandler)

	// List shapes tool
	s.mcp.AddTool(mcp.Tool{
		Name:        "visio_list_shapes",
		Description: "List all shapes on a page with basic information",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"fileAbsolutePath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the Visio file",
				},
				"pageName": map[string]interface{}{
					"type":        "string",
					"description": "Name of the page",
				},
			},
			Required: []string{"fileAbsolutePath", "pageName"},
		},
	}, tools.ListShapesHandler)

	// Write shape tool
	s.mcp.AddTool(mcp.Tool{
		Name:        "visio_write_shape",
		Description: "Create or modify shapes on a page",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"fileAbsolutePath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the Visio file",
				},
				"pageName": map[string]interface{}{
					"type":        "string",
					"description": "Target page name",
				},
				"shapeData": map[string]interface{}{
					"type":        "object",
					"description": "Shape properties (text, position, size, type)",
					"properties": map[string]interface{}{
						"text": map[string]interface{}{
							"type":        "string",
							"description": "Shape text content",
						},
						"pinX": map[string]interface{}{
							"type":        "number",
							"description": "X coordinate (in inches)",
						},
						"pinY": map[string]interface{}{
							"type":        "number",
							"description": "Y coordinate (in inches)",
						},
						"width": map[string]interface{}{
							"type":        "number",
							"description": "Shape width (in inches)",
						},
						"height": map[string]interface{}{
							"type":        "number",
							"description": "Shape height (in inches)",
						},
					},
				},
				"createPage": map[string]interface{}{
					"type":        "boolean",
					"description": "Create page if it doesn't exist",
					"default":     false,
				},
			},
			Required: []string{"fileAbsolutePath", "pageName", "shapeData"},
		},
	}, tools.WriteShapeHandler)

	fmt.Fprintf(os.Stderr, "Registered %d tools\n", 4)
}
