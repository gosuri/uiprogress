package uiprogress

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStoppingPrintout(t *testing.T) {
	progress := New()
	progress.SetRefreshInterval(time.Millisecond * 10)

	var buffer = &bytes.Buffer{}
	progress.SetOut(buffer)

	bar := progress.AddBar(100)
	progress.Start()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for i := 0; i <= 80; i = i + 10 {
			bar.Set(i)
			time.Sleep(time.Millisecond * 5)
		}

		wg.Done()
	}()

	wg.Wait()

	progress.Stop()
	fmt.Fprintf(buffer, "foo")

	var wantSuffix = "[======================================================>-------------]\nfoo"

	if !strings.HasSuffix(buffer.String(), wantSuffix) {
		t.Errorf("Content that should be printed after stop not appearing on buffer.")
	}
}

func TestIndexOf(t *testing.T) {
	progress := New()
	var buffer = &bytes.Buffer{}
	progress.SetOut(buffer)
	bar0 := progress.AddBar(100)
	bar1 := progress.AddBar(100)
	bar2 := progress.AddBar(100)
	if len(progress.Bars) != 3 {
		t.Errorf("expected Bars to have length %d, was %d", 3, len(progress.Bars))
	}
	if i := progress.indexOf(bar0); i != 0 {
		t.Errorf("expected bar0 to have index 0, was %d", i)
	}
	if i := progress.indexOf(bar1); i != 1 {
		t.Errorf("expected bar0 to have index 1, was %d", i)
	}
	if i := progress.indexOf(bar2); i != 2 {
		t.Errorf("expected bar0 to have index 2, was %d", i)
	}
}

func TestRemoveBar(t *testing.T) {
	progress := New()
	var buffer = &bytes.Buffer{}
	progress.SetOut(buffer)
	bar0 := progress.AddBar(100)
	bar1 := progress.AddBar(100)
	bar2 := progress.AddBar(100)

	if !progress.RemoveBar(bar1) {
		t.Error("expected RemoveBar to return the Bar")
	}
	if len(progress.Bars) != 2 {
		t.Error("expected len(Bars) to be 2")
	}
	if progress.Bars[0] != bar0 {
		t.Error("expected bar0 to be on the first position")
	}
	if progress.Bars[1] != bar2 {
		t.Error("expected bar2 to be on the second position")
	}
	if progress.RemoveBar(bar1) {
		t.Error("expected removal bar1 again to fail")
	}
}
