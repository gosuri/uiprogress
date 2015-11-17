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
	RefreshInterval = time.Millisecond
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
}

// New returns a new progress bar with defaults
func New() *Progress {
	return &Progress{
		Width:           Width,
		Out:             Out,
		Bars:            make([]*Bar, 0),
		RefreshInterval: RefreshInterval,
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

// AddBar creates a new progress bar and adds to the container
func (p *Progress) AddBar(total int) *Bar {
	bar := NewBar(total)
	bar.Width = p.Width
	p.Bars = append(p.Bars, bar)
	return bar
}

// Start starts the rendering the progress of progress bars. It listens for updates using `bar.Set(n)` and new bars when added using `AddBar`
func (p *Progress) Start() {
	lw := uilive.New()
	lw.Out = p.Out
	lw.RefreshInterval = p.RefreshInterval
	go func() {
		for {
			for _, bar := range p.Bars {
				fmt.Fprintln(lw, bar.String())
			}
			lw.Flush()
			lw.Wait()
		}
	}()
}
