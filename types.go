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
	IncrementClassWordCount(class, word string, i int64) error
	TotalClassWordCounts() (map[string]int64, error)
}
