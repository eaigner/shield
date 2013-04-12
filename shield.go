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

func (sh *shield) Learn(class, text string) (err error) {
	if len(class) == 0 {
		panic("invalid class")
	}
	if len(text) == 0 {
		panic("invalid text")
	}
	return sh.BulkLearn([]Set{Set{Class: class, Text: text}})
}

func (sh *shield) BulkLearn(sets []Set) (err error) {
	if len(sets) == 0 {
		panic("invalid data set")
	}
	m := make(map[string]map[string]int64)
	for _, set := range sets {
		if w, ok := m[set.Class]; ok {
			for word, count := range sh.tokenizer.Tokenize(set.Text) {
				w[word] += count
			}
		} else {
			m[set.Class] = sh.tokenizer.Tokenize(set.Text)
		}
	}
	for class, _ := range m {
		if err = sh.store.AddClass(class); err != nil {
			return
		}
	}
	return sh.store.IncrementClassWordCounts(m)
}

func (sh *shield) Forget(class, text string) (err error) {
	return nil // TODO: implement
}

func getKeys(m map[string]int64) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func getWordProb(freqs map[string]int64, word string, totalClassWordCount int64) float64 {
	var p float64
	if v, ok := freqs[word]; ok {
		p = float64(v) / float64(totalClassWordCount)
	}
	// We must not return 0, log(0) is not defined!
	if p == 0 {
		p = defaultProb
	}
	return p
}

func (s *shield) Classify(text string) (c string, err error) {
	// Get total class word counts
	totals, err := s.store.TotalClassWordCounts()
	if err != nil {
		return
	}
	classes := getKeys(totals)

	// Tokenize text
	wordFreqs := s.tokenizer.Tokenize(text)
	words := getKeys(wordFreqs)

	// Get word frequencies for each class
	classFreqs := make(map[string]map[string]int64)
	for _, class := range classes {
		freqs, err2 := s.store.ClassWordCounts(class, words)
		if err2 != nil {
			err = err2
			return
		}
		classFreqs[class] = freqs
	}

	// Calculate log scores for each class
	scores := make(map[string]float64, len(classes))
	for _, class := range classes {
		freqs := classFreqs[class]
		total := totals[class]

		// Because this classifier is not biased, we don't use prior probabilities
		score := float64(0)
		for _, word := range words {
			score += math.Log(getWordProb(freqs, word, total))
		}
		scores[class] = score
	}

	// Select class with highes prob
	var k string
	var i float64
	for k2, v2 := range scores {
		if i == 0 || v2 > i {
			k, i = k2, v2
		}
	}
	c = k
	return
}

func (sh *shield) Reset() error {
	return sh.store.Reset()
}
