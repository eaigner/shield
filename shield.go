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

func (sh *shield) Classify(text string) (c string, err error) {
	totalCounts, err := sh.store.TotalClassWordCounts()
	if err != nil {
		return
	}

	// Compute priors
	var sum int64
	for _, v := range totalCounts {
		sum += v
	}
	priors := make(map[string]float64)
	classes := make([]string, 0, len(totalCounts))
	for class, count := range totalCounts {
		classes = append(classes, class)
		priors[class] = float64(count) / float64(sum)
	}

	// Get class word counts in bulk
	tokens := sh.tokenizer.Tokenize(text)
	words := make([]string, 0, len(tokens))
	for word, _ := range tokens {
		words = append(words, word)
	}

	classWordCounts := make(map[string]map[string]int64)
	for _, class := range classes {
		wc, err2 := sh.store.ClassWordCounts(class, words)
		if err2 != nil {
			err = err2
			return
		}
		classWordCounts[class] = wc
	}

	// Compute score
	scores := make(map[string]float64)
	for class, v := range priors {
		score := math.Log(v)
		for _, count := range classWordCounts[class] {
			score += math.Log((float64(count) + defaultProb) / float64(totalCounts[class]))
		}
		scores[class] = score
	}

	// Select class with highes prob
	var k string = ""
	var i float64
	for k2, v2 := range scores {
		if i == 0 || v2 > i {
			k, i = k2, v2
		}
	}
	c = k
	return
}

func (sh *shield) BulkClassify(texts []string) (c []string, err error) {
	panic("TODO: impl!")
}

func (sh *shield) Reset() error {
	return sh.store.Reset()
}
