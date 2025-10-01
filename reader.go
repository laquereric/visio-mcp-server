package visio

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Reader handles reading Visio files
type Reader struct {
	filePath string
}

// NewReader creates a new Visio file reader
func NewReader(filePath string) *Reader {
	return &Reader{
		filePath: filePath,
	}
}

// ReadDocument reads the entire Visio document
func (r *Reader) ReadDocument() (*Document, error) {
	zipReader, err := zip.OpenReader(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open VSDX file: %w", err)
	}
	defer zipReader.Close()

	doc := &Document{
		Pages: make([]Page, 0),
	}

	// Read document properties
	props, err := r.readDocumentProperties(&zipReader.Reader)
	if err == nil {
		doc.Properties = props
	}

	// Read pages
	pages, err := r.readPages(&zipReader.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read pages: %w", err)
	}
	doc.Pages = pages

	return doc, nil
}

// ListPages returns basic information about all pages
func (r *Reader) ListPages() ([]PageInfo, error) {
	zipReader, err := zip.OpenReader(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open VSDX file: %w", err)
	}
	defer zipReader.Close()

	return r.listPagesFromZip(&zipReader.Reader)
}

// ReadPage reads a specific page by name
func (r *Reader) ReadPage(pageName string) (*Page, error) {
	zipReader, err := zip.OpenReader(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open VSDX file: %w", err)
	}
	defer zipReader.Close()

	pages, err := r.readPages(&zipReader.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read pages: %w", err)
	}

	for _, page := range pages {
		if page.Name == pageName {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("page not found: %s", pageName)
}

// readDocumentProperties reads document metadata
func (r *Reader) readDocumentProperties(zipReader *zip.Reader) (DocumentProperties, error) {
	props := DocumentProperties{}

	// Read core properties
	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, "docProps/core.xml") {
			rc, err := file.Open()
			if err != nil {
				return props, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return props, err
			}

			// Simple XML parsing for core properties
			// In production, use proper XML unmarshaling
			content := string(data)
			props.Title = extractXMLValue(content, "dc:title")
			props.Subject = extractXMLValue(content, "dc:subject")
			props.Creator = extractXMLValue(content, "dc:creator")
			props.Keywords = extractXMLValue(content, "cp:keywords")
			props.Description = extractXMLValue(content, "dc:description")

			break
		}
	}

	return props, nil
}

// listPagesFromZip lists pages from zip reader
func (r *Reader) listPagesFromZip(zipReader *zip.Reader) ([]PageInfo, error) {
	pageInfos := make([]PageInfo, 0)

	// Find and read document.xml to get page list
	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, "visio/document.xml") {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			// Parse document.xml to extract page information
			// This is a simplified version - production code should use proper XML parsing
			content := string(data)
			_ = content // TODO: Parse XML to extract page information

			break
		}
	}

	// For now, scan for page files
	pageFiles := make(map[string]bool)
	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "visio/pages/page") && strings.HasSuffix(file.Name, ".xml") {
			if !strings.Contains(file.Name, "_rels") {
				pageFiles[file.Name] = true
			}
		}
	}

	// Read each page file
	for pageFile := range pageFiles {
		pageInfo, err := r.readPageInfo(zipReader, pageFile)
		if err == nil {
			pageInfos = append(pageInfos, pageInfo)
		}
	}

	return pageInfos, nil
}

// readPages reads all pages from the document
func (r *Reader) readPages(zipReader *zip.Reader) ([]Page, error) {
	pages := make([]Page, 0)

	// Find page files
	pageFiles := make(map[string]bool)
	for _, file := range zipReader.File {
		if strings.Contains(file.Name, "visio/pages/page") && strings.HasSuffix(file.Name, ".xml") {
			if !strings.Contains(file.Name, "_rels") {
				pageFiles[file.Name] = true
			}
		}
	}

	// Read each page
	for pageFile := range pageFiles {
		page, err := r.readPageFromFile(zipReader, pageFile)
		if err == nil {
			pages = append(pages, page)
		}
	}

	return pages, nil
}

// readPageInfo reads basic page information
func (r *Reader) readPageInfo(zipReader *zip.Reader, pageFile string) (PageInfo, error) {
	info := PageInfo{
		ID:   filepath.Base(pageFile),
		Name: "Page",
	}

	for _, file := range zipReader.File {
		if file.Name == pageFile {
			rc, err := file.Open()
			if err != nil {
				return info, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return info, err
			}

			content := string(data)
			info.Name = extractXMLValue(content, "Name")
			if info.Name == "" {
				info.Name = strings.TrimSuffix(filepath.Base(pageFile), ".xml")
			}

			// Count shapes (simplified)
			info.ShapeCount = strings.Count(content, "<Shape ")

			break
		}
	}

	return info, nil
}

// readPageFromFile reads a complete page from a file
func (r *Reader) readPageFromFile(zipReader *zip.Reader, pageFile string) (Page, error) {
	page := Page{
		ID:     filepath.Base(pageFile),
		Name:   "Page",
		Shapes: make([]Shape, 0),
	}

	for _, file := range zipReader.File {
		if file.Name == pageFile {
			rc, err := file.Open()
			if err != nil {
				return page, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return page, err
			}

			content := string(data)
			page.Name = extractXMLValue(content, "Name")
			if page.Name == "" {
				page.Name = strings.TrimSuffix(filepath.Base(pageFile), ".xml")
			}

			// Parse shapes (simplified - production should use proper XML parsing)
			page.Shapes = r.parseShapes(content)

			break
		}
	}

	return page, nil
}

// parseShapes extracts shapes from page XML content
func (r *Reader) parseShapes(content string) []Shape {
	shapes := make([]Shape, 0)

	// This is a simplified parser
	// Production code should use encoding/xml with proper struct definitions
	lines := strings.Split(content, "\n")
	var currentShape *Shape

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "<Shape ") {
			currentShape = &Shape{
				Properties: make(map[string]string),
			}
			// Extract ID from Shape tag
			if id := extractAttribute(line, "ID"); id != "" {
				currentShape.ID = id
			}
		}

		if currentShape != nil {
			if strings.Contains(line, "<Cell N=\"PinX\"") {
				if val := extractAttribute(line, "V"); val != "" {
					currentShape.PinX, _ = strconv.ParseFloat(val, 64)
				}
			}
			if strings.Contains(line, "<Cell N=\"PinY\"") {
				if val := extractAttribute(line, "V"); val != "" {
					currentShape.PinY, _ = strconv.ParseFloat(val, 64)
				}
			}
			if strings.Contains(line, "<Cell N=\"Width\"") {
				if val := extractAttribute(line, "V"); val != "" {
					currentShape.Width, _ = strconv.ParseFloat(val, 64)
				}
			}
			if strings.Contains(line, "<Cell N=\"Height\"") {
				if val := extractAttribute(line, "V"); val != "" {
					currentShape.Height, _ = strconv.ParseFloat(val, 64)
				}
			}
			if strings.Contains(line, "<Text>") {
				text := extractXMLValue(line, "Text")
				currentShape.Text = text
			}
		}

		if strings.Contains(line, "</Shape>") && currentShape != nil {
			shapes = append(shapes, *currentShape)
			currentShape = nil
		}
	}

	return shapes
}

// Helper functions for simple XML parsing

func extractXMLValue(content, tag string) string {
	startTag := "<" + tag + ">"
	endTag := "</" + tag + ">"

	startIdx := strings.Index(content, startTag)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(startTag)

	endIdx := strings.Index(content[startIdx:], endTag)
	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(content[startIdx : startIdx+endIdx])
}

func extractAttribute(line, attr string) string {
	attrStr := attr + "=\""
	startIdx := strings.Index(line, attrStr)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(attrStr)

	endIdx := strings.Index(line[startIdx:], "\"")
	if endIdx == -1 {
		return ""
	}

	return line[startIdx : startIdx+endIdx]
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
