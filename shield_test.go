package shield

import (
	"testing"
)

func TestLearn(t *testing.T) {
	store := NewRedisStore("127.0.0.1:6379", "", 0)
	tokenizer := NewEnglishTokenizer()
	sh := New(tokenizer, store)

	sh.Learn("good", "sunshine drugs love sex lobster sloth")
	sh.Learn("bad", "fear death horror government zombie god")

	c, err := sh.Classify("sloths are so cute i love them")
	if err != nil {
		t.Fatal(err)
	}
	if c != "good" {
		t.Fatal(c)
	}

	c, err = sh.Classify("i fear god and love the government")
	if err != nil {
		t.Fatal(err)
	}
	if c != "bad" {
		t.Fatal(c)
	}
}
