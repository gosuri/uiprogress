package strutil

import (
	"testing"
	"time"
)

func TestResize(t *testing.T) {
	s := "foo"
	got := Resize(s, 5)
	if len(got) != 5 {
		t.Fatal("want", 5, "got", len(got))
	}
	s = "foobar"
	got = Resize(s, 5)

	if got != "fo..." {
		t.Fatal("want", "fo...", "got", got)
	}
}

func TestPadRight(t *testing.T) {
	got := PadRight("foo", 5, '-')
	if got != "foo--" {
		t.Fatal("want", "foo--", "got", got)
	}
}

func TestPadLeft(t *testing.T) {
	got := PadLeft("foo", 5, '-')
	if got != "--foo" {
		t.Fatal("want", "--foo", "got", got)
	}
}

func TestPrettyTime(t *testing.T) {
	d, _ := time.ParseDuration("")
	got := PrettyTime(d)
	if got != "---" {
		t.Fatal("want", "---", "got", got)
	}
}

func TestPrettyTimeFormat(t *testing.T) {
	d, _ := time.ParseDuration("300ms")
	got := PrettyTimeFormat(d, time.Millisecond)
	if got != "300ms" {
		t.Fatal("want", "300ms", "got", got)
	}
	got = PrettyTimeFormat(d, time.Second)
	if got != "0s" {
		t.Fatal("want", "0s", "got", got)
	}
	got = PrettyTimeFormat(d, time.Hour)
	if got != "0s" {
		t.Fatal("want", "0s", "got", got)
	}
}
