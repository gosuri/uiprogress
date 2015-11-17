package uiprogress_test

import (
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func Example() {
	uiprogress.Start() // start rendering

	bar := uiprogress.AddBar(100) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	bar.PrependElapsed()

	for i := 1; i <= bar.Total; i++ {
		bar.Set(i)
		time.Sleep(time.Millisecond * 10)
	}
}

func ExampleProgress_AddBar() {
	waitTime := time.Millisecond * 100
	uiprogress.Start()

	var wg sync.WaitGroup

	bar1 := uiprogress.AddBar(20).AppendCompleted().PrependElapsed()
	wg.Add(1)
	// update the progress bars concurrently using a go routine
	go func() {
		defer wg.Done()
		for i := 1; i <= bar1.Total; i++ {
			bar1.Set(i)
			time.Sleep(waitTime)
		}
	}()

	bar2 := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= bar2.Total; i++ {
			bar2.Set(i)
			time.Sleep(waitTime)
		}
	}()

	time.Sleep(time.Second)
	bar3 := uiprogress.AddBar(20).PrependElapsed().AppendCompleted()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= bar3.Total; i++ {
			bar3.Set(i)
			time.Sleep(waitTime)
		}
	}()

	// wait for a collection of goroutines to finish
	wg.Wait()
}

func ExampleDecoratorFunc() {
	var steps = []string{"downloading source", "installing deps", "compiling", "packaging", "seeding database", "deploying", "staring servers"}
	bar := uiprogress.AddBar(len(steps))

	// prepend the current step to the bar
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "app: " + steps[b.Current()-1]
	})

	for i := 0; i < bar.Total; i++ {
		bar.Set(i + 1)
		time.Sleep(time.Millisecond * 100)
	}
}
