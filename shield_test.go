package shield

import (
	"github.com/eaigner/goredis"
	"testing"
)

var redis = "127.0.0.1:6379"

func TestLearn(t *testing.T) {
	var client goredis.Client
	client.Addr = "127.0.0.1:6379"

	sh := New(&client, NewEnglishTokenizer())
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
