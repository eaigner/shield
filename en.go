package shield

import (
	"regexp"
)

type enTokenizer struct {
}

func NewEnglishTokenizer() Tokenizer {
	return &enTokenizer{}
}

func (t *enTokenizer) Tokenize(text string) (words map[string]int64) {
	words = make(map[string]int64)
	for _, w := range splitTokenRx.Split(text, -1) {
		words[w]++
	}
	return
}

var splitTokenRx = regexp.MustCompile(`\s+|able|all|already|and|any|are|because|both|can|come|each|email|even|few|first|for|from|give|has|have|http|information|into|it's|just|know|like|long|look|made|mail|mailing|mailto|make|many|more|most|much|need|not|now|number|off|one|only|out|own|people|place|right|same|see|such|that|the|this|through|time|using|web|where|why|with|without|work|world|year|years|you|you're|your`)
