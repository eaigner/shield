package shield

import (
	"github.com/eaigner/goredis"
	"math"
	"strconv"
	"strings"
)

const defaultProb float64 = 0.00000000001

type shield struct {
	t          Tokenizer
	redis      *goredis.Client
	rootKey    string
	sumKey     string
	classKey   string
	classesKey string
}

func New(client *goredis.Client, t Tokenizer) Shield {
	return NewPrefix(client, t, "shield")
}

func NewPrefix(client *goredis.Client, t Tokenizer, rootKey string) Shield {
	return &shield{
		t:          t,
		redis:      client,
		rootKey:    rootKey,
		sumKey:     rkey(rootKey, "sum"),
		classKey:   rkey(rootKey, "class"),
		classesKey: rkey(rootKey, "classes"),
	}
}

func (sh *shield) Learn(class string, text string) (err error) {
	if class == "" {
		panic("no class specified")
	}
	_, err = sh.redis.Sadd(sh.classesKey, []byte(class))
	if err != nil {
		return
	}

	// Update class word count
	var numIncrs int64
	defer func() {
		if numIncrs > 0 {
			err2 := sh.incTotalClassWordCount(class, numIncrs)
			if err == nil {
				err = err2
			}
		}
	}()

	// Update specific word count
	for word, count := range group(sh.t.Tokenize(text)) {
		err = sh.incClassWordCount(class, word, count)
		if err != nil {
			return
		}
		numIncrs++
	}

	return
}

func (sh *shield) Forget(class string, text string) error {
	// TODO: implement
	return nil
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
			k = k2
			i = v2
		}
	}
	c = k
	return
}

func (sh *shield) score(text string) (m map[string]float64, err error) {
	classes, err := sh.redis.Smembers(sh.classesKey)
	if err != nil {
		return
	}

	wcs, err := sh.classWordCounts()
	if err != nil {
		return
	}

	m = make(map[string]float64)
	grouped := group(sh.t.Tokenize(text))
	for _, v := range classes {
		class := string(v)
		wc := wcs[class]
		if wc == 0 {
			continue
		}
		m[class] = 0.0
		for word, _ := range grouped {
			score, err := sh.classWordCount(class, word)
			fscore := float64(score)
			if err != nil || score == 0 {
				fscore = defaultProb
			}
			m[class] += math.Log(fscore / float64(wc))
		}
	}
	return
}

func (sh *shield) classWordCount(class string, word string) (i int64, err error) {
	key := rkey(sh.classKey, class)
	b, err := sh.redis.Hget(key, word)
	if err != nil {
		return
	}
	i = btoi(b)
	return
}

func (sh *shield) incClassWordCount(class string, word string, inc int64) error {
	_, err := sh.redis.Hincrby(rkey(sh.classKey, class), word, inc)
	return err
}

func (sh *shield) classWordCounts() (m map[string]int64, err error) {
	v := make(map[string]int64)
	err = sh.redis.Hgetall(sh.sumKey, &v)
	if err == nil {
		m = v
	}
	return
}

func (sh *shield) incTotalClassWordCount(class string, inc int64) error {
	_, err := sh.redis.Hincrby(sh.sumKey, class, inc)
	return err
}

func rkey(s ...string) string {
	var a []string
	var sep = ":"
	for _, v := range s {
		a2 := strings.Split(v, sep)
		a = append(a, a2...)
	}
	return strings.Join(a, sep)
}

func btoi(b []byte) int64 {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func group(words []string) map[string]int64 {
	m := make(map[string]int64)
	for _, w := range words {
		m[w]++
	}
	return m
}
