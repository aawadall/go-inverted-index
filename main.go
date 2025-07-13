package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/google/uuid"
)

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

// Index configuration 
type Config struct {
	CaseSensitive bool // Whether the index is case sensitive
	MinWordLength int // Minimum length of words to be indexed
	StopWords     []string // List of stop words to ignore during indexing
}

// Inverted Index is a structure that holds the mapping of terms to documents.
type InvertedIndex struct {
	index map[string][]string // Maps terms to document IDs
	documents map[string]Document // Maps document IDs to Document objects
	config Config // Configuration for the index
	stopWords map[string]bool // Set of stop words for quick lookup
}

// NewInvertedIndex creates a new InvertedIndex with the given configuration.
func NewInvertedIndex(config Config) *InvertedIndex {
	idx := &InvertedIndex{
		index:     make(map[string][]string),
		documents: make(map[string]Document),
		config:    config,
		stopWords: make(map[string]bool),
	}

	// Initialize stop words set for quick lookup
	for _, word := range config.StopWords {
		if !idx.config.CaseSensitive {
			word = strings.ToLower(word)
		}
		idx.stopWords[word] = true
	}
	return idx
}

// AddDocument adds a new document to the index.
func (idx *InvertedIndex) AddDocument(doc Document) {
	doc.ID = doc.IDOrGenerate() // Ensure the document has an ID
	idx.documents[doc.ID] = doc

	// Convert content to string for tokenization
	content := idx.contentToString(doc.Content)
	tokens := idx.tokenize(content)

	for _, token := range tokens {
		if !idx.contains(idx.index[token], doc.ID) {
			idx.index[token] = append(idx.index[token], doc.ID)
		}
	}

	for term := range idx.index {
		sort.Strings(idx.index[term]) // Sort document IDs for consistent order
	}
}


// contentToString converts the content of a document to a string.
func (idx *InvertedIndex) contentToString(content interface{}) string {
	// TODO: implement custom content extraction given content type
	switch v := content.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v) // Fallback to fmt.Sprintf for other types
	}
}

// tokenize splits the content into tokens based on whitespace and punctuation.
func (idx *InvertedIndex) tokenize(content string) []string {
	re := regexp.MustCompile(`[^\w\s]+`) // Regex to match non-word characters
	text := re.ReplaceAllString(content, " ") // Replace non-word characters with space
	words := strings.Fields(text) // Split by whitespace
	var tokens []string

	for _, word := range words {
		if !idx.config.CaseSensitive {
			word = strings.ToLower(word) // Normalize case if not case sensitive
		}
		if len(word) >= idx.config.MinWordLength && !idx.stopWords[word] {
			tokens = append(tokens, word) // Only add valid tokens
		}
	}

	return tokens
}

// contains checks if a slice contains a specific element.
func (idx *InvertedIndex) contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Search one term 
func (idx *InvertedIndex) SearchTerm(term string) []string { 
	if !idx.config.CaseSensitive {
		term = strings.ToLower(term) // Normalize case if not case sensitive
	}
	return idx.index[term]
}

// SearchAND searches for documents that contain all of the specified terms.
func (idx *InvertedIndex) SearchAND(terms []string) []string {
	if len(terms) == 0 {
		return nil
	}

	result := idx.SearchTerm(terms[0])

	for i := 1; i < len(terms); i++ {
		termDocs := idx.SearchTerm(terms[i])
		result = idx.intersect(result, termDocs)
	}
	return result
}

// SearchOR searches for documents that contain any of the specified terms.
func (idx *InvertedIndex) SearchOR(terms []string) []string {
	if len(terms) == 0 {
		return nil
	}

	result := []string{}

	for _, term := range terms {
		termDocs := idx.SearchTerm(term)
		result = idx.union(result, termDocs)
	}

	return result
}

// SearchNOT searches for documents that do not contain the specified term.
func (idx *InvertedIndex) SearchNOT(includeTerm, excludeTerm string) []string {
	includeList := idx.SearchTerm(includeTerm)
	excludeList := idx.SearchTerm(excludeTerm)
	return idx.subtract(includeList, excludeList)
}

// intersect returns the intersection of two slices.
func (idx *InvertedIndex) intersect(a, b []string) []string {
	var result []string
	i, j := 0, 0
	
	for i < len(a) && j < len(b) {
		if a[i] == b[j]	{
			result = append(result, a[i])
			i++
			j++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}

	return result
}

// union returns the union of two slices.
func (idx *InvertedIndex) union(a, b []string) []string {
	set := make(map[string]bool)
	var result []string

	// TODO: refactor, reuse following
	for _, id := range a {
		if !set[id] {
			set[id] = true
			result = append(result, id)
		}
	}

	for _, id := range b {
		if !set[id] {
			set[id] = true
			result = append(result, id)
		}
	}

	sort.Strings(result) // Sort the result for consistent order
	return result
}

// subtract returns the difference between two slices.
func (idx *InvertedIndex) subtract(a, b []string) []string {
	bSet := make(map[string]bool)
	for _, id := range b {
		bSet[id] = true
	}

	var result []string
	for _, id := range a {
		if !bSet[id] {
			result = append(result, id)
		}
	}

	return result
}

// PrintIndex prints the inverted index in a readable format.
func (idx *InvertedIndex) PrintIndex() {
	fmt.Println("Inverted Index:")
	for term, docIDs := range idx.index {
		fmt.Printf("%s: %v\n", term, docIDs)
	}
}

// GetDocument retrieves a document by its ID.
func (idx *InvertedIndex) GetDocument(id string) (Document, bool) {
	doc, exists := idx.documents[id]
	return doc, exists
}

// main function for testing purposes
func main() {

// Example usage of the InvertedIndex

// Configuration for the index
	config := Config{
		CaseSensitive: false,
		MinWordLength: 2,
		StopWords:     []string{"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"},
	}

	// Create a new inverted index
	idx := NewInvertedIndex(config)

	// Sample documents with different content types
	documents := []Document{
		{ID: "doc1", Content: "The quick brown fox jumps over the lazy dog"},
		{ID: "doc2", Content: "A quick brown fox is running fast"},
		{ID: "doc3", Content: "The lazy dog is sleeping under the tree"},
		{ID: "doc4", Content: []byte("Fast running foxes are jumping high")},
		{ID: "doc5", Content: "Dogs and foxes are different animals"},
	}

	// Add documents to the index
	for _, doc := range documents {
		idx.AddDocument(doc)
	}

	// Print the index 
	idx.PrintIndex()
	fmt.Println("Index created successfully.")

	// Example searches
	fmt.Println("=== SEARCH EXAMPLES ===")
	
	// Single term search
	fmt.Println("Search 'fox':")
	results := idx.SearchTerm("fox")
	printResults(idx, results)
	
	// AND search
	fmt.Println("Search 'fox AND quick':")
	results = idx.SearchAND([]string{"fox", "quick"})
	printResults(idx, results)
	
	// OR search
	fmt.Println("Search 'quick OR lazy':")
	results = idx.SearchOR([]string{"quick", "lazy"})
	printResults(idx, results)
	
	// NOT search
	fmt.Println("Search 'fox NOT lazy':")
	results = idx.SearchNOT("fox", "lazy")
	printResults(idx, results)
}

func printResults(idx *InvertedIndex, results []string) {
	if len(results) == 0 {
		fmt.Println("No results found")
		fmt.Println()
		return
	}
	
	fmt.Printf("Found %d document(s): %v\n", len(results), results)
	for _, docID := range results {
		if doc, exists := idx.GetDocument(docID); exists {
			content := fmt.Sprintf("%v", doc.Content)
			fmt.Printf("  Doc %s: %s\n", docID, content)
		}
	}
	fmt.Println()
}