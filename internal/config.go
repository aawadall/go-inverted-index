package internal

// Index configuration 
type Config struct {
	CaseSensitive bool // Whether the index is case sensitive
	MinWordLength int // Minimum length of words to be indexed
	StopWords     []string // List of stop words to ignore during indexing
}