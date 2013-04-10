package shield

import (
	"math"
)

const defaultProb float64 = 0.00000000001

type shield struct {
	tokenizer Tokenizer
	store     Store
}

func New(t Tokenizer, s Store) Shield {
	return &shield{
		tokenizer: t,
		store:     s,
	}
}

func (sh *shield) Learn(class string, text string) (err error) {
	if class == "" {
		panic("no class specified")
	}
	if err = sh.store.AddClass(class); err != nil {
		return
	}
	for word, count := range sh.tokenizer.Tokenize(text) {
		if err = sh.store.IncrementClassWordCount(class, word, count); err != nil {
			return
		}
	}
	return
}

func (sh *shield) Forget(class string, text string) error {
	return nil // TODO: implement
}

func (sh *shield) Classify(text string) (c string, err error) {
	m, err := sh.score(text)
	if err != nil {
		return
	}
	var k string = ""
	var i float64
	for k2, v2 := range m {
		if i == 0 || v2 > i {
			k, i = k2, v2
		}
	}
	c = k
	return
}

func (sh *shield) Reset() error {
	return nil // TODO: implement
}

func (sh *shield) score(text string) (m map[string]float64, err error) {
	classes, err := sh.store.Classes()
	if err != nil {
		return
	}
	totalCounts, err := sh.store.TotalClassWordCounts()
	if err != nil {
		return
	}

	m = make(map[string]float64)
	grouped := sh.tokenizer.Tokenize(text)
	for _, class := range classes {
		wordCount := totalCounts[class]
		if wordCount == 0 {
			continue
		}
		m[class] = 0.0
		for word, _ := range grouped {
			c, err := sh.store.ClassWordCount(class, word)
			score := float64(c)
			if c == 0 || err != nil {
				score = defaultProb
			}
			m[class] += math.Log(score / float64(wordCount))
		}
	}
	return
}
