package main

import (
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {
	uiprogress.Start() // start rendering

	bar := uiprogress.AddBar(100) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	bar.PrependElapsed()

	for i := 1; i <= bar.Total; i++ {
		bar.Set(i)
		time.Sleep(time.Millisecond * 100)
	}
}
