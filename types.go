package shield

type Tokenizer interface {
	Tokenize(text string) (words []string)
}

type Shield interface {
	Learn(class, text string) error
	Forget(class, text string) error
	Classify(text string) (c string, err error)
}
