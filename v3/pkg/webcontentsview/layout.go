package webcontentsview

import (
	"math"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type LayoutAnchor int

const (
	LayoutAnchorStart LayoutAnchor = iota
	LayoutAnchorEnd
)

// LayoutValue resolves to an absolute value using a fixed pixel component plus
// an optional percentage of the parent dimension.
type LayoutValue struct {
	Pixels  int
	Percent float64
}

func (v LayoutValue) resolve(total int) int {
	return v.Pixels + int(math.Round((v.Percent/100)*float64(total)))
}

// Layout describes a WebContentsView rectangle relative to its parent.
//
// X and Y are offsets from the chosen anchor:
//   - LayoutAnchorStart means left/top
//   - LayoutAnchorEnd means right/bottom
//
// Width and Height support fixed pixels, percentages, or a mix of both.
type Layout struct {
	X                LayoutValue
	Y                LayoutValue
	Width            LayoutValue
	Height           LayoutValue
	HorizontalAnchor LayoutAnchor
	VerticalAnchor   LayoutAnchor
}

func (l Layout) Resolve(parent application.Rect) application.Rect {
	width := l.Width.resolve(parent.Width)
	height := l.Height.resolve(parent.Height)
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

	x := l.X.resolve(parent.Width)
	y := l.Y.resolve(parent.Height)

	if l.HorizontalAnchor == LayoutAnchorEnd {
		x = parent.X + parent.Width - width - x
	} else {
		x += parent.X
	}

	if l.VerticalAnchor == LayoutAnchorEnd {
		y = parent.Y + parent.Height - height - y
	} else {
		y += parent.Y
	}

	return application.Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}
