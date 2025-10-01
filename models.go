package visio

// Document represents a Visio document structure
type Document struct {
	Pages      []Page
	Properties DocumentProperties
}

// DocumentProperties contains document metadata
type DocumentProperties struct {
	Title       string
	Subject     string
	Creator     string
	Keywords    string
	Description string
	Created     string
	Modified    string
}

// Page represents a single page in a Visio document
type Page struct {
	ID         string
	Name       string
	Width      float64
	Height     float64
	Shapes     []Shape
	Background string
}

// Shape represents a shape on a Visio page
type Shape struct {
	ID         string
	Name       string
	Text       string
	Type       string
	PinX       float64 // X coordinate of rotation pin
	PinY       float64 // Y coordinate of rotation pin
	Width      float64
	Height     float64
	Master     string
	Properties map[string]string
}

// PageInfo contains basic page information
type PageInfo struct {
	ID         string
	Name       string
	Width      float64
	Height     float64
	ShapeCount int
	Background string
}

// ShapeData is used for creating or updating shapes
type ShapeData struct {
	Name       string
	Text       string
	Type       string
	PinX       float64
	PinY       float64
	Width      float64
	Height     float64
	Properties map[string]string
}
