package uiprogress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gosuri/uilive"
)

var (
	Out             = os.Stdout
	DefaultProgress = New()
	RefreshInterval = time.Millisecond
)

type Progress struct {
	Out             io.Writer
	Width           int
	Bars            []*Bar
	RefreshInterval time.Duration
}

func New() *Progress {
	return &Progress{
		Width:           Width,
		Out:             Out,
		Bars:            make([]*Bar, 0),
		RefreshInterval: RefreshInterval,
	}
}

func AddBar(total int) *Bar {
	return DefaultProgress.AddBar(total)
}

func Start() {
	DefaultProgress.Start()
}

func (p *Progress) AddBar(total int) *Bar {
	bar := NewBar(total)
	bar.Width = p.Width
	p.Bars = append(p.Bars, bar)
	return bar
}

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
