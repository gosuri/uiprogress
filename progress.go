package uiprogress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gosuri/uilive"
)

var (
	// Out is the default writer to render progress bars to
	Out = os.Stdout
	// DefaultProgress is the default progress
	DefaultProgress = New()
	// RefreshInterval in the default time duration to wait for refreshing the output
	RefreshInterval = time.Millisecond * 10
)

// Progress represents the container that renders progress bars
type Progress struct {
	// Out is the writer to render progress bars to
	Out io.Writer

	// Width is the width of the progress bars
	Width int

	// Bars is the collection of progress bars
	Bars []*Bar

	// RefreshInterval in the time duration to wait for refreshing the output
	RefreshInterval time.Duration

	lw *uilive.Writer

	stopChan chan struct{}
}

// New returns a new progress bar with defaults
func New() *Progress {
	return &Progress{
		Width:           Width,
		Out:             Out,
		Bars:            make([]*Bar, 0),
		RefreshInterval: RefreshInterval,

		lw:       uilive.New(),
		stopChan: make(chan struct{}),
	}
}

// AddBar creates a new progress bar and adds it to the default progress container
func AddBar(total int) *Bar {
	return DefaultProgress.AddBar(total)
}

// Start starts the rendering the progress of progress bars using the DefaultProgress. It listens for updates using `bar.Set(n)` and new bars when added using `AddBar`
func Start() {
	DefaultProgress.Start()
}

// Stop stops listening
func Stop() {
	DefaultProgress.Stop()
}

// Listen listens for updates and renders the progress bars
func Listen() {
	DefaultProgress.Listen()
}

// AddBar creates a new progress bar and adds to the container
func (p *Progress) AddBar(total int) *Bar {
	bar := NewBar(total)
	bar.Width = p.Width
	p.Bars = append(p.Bars, bar)
	return bar
}

// Listen listens for updates and renders the progress bars
func (p *Progress) Listen() {
	p.lw.Out = p.Out
	go func() {
		for {
			time.Sleep(p.RefreshInterval)
			for _, bar := range p.Bars {
				fmt.Fprintln(p.lw, bar.String())
			}
			p.lw.Flush()
		}
	}()
	<-p.stopChan
	time.Sleep(p.RefreshInterval)
	p.lw.Flush()
}

// Start starts the rendering the progress of progress bars. It listens for updates using `bar.Set(n)` and new bars when added using `AddBar`
func (p *Progress) Start() {
	go p.Listen()
}

// Stop stops listening
func (p *Progress) Stop() {
	p.stopChan <- struct{}{}
}
