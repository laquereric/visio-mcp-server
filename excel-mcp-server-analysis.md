# Excel MCP Server Analysis

## Repository Structure

### Key Directories:
- `.devcontainer/` - Development container configuration
- `.github/workflows/` - CI/CD workflows
- `cmd/excel-mcp-server/` - Main command/entry point
- `docs/` - Documentation
- `internal/` - Internal Go packages
- `launcher/` - Launcher scripts

### Key Files:
- `.dockerignore` - Docker ignore patterns
- `.editorconfig` - Editor configuration
- `.gitignore` - Git ignore patterns
- `.goreleaser.yaml` - GoReleaser configuration
- `.npmignore` - NPM ignore patterns
- `CLAUDE.md` - Claude Code setup
- `Dockerfile` - Docker container definition
- `LICENSE` - MIT License
- `README.md` - Main documentation

## Technology Stack:
- **Primary Language**: Go (93.5%)
- **Secondary**: Shell (3.5%), Dockerfile (2.1%), JavaScript (0.9%)
- **Runtime**: Node.js 20.x or later
- **Distribution**: NPM package (@negokaz/excel-mcp-server)

## Features:
1. Read/Write text values
2. Read/Write formulas
3. Create new sheets
4. Windows-specific: Live editing, Screen capture

## Tools Provided:
1. `excel_describe_sheets` - List sheet information
2. `excel_read_sheet` - Read values with pagination
3. `excel_screen_capture` - Screenshot (Windows only)
4. `excel_write_to_sheet` - Write values
5. `excel_create_table` - Create tables
6. `excel_copy_sheet` - Copy sheets
7. `excel_format_range` - Format cells

## Configuration:
- Environment variable: `EXCEL_MCP_PAGING_CELLS_LIMIT` (default: 4000)

## Supported Formats:
- xlsx, xlsm, xltx, xltm

## package.json Content:

```json
{
  "name": "@negokaz/excel-mcp-server",
  "version": "0.12.0",
  "description": "An MCP server that reads and writes spreadsheet data to MS Excel file",
  "author": "negokaz",
  "license": "MIT",
  "bin": {
    "excel-mcp-server": "dist/launcher.js"
  },
  "scripts": {
    "build": "goreleaser build --snapshot --clean && tsc",
    "watch": "tsc --watch",
    "debug": "npx @modelcontextprotocol/inspector dist/launcher.js"
  },
  "devDependencies": {
    "@types/node": "^22.13.4",
    "typescript": "^5.7.3"
  },
  "publishConfig": {
    "access": "public"
  }
}
```

## cmd/excel-mcp-server/main.go Content:

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/negokaz/excel-mcp-server/internal/server"
)

var (
    version = "dev"
)

func main() {
    s := server.New(version)
    err := s.Start()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to start the server: %v\n", err)
        os.Exit(1)
    }
}
```

**Key observations:**
- Very simple main entry point
- Imports internal server package
- Creates server instance with version
- Starts the server
- Handles errors with proper exit codes

## Internal Directory Structure:

The internal directory contains four subdirectories:
1. **excel/** - Excel file handling logic
2. **mcp/** - MCP protocol implementation
3. **server/** - Server implementation
4. **tools/** - Tool definitions and implementations

## launcher/launcher.ts Content:

```typescript
#!/usr/bin/env node
import * as path from 'path'
import * as childProcess from 'child_process'

const BINARY_DISTRIBUTION_PACKAGES: any = {
    win32_ia32: "excel-mcp-server_windows_386_sse2",
    win32_x64: "excel-mcp-server_windows_amd64_v1",
    win32_arm64: "excel-mcp-server_windows_arm64_v8.0",
    darwin_x64: "excel-mcp-server_darwin_amd64_v1",
    darwin_arm64: "excel-mcp-server_darwin_arm64_v8.0",
    linux_ia32: "excel-mcp-server_linux_386_sse2",
    linux_x64: "excel-mcp-server_linux_amd64_v1",
    linux_arm64: "excel-mcp-server_linux_arm64_v8.0",
}

function getBinaryPath(): string {
    const suffix = process.platform === 'win32' ? '.exe' : '';
    const pkg = BINARY_DISTRIBUTION_PACKAGES[ `${process.platform}_${process.arch}` ];
    if (!pkg) {
        return path.resolve(__dirname, pkg, `excel-mcp-server${suffix}`);
    } else {
        throw new Error(`Unsupported platform: ${process.platform}_${process.arch}`);
    }
}

childProcess.execFileSync(getBinaryPath(), process.argv, {
    stdio: 'inherit',
});
```

**Key observations:**
- Node.js launcher script that wraps the Go binary
- Detects platform and architecture to select correct binary
- Supports multiple platforms: Windows (win32), macOS (darwin), Linux
- Supports multiple architectures: ia32, x64, arm64
- Uses execFileSync to run the binary with inherited stdio
- Binaries are distributed as separate packages per platform

## go.mod Content:

```go
module github.com/negokaz/excel-mcp-server

go 1.23.0

toolchain go1.24.0

require (
    github.com/0udwins/zog v0.21.4
    github.com/go-ole/go-ole v1.3.0
    github.com/goccy/go-yaml v1.18.0
    github.com/mark3labs/mcp-go v0.34.0
    github.com/skanehira/clipboard-image v1.0.0
    github.com/xuri/excelize/v2 v2.9.2-0.20250717000717-d0d7139785fe
)

require (
    github.com/google/uuid v1.6.0 // indirect
    github.com/richardlehane/mscfb v1.0.4 // indirect
    github.com/richardlehane/msoleps v1.0.4 // indirect
    github.com/spf13/cast v1.9.2 // indirect
    github.com/tiende/go-deepcopy v1.6.1 // indirect
    github.com/xuri/efp v0.0.1 // indirect
    github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a0 // indirect
    github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
    golang.org/x/crypto v0.40.0 // indirect
    golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
    golang.org/x/net v0.42.0 // indirect
    golang.org/x/sys v0.34.0 // indirect
)
```

**Key dependencies:**
- **github.com/mark3labs/mcp-go** - MCP protocol implementation for Go
- **github.com/xuri/excelize/v2** - Excel file manipulation library
- **github.com/go-ole/go-ole** - OLE automation (for Windows Excel interaction)
- **github.com/0udwins/zog** - Schema validation library
- **github.com/goccy/go-yaml** - YAML parsing
- **github.com/skanehira/clipboard-image** - Clipboard image handling

## Visio File Format (.vsdx) Information:

### File Format Overview:
- **VSDX** is based on Open Packaging Conventions (OPC) and XML
- It's essentially a ZIP container with XML files and other resources
- Replaces older binary (.vsd) and XML (.vdx) formats

### File Types:
- .vsdx (Visio drawing)
- .vsdm (Visio macro-enabled drawing)
- .vssx (Visio stencil)
- .vssm (Visio macro-enabled stencil)
- .vstx (Visio template)
- .vstm (Visio macro-enabled template)

### Structure:
- **Package**: ZIP container
- **Package Parts**: XML files, images, VBA solutions
- **Document Parts**: Actual content and metadata
- **Relationship Parts**: Define how parts relate to each other (*.rels files)

### Key Differences from VDX:
- Uses packaging (ZIP) instead of standalone XML
- XML divided into multiple parts
- Uses Cell, Row, Section elements instead of named elements
- ShapeSheet cells represented uniformly as Cell elements with N attribute

### Programmatic Access:
- Can be opened as ZIP file
- XML manipulation for content
- Similar to other Office Open XML formats
- No direct Go library found (unlike Excel with excelize)
