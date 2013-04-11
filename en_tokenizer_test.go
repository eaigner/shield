package shield

import (
	"fmt"
	"testing"
)

func TestTokenize(t *testing.T) {
	tokenizer := NewEnglishTokenizer()
	text := "lorem    ipsum able hello erik    can do hi there  \t  spaaace! lorem"
	m := tokenizer.Tokenize(text)
	x := fmt.Sprintf("%v", m)
	if x != `map[lorem:2 ipsum:1 hello:1 erik:1 spaaace:1]` {
		t.Fatal(x)
	}
}
