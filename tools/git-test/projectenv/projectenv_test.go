package projectenv

import "testing"

func TestReplaceIfEmpty(t *testing.T) {
	s := ""
	replacer := "hola"
	replaceIfEmpty(&s, replacer)
	if s != replacer {
		t.Errorf("Expecting %s but got %s", replacer, s)
	}
}
