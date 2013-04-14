package shield

import (
	"math"
)

const defaultProb float64 = 1e-11

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
	return sh.increment(class, text, 1)
}

func (sh *shield) BulkLearn(sets []Set) (err error) {
	return sh.bulkIncrement(sets, 1)
}

func (sh *shield) Forget(class, text string) (err error) {
	return sh.increment(class, text, -1)
}

func (sh *shield) increment(class, text string, sign int64) (err error) {
	if len(class) == 0 {
		panic("invalid class")
	}
	if len(text) == 0 {
		panic("invalid text")
	}
	return sh.bulkIncrement([]Set{Set{Class: class, Text: text}}, sign)
}

func (sh *shield) bulkIncrement(sets []Set, sign int64) (err error) {
	if len(sets) == 0 {
		panic("invalid data set")
	}
	m := make(map[string]map[string]int64)
	for _, set := range sets {
		tokens := sh.tokenizer.Tokenize(set.Text)
		for k, _ := range tokens {
			tokens[k] *= sign
		}
		if w, ok := m[set.Class]; ok {
			for word, count := range tokens {
				w[word] += count
			}
		} else {
			m[set.Class] = tokens
		}
	}
	for class, _ := range m {
		if err = sh.store.AddClass(class); err != nil {
			return
		}
	}
	return sh.store.IncrementClassWordCounts(m)
}

func getKeys(m map[string]int64) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

func (s *shield) Score(text string) (scores map[string]float64, err error) {
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
	logScores := make(map[string]float64, len(classes))
	for _, class := range classes {
		freqs := classFreqs[class]
		total := totals[class]

		// Because this classifier is not biased, we don't use prior probabilities
		score := float64(0)
		for _, word := range words {
			// Compute the probability that this word belongs to that class
			wordProb := float64(freqs[word]) / float64(total)
			if wordProb == 0 {
				wordProb = defaultProb
			}
			score += math.Log(wordProb)
		}
		logScores[class] = score
	}

	// Normalize the scores
	var min = math.MaxFloat64
	var max = -math.MaxFloat64
	for _, score := range logScores {
		if score > max {
			max = score
		}
		if score < min {
			min = score
		}
	}
	r := max - min
	scores = make(map[string]float64, len(classes))
	for class, score := range logScores {
		if r == 0 {
			scores[class] = 1
		} else {
			scores[class] = (score - min) / r
		}
	}
	return
}

func (s *shield) Classify(text string) (class string, err error) {
	scores, err := s.Score(text)
	if err != nil {
		return
	}

	// Select class with highes prob
	var score float64
	for k, v := range scores {
		if v > score {
			class, score = k, v
		}
	}
	return
}

func (sh *shield) Reset() error {
	return sh.store.Reset()
}
