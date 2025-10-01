# Architecture Documentation

## Overview

The visio-mcp-server is designed as a Model Context Protocol (MCP) server that provides programmatic access to Microsoft Visio diagram files. The architecture follows a modular design with clear separation of concerns.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     MCP Client (Claude)                      │
└────────────────────────┬────────────────────────────────────┘
                         │ MCP Protocol (stdio)
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Node.js Launcher                           │
│  - Platform detection                                        │
│  - Binary selection                                          │
│  - Process spawning                                          │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                 Go MCP Server                                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Server Layer                            │   │
│  │  - Tool registration                                 │   │
│  │  - Request routing                                   │   │
│  │  - Response formatting                               │   │
│  └────────────┬─────────────────────────────────────────┘   │
│               │                                              │
│  ┌────────────▼─────────────────────────────────────────┐   │
│  │              Tools Layer                             │   │
│  │  - visio_describe_pages                              │   │
│  │  - visio_read_page                                   │   │
│  │  - visio_list_shapes                                 │   │
│  │  - visio_write_shape                                 │   │
│  └────────────┬─────────────────────────────────────────┘   │
│               │                                              │
│  ┌────────────▼─────────────────────────────────────────┐   │
│  │              Visio Layer                             │   │
│  │  - Reader: Parse VSDX files                          │   │
│  │  - Writer: Modify VSDX files                         │   │
│  │  - Models: Data structures                           │   │
│  └────────────┬─────────────────────────────────────────┘   │
└───────────────┼──────────────────────────────────────────────┘
                │
┌───────────────▼──────────────────────────────────────────────┐
│                    File System                               │
│  - VSDX files (ZIP containers)                               │
│  - XML content                                               │
│  - Embedded resources                                        │
└──────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. Node.js Launcher

**Location**: `launcher/launcher.ts`

**Responsibilities**:
- Detect the operating system and architecture
- Select the appropriate pre-built Go binary
- Spawn the binary as a child process
- Pass through stdin/stdout for MCP communication

**Why Node.js?**:
- Easy NPM distribution
- Cross-platform compatibility
- Simple installation via `npx`

### 2. MCP Server (Go)

**Location**: `internal/server/server.go`

**Responsibilities**:
- Initialize the MCP server
- Register available tools
- Handle MCP protocol communication
- Route requests to appropriate handlers

**Key Dependencies**:
- `github.com/mark3labs/mcp-go`: MCP protocol implementation

### 3. Tools Layer

**Location**: `internal/tools/handlers.go`

**Responsibilities**:
- Implement business logic for each tool
- Validate input arguments
- Call Visio layer functions
- Format responses as JSON

**Tools**:
1. **visio_describe_pages**: List all pages with metadata
2. **visio_read_page**: Read shapes from a specific page
3. **visio_list_shapes**: Get basic shape information
4. **visio_write_shape**: Create or modify shapes

### 4. Visio Layer

**Location**: `internal/visio/`

**Components**:

#### Reader (`reader.go`)
- Opens VSDX files as ZIP archives
- Parses XML content
- Extracts page and shape information
- Builds in-memory data structures

**Key Methods**:
- `ReadDocument()`: Read entire document
- `ListPages()`: Get page metadata
- `ReadPage(name)`: Read specific page

#### Writer (`writer.go`)
- Modifies existing VSDX files
- Creates new VSDX files
- Manipulates XML content
- Maintains file integrity

**Key Methods**:
- `WriteShape()`: Add/modify shapes
- `CreateNewDocument()`: Create new file

#### Models (`models.go`)
- Define data structures
- Document, Page, Shape types
- Properties and metadata

## Data Flow

### Read Operation Flow

```
1. MCP Client sends request
   ↓
2. Launcher spawns Go binary
   ↓
3. Server receives request
   ↓
4. Server routes to tool handler
   ↓
5. Handler calls Visio Reader
   ↓
6. Reader opens VSDX file (ZIP)
   ↓
7. Reader parses XML content
   ↓
8. Reader builds data structures
   ↓
9. Handler formats response
   ↓
10. Server sends response to client
```

### Write Operation Flow

```
1. MCP Client sends write request
   ↓
2. Server routes to write handler
   ↓
3. Handler validates input
   ↓
4. Handler calls Visio Writer
   ↓
5. Writer opens VSDX file
   ↓
6. Writer creates temp file
   ↓
7. Writer modifies XML content
   ↓
8. Writer writes to temp file
   ↓
9. Writer replaces original file
   ↓
10. Handler sends success response
```

## File Format Handling

### VSDX Structure

```
document.vsdx (ZIP)
├── [Content_Types].xml          # MIME types
├── _rels/
│   └── .rels                    # Package relationships
├── docProps/
│   ├── app.xml                  # Application properties
│   ├── core.xml                 # Core properties
│   └── custom.xml               # Custom properties
└── visio/
    ├── document.xml             # Document structure
    ├── pages/
    │   ├── page1.xml           # Page 1 content
    │   ├── page2.xml           # Page 2 content
    │   └── _rels/
    │       └── page1.xml.rels  # Page relationships
    ├── masters/                 # Master shapes
    └── windows.xml              # Window settings
```

### XML Parsing Strategy

**Current Implementation** (Simplified):
- String-based parsing
- Regular expressions for extraction
- Quick prototyping

**Production Recommendation**:
- Use `encoding/xml` package
- Define proper struct tags
- Marshal/Unmarshal with type safety

Example:
```go
type Shape struct {
    XMLName xml.Name `xml:"Shape"`
    ID      string   `xml:"ID,attr"`
    Cells   []Cell   `xml:"Cell"`
}

type Cell struct {
    Name  string  `xml:"N,attr"`
    Value float64 `xml:"V,attr"`
}
```

## Error Handling

### Strategy

1. **Validation Errors**: Return immediately with descriptive message
2. **File Errors**: Wrap with context (file path, operation)
3. **Parse Errors**: Include line/position information
4. **Write Errors**: Clean up temporary files

### Example

```go
if !visio.FileExists(filePath) {
    return nil, fmt.Errorf("file not found: %s", filePath)
}

page, err := reader.ReadPage(pageName)
if err != nil {
    return nil, fmt.Errorf("failed to read page %s: %w", pageName, err)
}
```

## Performance Considerations

### Current Limitations

1. **Full File Loading**: Entire VSDX loaded into memory
2. **String Parsing**: Inefficient XML parsing
3. **No Caching**: Re-parse on every request

### Optimization Opportunities

1. **Streaming**: Parse XML incrementally
2. **Caching**: Cache parsed documents
3. **Lazy Loading**: Load pages on demand
4. **Parallel Processing**: Parse multiple pages concurrently

## Security Considerations

### File Access

- Validate file paths (no directory traversal)
- Check file extensions
- Limit file sizes
- Sandbox file operations

### XML Parsing

- Prevent XML bomb attacks
- Limit entity expansion
- Validate XML structure
- Sanitize user input

## Testing Strategy

### Unit Tests

- Test individual functions
- Mock file system operations
- Test error conditions

### Integration Tests

- Test with real VSDX files
- Test read/write cycles
- Test edge cases

### End-to-End Tests

- Test MCP protocol
- Test with MCP client
- Test tool invocations

## Deployment

### Distribution

1. **Go Binaries**: Built for multiple platforms via GoReleaser
2. **NPM Package**: Contains launcher and binaries
3. **Installation**: Via `npx` or Smithery

### Platforms Supported

- Windows (x86, x64, ARM64)
- macOS (x64, ARM64)
- Linux (x86, x64, ARM64)

## Future Architecture Enhancements

### 1. Plugin System

Allow custom shape handlers and formatters.

### 2. Caching Layer

Cache parsed documents for better performance.

### 3. Streaming API

Support large files with streaming.

### 4. Master Shape Library

Pre-load common master shapes.

### 5. Validation Layer

Validate VSDX files before processing.

## References

- [MCP Protocol Specification](https://modelcontextprotocol.io/)
- [Visio File Format](https://learn.microsoft.com/en-us/office/client-developer/visio/introduction-to-the-visio-file-formatvsdx)
- [Office Open XML](https://www.ecma-international.org/publications-and-standards/standards/ecma-376/)
- [Go Archive/Zip Package](https://pkg.go.dev/archive/zip)
- [Go Encoding/XML Package](https://pkg.go.dev/encoding/xml)
