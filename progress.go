package uiprogress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gosuri/uilive"
	"golang.org/x/term"
)

// Out is the default writer to render progress bars to
var Out = os.Stdout

// RefreshInterval in the default time duration to wait for refreshing the output
var RefreshInterval = time.Millisecond * 10

// defaultProgress is the default progress
var defaultProgress = New()

// getTerminalWidth returns the width of the terminal window.
// If detection fails, it returns the default Width value.
func getTerminalWidth() int {
	if width, _, err := term.GetSize(int(Out.Fd())); err == nil && width > 0 {
		return width
	}
	return Width
}

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
	lw := uilive.New()
	lw.Out = Out

	return &Progress{
		Width:           Width,
		Out:             Out,
		Bars:            make([]*Bar, 0),
		RefreshInterval: RefreshInterval,

		tdone: make(chan bool),
		lw:    lw,
		mtx:   &sync.RWMutex{},
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

func (p *Progress) SetOut(o io.Writer) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.Out = o
	p.lw.Out = o
}

func (p *Progress) SetRefreshInterval(interval time.Duration) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.RefreshInterval = interval
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
	for {

		p.mtx.Lock()
		interval := p.RefreshInterval
		p.mtx.Unlock()

		select {
		case <-time.After(interval):
			p.print()
		case <-p.tdone:
			p.print()
			close(p.tdone)
			return
		}
	}
}

func (p *Progress) print() {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// Auto-detect terminal width and adjust bar widths
	termWidth := getTerminalWidth()

	for _, bar := range p.Bars {
		// Calculate available width for the progress bar itself
		// Account for decorators by checking the current string length
		bar.mtx.RLock()
		prependFuncs := make([]DecoratorFunc, len(bar.prependFuncs))
		copy(prependFuncs, bar.prependFuncs)
		appendFuncs := make([]DecoratorFunc, len(bar.appendFuncs))
		copy(appendFuncs, bar.appendFuncs)
		bar.mtx.RUnlock()

		decoratorsWidth := 0

		// Calculate prepend decorators width
		for _, f := range prependFuncs {
			decoratorsWidth += len(f(bar)) + 1 // +1 for space
		}

		// Calculate append decorators width
		for _, f := range appendFuncs {
			decoratorsWidth += len(f(bar)) + 1 // +1 for space
		}

		// Set bar width to terminal width minus decorators
		// Ensure minimum width of 10 characters for the bar
		barWidth := termWidth - decoratorsWidth
		if barWidth < 10 {
			barWidth = 10
		}

		bar.mtx.Lock()
		bar.Width = barWidth
		bar.mtx.Unlock()

		fmt.Fprintln(p.lw, bar.String())
	}
	p.lw.Flush()
}

// Start starts the rendering the progress of progress bars. It listens for updates using `bar.Set(n)` and new bars when added using `AddBar`
func (p *Progress) Start() {
	go p.Listen()
}

// Stop stops listening
func (p *Progress) Stop() {
	p.tdone <- true
	<-p.tdone
}

// Bypass returns a writer which allows non-buffered data to be written to the underlying output
func (p *Progress) Bypass() io.Writer {
	return p.lw.Bypass()
}
