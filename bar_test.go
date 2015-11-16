package uiprogress

import (
	"strings"
	"testing"
)

func TestBarPrepend(t *testing.T) {
	b := NewBar(100)
	b.PrependCompleted()
	b.Set(50)
	if !strings.Contains(b.String(), "50") {
		t.Fatal("want", "50%", "in", b.String())
	}
}
