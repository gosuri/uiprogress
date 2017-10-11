package uiprogress

import (
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestBarPrepend(t *testing.T) {
	b := NewBar(100)
	b.PrependCompleted()
	b.Set(50)
	if !strings.Contains(b.String(), "50") {
		t.Fatal("want", "50%", "in", b.String())
	}
}

// BarTest creates a bar with the given total then kicks off 10000 goroutines, each running the given func. At the end, the Current value of the bar is checked against expected.
func BarTest(t *testing.T, total, expected int, f func(*Bar, int)) {
	b := NewBar(total)
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			f(b, i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
		}(i)
	}
	wg.Wait()
	if b.Current() != expected {
		t.Fatal("need", expected, "got", b.Current())
	}
}

func TestBarIncr(t *testing.T) {
	BarTest(t, 10000, 10000, func(b *Bar, i int) {
		b.Incr()
	})
}

func TestBarIncrWithArg(t *testing.T) {
	BarTest(t, 100000000, 50005000, func(b *Bar, i int) {
		b.Incr()
		b.Incr(i)
	})
}
