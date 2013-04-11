Shield is a bayesian text classifier with flexible tokenizer and backend store support

Currently implemented:

- Redis backend
- English tokenizer

## Example

```go
package main

import (
  "github.com/eaigner/shield"
)

func main() {
  sh := shield.New(
    shield.NewEnglishTokenizer(),
    shield.NewRedisStore("127.0.0.1:6379", "", 0),
  )

  sh.Learn("good", "sunshine drugs love sex lobster sloth")
  sh.Learn("bad", "fear death horror government zombie god")

  c, _ := sh.Classify("sloths are so cute i love them")
  if c != "good" {
    panic(c)
  }

  c, _ = sh.Classify("i fear god and love the government")
  if c != "bad" {
    panic(c)
  }
}
```