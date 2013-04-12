package shield

type Tokenizer interface {
	Tokenize(text string) (words map[string]int64)
}

type Set struct {
	Class string
	Text  string
}

type Shield interface {
	// Learn learns a single document
	Learn(class, text string) (err error)

	// BulkLearn learns many documents at once
	BulkLearn(sets []Set) (err error)

	// Forget forgets the document in the specified class
	Forget(class, text string) (err error)

	// Score returns the scores for each class normalized from 0 to 1
	Score(text string) (scores map[string]float64, err error)

	// Classify returns the class with the highest score
	Classify(text string) (c string, err error)

	// Reset clears the storage
	Reset() error
}

type Store interface {
	Classes() ([]string, error)
	AddClass(class string) error
	ClassWordCounts(class string, words []string) (mc map[string]int64, err error)
	IncrementClassWordCounts(m map[string]map[string]int64) error
	TotalClassWordCounts() (map[string]int64, error)
	Reset() error
}
