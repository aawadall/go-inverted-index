package main

import (
	"fmt"
	"github.com/aawadall/go-inverted-index/internal"
)

// main function for testing purposes
func main() {

// Example usage of the InvertedIndex

// Configuration for the index
	config := internal.Config{
		CaseSensitive: false,
		MinWordLength: 2,
		StopWords:     []string{"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for", "of", "with", "by"},
	}

	// Create a new inverted index
	idx := internal.NewInvertedIndex(config)

	// Sample documents with different content types
	documents := []internal.Document{
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

func printResults(idx *internal.InvertedIndex, results []string) {
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