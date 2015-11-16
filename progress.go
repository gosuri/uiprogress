package uiprogress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gosuri/uilive"
)

var (
	Writer = os.Stdout
)

type Progress struct {
	Writer io.Writer
	Width  int
	Bars   []*Bar
}

func New() *Progress {
	return &Progress{
		Width:  100,
		Writer: Writer,
		Bars:   make([]*Bar, 0),
	}
}

func (p *Progress) AddBar(total int) *Bar {
	bar := NewBar(total)
	bar.Width = p.Width
	p.Bars = append(p.Bars, bar)
	return bar
}

func (p *Progress) Start() {
	lw := uilive.New()
	lw.Out = p.Writer
	go func() {
		for {
			for _, bar := range p.Bars {
				fmt.Fprintln(lw, bar.String())
			}
			time.Sleep(time.Millisecond * 10)
			lw.Flush()
		}
	}()
}
