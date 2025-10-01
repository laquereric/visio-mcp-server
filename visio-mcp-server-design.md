# Visio MCP Server Design Document

## Overview

The **visio-mcp-server** is a Model Context Protocol (MCP) server that enables reading and writing Microsoft Visio diagram files programmatically. It follows the same architectural patterns as the excel-mcp-server but is adapted for Visio's VSDX file format.

## Architecture

### Technology Stack
- **Primary Language**: Go (for core server and file manipulation)
- **Runtime Distribution**: Node.js (for NPM package distribution)
- **File Format**: VSDX (Office Open XML based ZIP container)
- **MCP Protocol**: github.com/mark3labs/mcp-go
- **File Handling**: Standard library (archive/zip, encoding/xml)

### Project Structure

```
visio-mcp-server/
├── cmd/
│   └── visio-mcp-server/
│       └── main.go                 # Entry point
├── internal/
│   ├── visio/                      # Visio file manipulation
│   │   ├── reader.go              # Read VSDX files
│   │   ├── writer.go              # Write VSDX files
│   │   └── models.go              # Data structures
│   ├── mcp/                        # MCP protocol handling
│   │   └── server.go
│   ├── server/                     # Server implementation
│   │   └── server.go
│   └── tools/                      # MCP tool definitions
│       ├── describe_pages.go
│       ├── read_page.go
│       ├── write_shape.go
│       └── export_image.go
├── launcher/
│   └── launcher.ts                 # Node.js launcher
├── .goreleaser.yaml                # Build configuration
├── Dockerfile
├── go.mod
├── go.sum
├── package.json
├── tsconfig.json
├── LICENSE
└── README.md
```

## Core Features

### 1. Describe Pages
**Tool**: `visio_describe_pages`

Lists all pages in a Visio file with metadata:
- Page names
- Page dimensions
- Number of shapes
- Background page references

**Arguments**:
- `fileAbsolutePath`: Absolute path to VSDX file

### 2. Read Page
**Tool**: `visio_read_page`

Reads shapes and their properties from a specific page:
- Shape IDs
- Shape text
- Shape positions (X, Y coordinates)
- Shape dimensions (Width, Height)
- Shape types
- Connection information

**Arguments**:
- `fileAbsolutePath`: Absolute path to VSDX file
- `pageName`: Name of the page to read
- `includeConnections`: Include connector information (default: false)

### 3. Write Shape
**Tool**: `visio_write_shape`

Creates or modifies shapes on a page:
- Add new shapes
- Update shape text
- Modify shape position
- Change shape dimensions
- Set shape properties

**Arguments**:
- `fileAbsolutePath`: Absolute path to VSDX file
- `pageName`: Target page name
- `shapeData`: Shape properties (text, position, size, type)
- `createPage`: Create page if it doesn't exist (default: false)

### 4. List Shapes
**Tool**: `visio_list_shapes`

Lists all shapes on a page with basic information:
- Shape ID
- Shape name/text
- Shape type
- Layer information

**Arguments**:
- `fileAbsolutePath`: Absolute path to VSDX file
- `pageName`: Page name

### 5. Export Page Image
**Tool**: `visio_export_page_image`

Exports a page as an image (if supported):
- Extract embedded images
- Return base64 encoded image data

**Arguments**:
- `fileAbsolutePath`: Absolute path to VSDX file
- `pageName`: Page name
- `format`: Image format (png, jpg)

## Technical Implementation

### VSDX File Structure

The VSDX format is a ZIP container with the following key parts:

```
document.vsdx (ZIP)
├── [Content_Types].xml          # Content type definitions
├── _rels/
│   └── .rels                    # Package relationships
├── docProps/
│   ├── app.xml                  # Application properties
│   ├── core.xml                 # Core properties
│   └── custom.xml               # Custom properties
└── visio/
    ├── document.xml             # Document structure
    ├── pages/
    │   ├── page1.xml           # Page content
    │   ├── page2.xml
    │   └── _rels/
    │       ├── page1.xml.rels
    │       └── page2.xml.rels
    ├── masters/                 # Master shapes
    ├── windows.xml              # Window settings
    └── _rels/
        └── document.xml.rels
```

### Key Go Packages

```go
import (
    "archive/zip"
    "encoding/xml"
    "io"
    "os"
    "path/filepath"
)
```

### Data Models

```go
// Page represents a Visio page
type Page struct {
    ID          string
    Name        string
    Width       float64
    Height      float64
    Shapes      []Shape
    Background  string
}

// Shape represents a Visio shape
type Shape struct {
    ID          string
    Name        string
    Text        string
    Type        string
    PinX        float64  // X coordinate
    PinY        float64  // Y coordinate
    Width       float64
    Height      float64
    Master      string
    Properties  map[string]string
}

// Document represents a Visio document
type Document struct {
    Pages       []Page
    Masters     []Master
    Properties  DocumentProperties
}
```

## Configuration

Environment variables:
- `VISIO_MCP_MAX_SHAPES`: Maximum shapes to read per page (default: 1000)
- `VISIO_MCP_INCLUDE_HIDDEN`: Include hidden shapes (default: false)

## Supported File Formats

- .vsdx (Visio drawing)
- .vsdm (Visio macro-enabled drawing)
- .vstx (Visio template)
- .vstm (Visio macro-enabled template)

Note: .vssx and .vssm (stencil files) are not supported in initial version.

## Installation

### Via NPM

```json
{
    "mcpServers": {
        "visio": {
            "command": "npx",
            "args": ["--yes", "@negokaz/visio-mcp-server"],
            "env": {
                "VISIO_MCP_MAX_SHAPES": "1000"
            }
        }
    }
}
```

### Via Smithery

```bash
npx -y @smithery/cli install @negokaz/visio-mcp-server --client claude
```

## Limitations

1. **No Windows COM Automation**: Unlike Excel server which supports live editing on Windows, Visio server works purely with file manipulation
2. **Limited Shape Creation**: Can create basic shapes but not complex master-based shapes initially
3. **No Macro Execution**: Cannot execute VBA macros in .vsdm files
4. **Read-Only Stencils**: Stencil files (.vssx, .vssm) are not supported
5. **No Rendering**: Cannot render diagrams to images (requires Visio application or third-party library)

## Future Enhancements

1. Support for stencil files
2. Shape styling (colors, line styles, fills)
3. Connector manipulation (connect shapes)
4. Layer management
5. Data graphics and data linking
6. Theme support
7. Page ordering and management
8. Shape search by properties

## License

MIT License (same as excel-mcp-server)
