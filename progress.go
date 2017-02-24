package uiprogress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gosuri/uilive"
)

// Out is the default writer to render progress bars to
var Out = os.Stdout

// RefreshInterval in the default time duration to wait for refreshing the output
var RefreshInterval = time.Millisecond * 10

// defaultProgress is the default progress
var defaultProgress = New()

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

	lw     *uilive.Writer
	ticker *time.Ticker
	tdone  chan bool
	mtx    *sync.RWMutex
}

// New returns a new progress bar with defaults
func New() *Progress {
	return &Progress{
		Width:           Width,
		Out:             Out,
		Bars:            make([]*Bar, 0),
		RefreshInterval: RefreshInterval,

		lw:  uilive.New(),
		mtx: &sync.RWMutex{},
	}
}

// AddBar creates a new progress bar and adds it to the default progress container
func AddBar(total int) *Bar {
	return defaultProgress.AddBar(total)
}

// Start starts the rendering the progress of progress bars using the DefaultProgress. It listens for updates using `bar.Set(n)` and new bars when added using `AddBar`
func Start() {
	defaultProgress.Start()
}

// Stop stops listening
func Stop() {
	defaultProgress.Stop()
}

// Listen listens for updates and renders the progress bars
func Listen() {
	defaultProgress.Listen()
}

// AddBar creates a new progress bar and adds to the container
func (p *Progress) AddBar(total int) *Bar {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	bar := NewBar(total)
	bar.Width = p.Width
	p.Bars = append(p.Bars, bar)
	return bar
}

// Listen listens for updates and renders the progress bars
func (p *Progress) Listen() {
	var tickChan = p.ticker.C
	p.lw.Out = p.Out

	for {
		select {
		case <-tickChan:
			p.mtx.RLock()

			if p.ticker != nil {
				p.print()
				p.lw.Flush()
			}

			p.mtx.RUnlock()
		case <-p.tdone:
			if p.ticker != nil {
				p.ticker.Stop()
				p.ticker = nil
				return
			}
		}
	}
}

func (p *Progress) print() {
	for _, bar := range p.Bars {
		fmt.Fprintln(p.lw, bar.String())
	}
}

// Start starts the rendering the progress of progress bars. It listens for updates using `bar.Set(n)` and new bars when added using `AddBar`
func (p *Progress) Start() {
	p.mtx.Lock()
	if p.ticker == nil {
		p.ticker = time.NewTicker(RefreshInterval)
		p.tdone = make(chan bool, 1)
	}
	p.mtx.Unlock()

	go p.Listen()
}

// Stop stops listening
func (p *Progress) Stop() {
	p.mtx.Lock()
	close(p.tdone)
	p.print()
	p.lw.Flush()
	p.mtx.Unlock()
}

// Bypass returns a writer which allows non-buffered data to be written to the underlying output
func (p *Progress) Bypass() io.Writer {
	return p.lw.Bypass()
}
