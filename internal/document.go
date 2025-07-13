package internal

import "github.com/google/uuid"

// Document is the building block for indexing and searching.
type Document struct {
	ID   string // Unique identifier for the document
	Content interface{} // Content of the document, can be any type
}

// document ID helper, if no ID is provided, generate a new one
func (d *Document) IDOrGenerate() string {
	if d.ID == "" {
		d.ID = generateNewID() // Assume generateNewID is a function that generates a unique ID
	}
	return d.ID
}

// generateNewID returns uuid
func generateNewID() string {
	return uuid.NewString()
}
