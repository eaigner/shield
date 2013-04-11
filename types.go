package shield

type Tokenizer interface {
	Tokenize(text string) (words map[string]int64)
}

type Shield interface {
	Learn(class, text string) error
	Forget(class, text string) error
	Classify(text string) (c string, err error)
	Reset() error
}

type Store interface {
	Classes() ([]string, error)
	AddClass(class string) error
	ClassWordCount(class, word string) (int64, error)
	IncrementClassWordCounts(m map[string]map[string]int64) error
	TotalClassWordCounts() (map[string]int64, error)
	Reset() error
}
