package uiprogress

import (
	"bytes"
	"fmt"
	"regexp"
	"time"
)

var (
	Fill     byte = '='
	Head     byte = '>'
	Empty    byte = '-'
	LeftEnd  byte = '['
	RightEnd byte = ']'
)

type Bar struct {
	Total   int
	Current int

	// LeftEnd is character in the left most part of the progress indicator. Defaults to '['
	LeftEnd byte

	// RightEnd is character in the right most part of the progress indicator. Defaults to ']'
	RightEnd byte

	// Fill is the character representing completed progress. Defaults to '='
	Fill byte

	// Head is the character that moves when progress is updated.  Defaults to '>'
	Head byte

	// Empty is the character that represents the empty progress. Default is '-'
	Empty byte

	Width int

	appendFuncs []func() string

	prepends []string
}

func NewBar(total int) *Bar {
	return &Bar{
		Total:    total,
		Width:    100,
		LeftEnd:  LeftEnd,
		RightEnd: RightEnd,
		Head:     Head,
		Fill:     Fill,
		Empty:    Empty,
	}
}

func (b *Bar) AppendFunc(f func() string) *Bar {
	b.appendFuncs = append(b.appendFuncs, f)
	return b
}

func (b *Bar) AppendCompleted() *Bar {
	f := func() string { return fmt.Sprintf(" %3.f %%", b.CompletedPercent()) }
	b.AppendFunc(f)
	return b
}

func (b *Bar) AppendElapsed() *Bar {
	return b
}

func (b *Bar) Prepend(p string) *Bar {
	b.prepends = append(b.prepends, p)
	return b
}

func (b *Bar) Bytes() []byte {
	completedWidthF := float64(b.Width) * (b.CompletedPercent() / 100.00)
	completedWidth := int(completedWidthF)

	// add fill and empty bits
	var buf bytes.Buffer
	for i := 0; i < completedWidth; i++ {
		buf.WriteByte(b.Fill)
	}
	for i := 0; i < b.Width-completedWidth; i++ {
		buf.WriteByte(b.Empty)
	}

	// add head bit
	pb := buf.Bytes()
	if completedWidth > 0 && completedWidth < b.Width {
		pb[completedWidth-1] = b.Head
	}

	// add left and right ends bits
	pb[0], pb[len(pb)-1] = b.LeftEnd, b.RightEnd

	for _, f := range b.appendFuncs {
		pb = append(pb, []byte(f())...)
	}
	return pb
}

func (b *Bar) String() string {
	return string(b.Bytes())
}

func (b *Bar) CompletedPercent() float64 {
	return (float64(b.Current) / float64(b.Total)) * 100.00
}

func prettyTime(t time.Duration) string {
	re, err := regexp.Compile(`(\d+).(\d+)(\w+)`)
	if err != nil {
		return err.Error()
	}
	parts := re.FindSubmatch([]byte(t.String()))
	if len(parts) != 4 {
		return "---"
	}
	return string(parts[1]) + string(parts[3])
}
