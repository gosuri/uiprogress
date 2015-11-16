package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	count1 := 10
	count2 := 20
	p := uiprogress.New()

	bar1 := p.AddBar(count1)
	bar2 := p.AddBar(count2)

	p.Start()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= count1; i++ {
			bar1.Current = i
			time.Sleep(time.Millisecond * 100)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= count2; i++ {
			bar2.Current = i
			time.Sleep(time.Millisecond * 200)
		}
	}()
	wg.Wait()
}
