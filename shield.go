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
	wordMap := map[string]int64{}
	for word, count := range sh.tokenizer.Tokenize(text) {
		wordMap[word] += count
	}
	return sh.store.IncrementClassWordCounts(map[string]map[string]int64{
		class: wordMap,
	})
}

func (sh *shield) Forget(class string, text string) error {
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
	for class, count := range totalCounts {
		priors[class] = float64(count) / float64(sum)
	}

	// Compute score
	tokens := sh.tokenizer.Tokenize(text)
	scores := make(map[string]float64)
	for class, v := range priors {
		score := math.Log(v)
		for word, _ := range tokens {
			cwc, _ := sh.store.ClassWordCount(class, word)
			score += math.Log((float64(cwc) + defaultProb) / float64(totalCounts[class]))
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

func (sh *shield) Reset() error {
	return sh.store.Reset()
}
