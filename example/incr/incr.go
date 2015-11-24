package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	count := 1000
	bar := uiprogress.AddBar(count).AppendCompleted().PrependElapsed()
	// prepend the task progress to the bar
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Task (%d/%d)", b.Current(), count)
	})

	uiprogress.Start()
	var wg sync.WaitGroup
	// Fanout into 1k go routines
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bar.Incr()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		}()
	}
	wg.Wait()
	uiprogress.Stop()
}
