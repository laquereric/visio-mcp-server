package tools

import (
	"encoding/json"
	"fmt"

	"github.com/negokaz/visio-mcp-server/internal/visio"
)

// DescribePagesHandler handles the visio_describe_pages tool
func DescribePagesHandler(arguments map[string]interface{}) (*string, error) {
	fileAbsolutePath, ok := arguments["fileAbsolutePath"].(string)
	if !ok {
		return nil, fmt.Errorf("fileAbsolutePath is required")
	}

	// Check if file exists
	if !visio.FileExists(fileAbsolutePath) {
		return nil, fmt.Errorf("file not found: %s", fileAbsolutePath)
	}

	// Read pages
	reader := visio.NewReader(fileAbsolutePath)
	pages, err := reader.ListPages()
	if err != nil {
		return nil, fmt.Errorf("failed to list pages: %w", err)
	}

	// Format response
	response := map[string]interface{}{
		"file":      fileAbsolutePath,
		"pageCount": len(pages),
		"pages":     pages,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	result := string(jsonData)
	return &result, nil
}

// ReadPageHandler handles the visio_read_page tool
func ReadPageHandler(arguments map[string]interface{}) (*string, error) {
	fileAbsolutePath, ok := arguments["fileAbsolutePath"].(string)
	if !ok {
		return nil, fmt.Errorf("fileAbsolutePath is required")
	}

	pageName, ok := arguments["pageName"].(string)
	if !ok {
		return nil, fmt.Errorf("pageName is required")
	}

	// Check if file exists
	if !visio.FileExists(fileAbsolutePath) {
		return nil, fmt.Errorf("file not found: %s", fileAbsolutePath)
	}

	// Read page
	reader := visio.NewReader(fileAbsolutePath)
	page, err := reader.ReadPage(pageName)
	if err != nil {
		return nil, fmt.Errorf("failed to read page: %w", err)
	}

	// Format response
	response := map[string]interface{}{
		"file":       fileAbsolutePath,
		"pageName":   page.Name,
		"width":      page.Width,
		"height":     page.Height,
		"shapeCount": len(page.Shapes),
		"shapes":     page.Shapes,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	result := string(jsonData)
	return &result, nil
}

// ListShapesHandler handles the visio_list_shapes tool
func ListShapesHandler(arguments map[string]interface{}) (*string, error) {
	fileAbsolutePath, ok := arguments["fileAbsolutePath"].(string)
	if !ok {
		return nil, fmt.Errorf("fileAbsolutePath is required")
	}

	pageName, ok := arguments["pageName"].(string)
	if !ok {
		return nil, fmt.Errorf("pageName is required")
	}

	// Check if file exists
	if !visio.FileExists(fileAbsolutePath) {
		return nil, fmt.Errorf("file not found: %s", fileAbsolutePath)
	}

	// Read page
	reader := visio.NewReader(fileAbsolutePath)
	page, err := reader.ReadPage(pageName)
	if err != nil {
		return nil, fmt.Errorf("failed to read page: %w", err)
	}

	// Create simplified shape list
	shapeList := make([]map[string]interface{}, 0, len(page.Shapes))
	for _, shape := range page.Shapes {
		shapeList = append(shapeList, map[string]interface{}{
			"id":   shape.ID,
			"text": shape.Text,
			"type": shape.Type,
			"x":    shape.PinX,
			"y":    shape.PinY,
		})
	}

	// Format response
	response := map[string]interface{}{
		"file":       fileAbsolutePath,
		"pageName":   page.Name,
		"shapeCount": len(shapeList),
		"shapes":     shapeList,
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	result := string(jsonData)
	return &result, nil
}

// WriteShapeHandler handles the visio_write_shape tool
func WriteShapeHandler(arguments map[string]interface{}) (*string, error) {
	fileAbsolutePath, ok := arguments["fileAbsolutePath"].(string)
	if !ok {
		return nil, fmt.Errorf("fileAbsolutePath is required")
	}

	pageName, ok := arguments["pageName"].(string)
	if !ok {
		return nil, fmt.Errorf("pageName is required")
	}

	shapeDataRaw, ok := arguments["shapeData"]
	if !ok {
		return nil, fmt.Errorf("shapeData is required")
	}

	createPage := false
	if cp, ok := arguments["createPage"].(bool); ok {
		createPage = cp
	}

	// Parse shape data
	shapeDataMap, ok := shapeDataRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("shapeData must be an object")
	}

	shapeData := visio.ShapeData{
		Text:   getStringValue(shapeDataMap, "text"),
		PinX:   getFloatValue(shapeDataMap, "pinX"),
		PinY:   getFloatValue(shapeDataMap, "pinY"),
		Width:  getFloatValue(shapeDataMap, "width"),
		Height: getFloatValue(shapeDataMap, "height"),
	}

	// Check if file exists
	if !visio.FileExists(fileAbsolutePath) {
		return nil, fmt.Errorf("file not found: %s", fileAbsolutePath)
	}

	// Write shape
	writer := visio.NewWriter(fileAbsolutePath)
	err := writer.WriteShape(pageName, shapeData, createPage)
	if err != nil {
		return nil, fmt.Errorf("failed to write shape: %w", err)
	}

	// Format response
	response := map[string]interface{}{
		"success": true,
		"file":    fileAbsolutePath,
		"page":    pageName,
		"message": "Shape written successfully",
	}

	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	result := string(jsonData)
	return &result, nil
}

// Helper functions

func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloatValue(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	if val, ok := m[key].(int); ok {
		return float64(val)
	}
	return 0.0
}
