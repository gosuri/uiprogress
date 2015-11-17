# uiprogress [![GoDoc](https://godoc.org/github.com/gosuri/uiprogress?status.svg)](https://godoc.org/github.com/gosuri/uiprogress) [![Build Status](https://travis-ci.org/gosuri/uiprogress.svg?branch=master)](https://travis-ci.org/gosuri/uiprogress)

A Go library to render progress bars in terminal applications. It provides a set of powerful features that are customizable using a simple interface.

![example](doc/example_full.gif)

Progress bars improve readability for terminal applications with long outputs by providing a concise feedback loop.

## Features

* __Multiple Bars__: uiprogress can render multiple progress bars that can be tracked concurrently
* __Dynamic Addition__:  Add additional progress bars any time, even after the progress tracking has started
* __Prepend and Append Functions__: Append or prepend completion percent and time elapsed to the progress bars
* __Custom Decorator Functions__: Add custom functions around the bar along with helper functions

## Usage

To start listening for progress bars, call `uiprogress.Start()` and add a progress bar using `uiprogress.AddBar(total)`. Update the progress using `bar.Set(n)`. Full source code for the below example is available at [example/simple/simple.go](example/simple/simple.go) 

```go
uiprogress.Start() // start rendering
bar := uiprogress.AddBar(100) // Add a new bar
// optionally, append and prepend completion and elapsed time
bar.AppendCompleted()
bar.PrependElapsed()

for i := 1; i <= bar.Total; i++ {
  bar.Set(i)
  time.Sleep(time.Millisecond * 10)
}
```

This will render the below in the terminal

![example](doc/example_simple.gif)

### Using Custom Decorators

You can also add a custom decorator function in addition to default `bar.AppendCompleted()` and `bar.PrependElapsed()` decorators. The below example tracks the current step for an application deploy progress. Source code for the below example is available at [example/full/full.go](example/full/full.go) 

```go
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
```

### Rendering Multiple bars

You can add multiple bars using `uiprogress.AddBar(n)`. The below example demonstrates updating multiple bars concurrently and adding a new bar later in the pipeline. Source for this example is available at [example/multi/multi.go](example/multi/multi.go) 

```go
waitTime := time.Millisecond * 100
uiprogress.Start()
var wg sync.WaitGroup

// update the progress bars concurrently using a go routine
bar1 := uiprogress.AddBar(20).AppendCompleted().PrependElapsed()
wg.Add(1)
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

time.Sleep(time.Second) // some logic
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
```

This will produce

![example](doc/example_multi.gif)

## Installation

```sh
$ go get -v github.com/gosuri/uiprogress
```
