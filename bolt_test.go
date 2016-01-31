package shield

import (
	"io/ioutil"
	"testing"
)

func TempFileName() string {
	f, _ := ioutil.TempFile("", "")
	return f.Name()
}

var (
	boltStore = NewBoltStore(TempFileName(), logger, "")
)

func TestBoltLearn(t *testing.T) {
	sh := newShield(boltStore)
	testLearn(t, sh)
}

func TestBoltDecrement(t *testing.T) {
	sh := newShield(boltStore)
	testDecrement(t, sh)
}
