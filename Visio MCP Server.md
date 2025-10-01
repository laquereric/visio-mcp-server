# Visio MCP Server

A Model Context Protocol (MCP) server that reads and writes MS Visio diagram data.

![Visio MCP Server](https://img.shields.io/badge/MCP-Server-blue)
![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)

## Features

- **Read diagram data**: Extract shapes, text, and properties from Visio files
- **List pages**: Get information about all pages in a document
- **Read shapes**: Access shape properties including position, size, and text
- **Write shapes**: Create and modify shapes programmatically

## Requirements

- Node.js 20.x or later

## Supported File Formats

- `.vsdx` (Visio drawing)
- `.vsdm` (Visio macro-enabled drawing)
- `.vstx` (Visio template)
- `.vstm` (Visio macro-enabled template)

## Installation

### Installing via NPM

The visio-mcp-server is automatically installed by adding the following configuration to your MCP servers configuration.

**For Windows:**

```json
{
    "mcpServers": {
        "visio": {
            "command": "cmd",
            "args": ["/c", "npx", "--yes", "@negokaz/visio-mcp-server"],
            "env": {
                "VISIO_MCP_MAX_SHAPES": "1000"
            }
        }
    }
}
```

**For macOS/Linux:**

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

### Installing via Smithery

To install Visio MCP Server for Claude Desktop automatically via [Smithery](https://smithery.ai):

```shell
npx -y @smithery/cli install @negokaz/visio-mcp-server --client claude
```

## Tools

### `visio_describe_pages`

List all pages in a Visio file with metadata.

**Arguments:**

- `fileAbsolutePath` (string, required)
  - Absolute path to the Visio file

**Example Response:**

```json
{
  "file": "/path/to/diagram.vsdx",
  "pageCount": 3,
  "pages": [
    {
      "id": "page1.xml",
      "name": "Network Diagram",
      "width": 11.0,
      "height": 8.5,
      "shapeCount": 15,
      "background": ""
    }
  ]
}
```

### `visio_read_page`

Read shapes and their properties from a specific page.

**Arguments:**

- `fileAbsolutePath` (string, required)
  - Absolute path to the Visio file
- `pageName` (string, required)
  - Name of the page to read
- `includeConnections` (boolean, optional)
  - Include connector information [default: false]

**Example Response:**

```json
{
  "file": "/path/to/diagram.vsdx",
  "pageName": "Network Diagram",
  "width": 11.0,
  "height": 8.5,
  "shapeCount": 5,
  "shapes": [
    {
      "id": "1",
      "name": "Server",
      "text": "Web Server",
      "type": "Shape",
      "pinX": 5.5,
      "pinY": 4.25,
      "width": 2.0,
      "height": 1.5
    }
  ]
}
```

### `visio_list_shapes`

List all shapes on a page with basic information.

**Arguments:**

- `fileAbsolutePath` (string, required)
  - Absolute path to the Visio file
- `pageName` (string, required)
  - Name of the page

**Example Response:**

```json
{
  "file": "/path/to/diagram.vsdx",
  "pageName": "Network Diagram",
  "shapeCount": 5,
  "shapes": [
    {
      "id": "1",
      "text": "Web Server",
      "type": "Shape",
      "x": 5.5,
      "y": 4.25
    }
  ]
}
```

### `visio_write_shape`

Create or modify shapes on a page.

**Arguments:**

- `fileAbsolutePath` (string, required)
  - Absolute path to the Visio file
- `pageName` (string, required)
  - Target page name
- `shapeData` (object, required)
  - Shape properties:
    - `text` (string): Shape text content
    - `pinX` (number): X coordinate in inches
    - `pinY` (number): Y coordinate in inches
    - `width` (number): Shape width in inches
    - `height` (number): Shape height in inches
- `createPage` (boolean, optional)
  - Create page if it doesn't exist [default: false]

**Example Request:**

```json
{
  "fileAbsolutePath": "/path/to/diagram.vsdx",
  "pageName": "Network Diagram",
  "shapeData": {
    "text": "Database Server",
    "pinX": 8.0,
    "pinY": 4.0,
    "width": 2.0,
    "height": 1.5
  }
}
```

## Configuration

You can customize the MCP server behavior using environment variables:

### `VISIO_MCP_MAX_SHAPES`

The maximum number of shapes to read from a single page.  
**Default:** 1000

### `VISIO_MCP_INCLUDE_HIDDEN`

Include hidden shapes when reading pages.  
**Default:** false

## Development

### Prerequisites

- Go 1.23 or later
- Node.js 20.x or later
- TypeScript

### Building from Source

```bash
# Clone the repository
git clone https://github.com/negokaz/visio-mcp-server.git
cd visio-mcp-server

# Install dependencies
npm install

# Build the project
npm run build
```

### Running Tests

```bash
go test ./...
```

### Debugging

You can debug the MCP server using the MCP Inspector:

```bash
npm run debug
```

## Architecture

The Visio MCP Server is built with:

- **Go**: Core server and file manipulation logic
- **TypeScript/Node.js**: NPM package distribution and launcher
- **MCP Protocol**: Communication with AI assistants

### Project Structure

```
visio-mcp-server/
├── cmd/visio-mcp-server/    # Main entry point
├── internal/
│   ├── visio/               # Visio file handling
│   ├── server/              # MCP server implementation
│   └── tools/               # Tool handlers
├── launcher/                # Node.js launcher
└── docs/                    # Documentation
```

## Limitations

1. **File Format**: Works with VSDX (Office Open XML) format only. Legacy VSD format is not supported.
2. **Shape Creation**: Can create basic shapes but not complex master-based shapes.
3. **No Rendering**: Cannot render diagrams to images without Visio application.
4. **Stencils**: Stencil files (.vssx, .vssm) are not supported in the current version.
5. **Macros**: Cannot execute VBA macros in .vsdm files.

## Roadmap

- [ ] Support for stencil files
- [ ] Shape styling (colors, line styles, fills)
- [ ] Connector manipulation
- [ ] Layer management
- [ ] Data graphics and data linking
- [ ] Theme support
- [ ] Page management (add, delete, reorder)
- [ ] Shape search by properties

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Copyright (c) 2025 Kazuki Negoro

visio-mcp-server is released under the [MIT License](LICENSE).

## Acknowledgments

- Inspired by [excel-mcp-server](https://github.com/negokaz/excel-mcp-server)
- Built with [mcp-go](https://github.com/mark3labs/mcp-go)
- Uses the Model Context Protocol by Anthropic

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/negokaz/visio-mcp-server).

---

**Note:** This is a community project and is not officially affiliated with Microsoft or Visio.
