package shield

type Tokenizer interface {
	Tokenize(text string) (words map[string]int64)
}

type Set struct {
	Class string
	Text  string
}

type Shield interface {
	Learn(class, text string) (err error)
	BulkLearn(sets []Set) (err error)
	Forget(class, text string) (err error)
	Classify(text string) (c string, err error)
	BulkClassify(texts []string) (c []string, err error)
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
