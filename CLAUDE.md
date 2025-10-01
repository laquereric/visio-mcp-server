# Claude Code Setup

This document provides guidance for working with the visio-mcp-server codebase using Claude Code.

## Project Overview

The visio-mcp-server is a Model Context Protocol (MCP) server that enables reading and writing Microsoft Visio diagram files programmatically. It's built with Go for core functionality and distributed via NPM using a Node.js launcher.

## Architecture

### Key Components

1. **cmd/visio-mcp-server/main.go**: Entry point that initializes and starts the server
2. **internal/server/server.go**: MCP server implementation and tool registration
3. **internal/visio/**: Visio file handling (reader, writer, models)
4. **internal/tools/**: Tool handlers for MCP operations
5. **launcher/launcher.ts**: Node.js launcher for cross-platform binary execution

### Technology Stack

- **Go 1.23+**: Core server and file manipulation
- **TypeScript**: NPM package launcher
- **MCP Protocol**: github.com/mark3labs/mcp-go
- **File Format**: VSDX (Office Open XML based)

## Development Workflow

### Prerequisites

```bash
# Install Go 1.23 or later
# Install Node.js 20.x or later
# Install TypeScript
npm install -g typescript
```

### Building

```bash
# Install Node dependencies
npm install

# Build Go binaries and TypeScript
npm run build

# Or build separately
go build ./cmd/visio-mcp-server
tsc
```

### Testing

```bash
# Run Go tests
go test ./...

# Run specific package tests
go test ./internal/visio

# Run with coverage
go test -cover ./...
```

### Debugging

```bash
# Debug with MCP Inspector
npm run debug

# Run directly
go run ./cmd/visio-mcp-server
```

## Code Organization

### internal/visio Package

Handles all Visio file operations:

- **models.go**: Data structures for documents, pages, and shapes
- **reader.go**: Reading VSDX files (ZIP + XML parsing)
- **writer.go**: Writing/modifying VSDX files

### internal/tools Package

Implements MCP tool handlers:

- **handlers.go**: Tool implementations for describe_pages, read_page, list_shapes, write_shape

### internal/server Package

MCP server setup and tool registration.

## VSDX File Format

VSDX files are ZIP containers with XML content:

```
document.vsdx (ZIP)
├── [Content_Types].xml
├── _rels/
├── docProps/
└── visio/
    ├── document.xml
    └── pages/
        └── page1.xml
```

Key XML elements:
- `<Shape>`: Represents a shape
- `<Cell N="PinX">`: Shape properties
- `<Text>`: Shape text content

## Common Tasks

### Adding a New Tool

1. Define tool in `internal/server/server.go`:
   ```go
   s.mcp.AddTool(mcp.Tool{
       Name: "visio_new_tool",
       Description: "...",
       InputSchema: mcp.ToolInputSchema{...},
   }, tools.NewToolHandler)
   ```

2. Implement handler in `internal/tools/handlers.go`:
   ```go
   func NewToolHandler(arguments map[string]interface{}) (*string, error) {
       // Implementation
   }
   ```

### Modifying Visio File Handling

1. Update models in `internal/visio/models.go`
2. Implement reading logic in `internal/visio/reader.go`
3. Implement writing logic in `internal/visio/writer.go`

### Testing File Operations

Create test VSDX files or use existing ones:

```go
func TestReadPage(t *testing.T) {
    reader := visio.NewReader("testdata/sample.vsdx")
    page, err := reader.ReadPage("Page-1")
    // Assertions
}
```

## Best Practices

1. **Error Handling**: Always return descriptive errors with context
2. **XML Parsing**: Use encoding/xml for production code, not string manipulation
3. **File Safety**: Use temporary files when modifying VSDX files
4. **Testing**: Add tests for new features
5. **Documentation**: Update README.md for user-facing changes

## Known Limitations

1. Simplified XML parsing (should use proper encoding/xml)
2. Limited shape creation (basic shapes only)
3. No stencil support
4. No image rendering
5. No connector manipulation

## Future Enhancements

- Proper XML marshaling/unmarshaling
- Master shape support
- Connector handling
- Style and formatting
- Layer management
- Data graphics

## Resources

- [MCP Protocol](https://modelcontextprotocol.io/)
- [Visio File Format](https://learn.microsoft.com/en-us/office/client-developer/visio/introduction-to-the-visio-file-formatvsdx)
- [Office Open XML](https://learn.microsoft.com/en-us/office/open-xml/open-xml-sdk)
- [mcp-go Library](https://github.com/mark3labs/mcp-go)
