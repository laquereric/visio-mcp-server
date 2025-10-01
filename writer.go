package visio

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Writer handles writing to Visio files
type Writer struct {
	filePath string
}

// NewWriter creates a new Visio file writer
func NewWriter(filePath string) *Writer {
	return &Writer{
		filePath: filePath,
	}
}

// WriteShape writes or updates a shape on a page
func (w *Writer) WriteShape(pageName string, shapeData ShapeData, createPage bool) error {
	// Check if file exists
	if !FileExists(w.filePath) {
		return fmt.Errorf("file does not exist: %s", w.filePath)
	}

	// Read existing file
	zipReader, err := zip.OpenReader(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to open VSDX file: %w", err)
	}
	defer zipReader.Close()

	// Create temporary file for writing
	tempFile := w.filePath + ".tmp"
	outFile, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	pageFound := false
	pageModified := false

	// Copy existing files and modify target page
	for _, file := range zipReader.File {
		if w.isTargetPageFile(file.Name, pageName) {
			pageFound = true
			// Modify this page
			err := w.modifyPageWithShape(file, zipWriter, shapeData)
			if err != nil {
				return fmt.Errorf("failed to modify page: %w", err)
			}
			pageModified = true
		} else {
			// Copy file as-is
			err := w.copyZipFile(file, zipWriter)
			if err != nil {
				return fmt.Errorf("failed to copy file: %w", err)
			}
		}
	}

	if !pageFound && !createPage {
		return fmt.Errorf("page not found: %s", pageName)
	}

	if !pageModified && createPage {
		// TODO: Create new page
		return fmt.Errorf("page creation not yet implemented")
	}

	// Close writers
	zipWriter.Close()
	outFile.Close()
	zipReader.Close()

	// Replace original file with modified file
	err = os.Rename(tempFile, w.filePath)
	if err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to replace original file: %w", err)
	}

	return nil
}

// isTargetPageFile checks if a file is the target page
func (w *Writer) isTargetPageFile(fileName, pageName string) bool {
	// Simplified check - should match page by name from document.xml
	return strings.Contains(fileName, "visio/pages/page") &&
		strings.HasSuffix(fileName, ".xml") &&
		!strings.Contains(fileName, "_rels")
}

// modifyPageWithShape modifies a page file to add/update a shape
func (w *Writer) modifyPageWithShape(file *zip.File, zipWriter *zip.Writer, shapeData ShapeData) error {
	// Read original page content
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	content := string(data)

	// Add shape to page XML
	// This is a simplified implementation
	// Production code should use proper XML manipulation
	modifiedContent := w.addShapeToPageXML(content, shapeData)

	// Write modified content
	writer, err := zipWriter.Create(file.Name)
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(modifiedContent))
	return err
}

// addShapeToPageXML adds a shape to page XML content
func (w *Writer) addShapeToPageXML(content string, shapeData ShapeData) string {
	// Find the Shapes section
	shapesEndTag := "</Shapes>"
	idx := strings.Index(content, shapesEndTag)

	if idx == -1 {
		// No Shapes section, return as-is
		return content
	}

	// Create shape XML
	shapeXML := w.generateShapeXML(shapeData)

	// Insert before </Shapes>
	return content[:idx] + shapeXML + content[idx:]
}

// generateShapeXML generates XML for a shape
func (w *Writer) generateShapeXML(shapeData ShapeData) string {
	// Simplified shape XML generation
	// Production code should use proper XML encoding
	xml := fmt.Sprintf(`
    <Shape ID="1" Type="Shape">
        <Cell N="PinX" V="%.2f"/>
        <Cell N="PinY" V="%.2f"/>
        <Cell N="Width" V="%.2f"/>
        <Cell N="Height" V="%.2f"/>
        <Text>%s</Text>
    </Shape>
`, shapeData.PinX, shapeData.PinY, shapeData.Width, shapeData.Height, shapeData.Text)

	return xml
}

// copyZipFile copies a file from source zip to destination zip
func (w *Writer) copyZipFile(file *zip.File, zipWriter *zip.Writer) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	writer, err := zipWriter.Create(file.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, rc)
	return err
}

// CreateNewDocument creates a new Visio document
func (w *Writer) CreateNewDocument() error {
	// Create a minimal VSDX file structure
	outFile, err := os.Create(w.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Create [Content_Types].xml
	err = w.writeContentTypes(zipWriter)
	if err != nil {
		return err
	}

	// Create _rels/.rels
	err = w.writeRootRels(zipWriter)
	if err != nil {
		return err
	}

	// Create docProps/core.xml
	err = w.writeCoreProperties(zipWriter)
	if err != nil {
		return err
	}

	// Create visio/document.xml
	err = w.writeDocument(zipWriter)
	if err != nil {
		return err
	}

	// Create visio/pages/page1.xml
	err = w.writeDefaultPage(zipWriter)
	if err != nil {
		return err
	}

	return nil
}

// Helper methods to write minimal VSDX structure

func (w *Writer) writeContentTypes(zipWriter *zip.Writer) error {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Default Extension="xml" ContentType="application/xml"/>
    <Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
    <Override PartName="/visio/document.xml" ContentType="application/vnd.ms-visio.drawing.main+xml"/>
    <Override PartName="/visio/pages/page1.xml" ContentType="application/vnd.ms-visio.page+xml"/>
</Types>`

	writer, err := zipWriter.Create("[Content_Types].xml")
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(content))
	return err
}

func (w *Writer) writeRootRels(zipWriter *zip.Writer) error {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="visio/document.xml"/>
    <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
</Relationships>`

	writer, err := zipWriter.Create("_rels/.rels")
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(content))
	return err
}

func (w *Writer) writeCoreProperties(zipWriter *zip.Writer) error {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
                   xmlns:dc="http://purl.org/dc/elements/1.1/"
                   xmlns:dcterms="http://purl.org/dc/terms/"
                   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <dc:creator>Visio MCP Server</dc:creator>
    <dc:title>New Diagram</dc:title>
</cp:coreProperties>`

	writer, err := zipWriter.Create("docProps/core.xml")
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(content))
	return err
}

func (w *Writer) writeDocument(zipWriter *zip.Writer) error {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<VisioDocument xmlns="http://schemas.microsoft.com/office/visio/2012/main">
    <Pages>
        <Page ID="0" Name="Page-1"/>
    </Pages>
</VisioDocument>`

	writer, err := zipWriter.Create("visio/document.xml")
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(content))
	return err
}

func (w *Writer) writeDefaultPage(zipWriter *zip.Writer) error {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<PageContents xmlns="http://schemas.microsoft.com/office/visio/2012/main">
    <Shapes>
    </Shapes>
</PageContents>`

	writer, err := zipWriter.Create("visio/pages/page1.xml")
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(content))
	return err
}
