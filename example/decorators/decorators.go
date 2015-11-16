package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

var steps = []string{
	"downloading source",
	"installing deps",
	"compiling",
	"packaging",
	"seeding database",
	"deploying",
	"staring servers",
}

func main() {
	fmt.Println("apps: deployment started: app1, app2")
	uiprogress.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go app("app1", &wg)
	wg.Add(1)
	go app("app2", &wg)
	wg.Wait()

	fmt.Println("apps: successfully deployed: app1, app2")
}

func app(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	bar := uiprogress.AddBar(len(steps)).
		AppendCompleted().
		PrependElapsed()
	bar.Width = 50

	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return strutil.Resize(name+": "+steps[b.Current()-1], 22)
	})

	rand.Seed(500)
	for i := 0; i < bar.Total; i++ {
		bar.Set(i + 1)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
	}
}
