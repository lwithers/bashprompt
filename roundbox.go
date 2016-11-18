package main

import (
	"fmt"
	"strings"
)

const (
	FirstLine  = -1
	SecondLine = -2

	RoundBoxPadding = 2
	RoundBoxSpacing = 2
)

type RoundBoxInfo struct {
	leftColour, rightColour int
	text                    string
}

func RoundBox(s string) *RoundBoxInfo {
	return &RoundBoxInfo{
		leftColour:  7,
		rightColour: 7,
		text:        s,
	}
}

func (r *RoundBoxInfo) SetColour(leftColour, rightColour int) {
	if leftColour >= 0 && leftColour <= 7 {
		r.leftColour = leftColour
	}
	if rightColour >= 0 && rightColour <= 7 {
		r.rightColour = rightColour
	}
}

func (r *RoundBoxInfo) Width() int {
	return PrintableLength(r.text) + RoundBoxPadding
}

func (r *RoundBoxInfo) Draw(w TerminalWriter, leading int) {
	SetColour(w, "37")
	switch leading {
	case FirstLine:
		w.WriteString("┌")
	case SecondLine:
		w.WriteString("└─")
	default:
		w.WriteString(strings.Repeat("─", leading))
	}

	SetColour(w, fmt.Sprintf("3%d", r.leftColour))
	w.WriteRune('') // \uE0B6 (powerline extra symbol)
	SetColour(w, fmt.Sprintf("0;4%d", r.leftColour))
	w.WriteString(r.text)
	SetColour(w, fmt.Sprintf("0;3%d", r.rightColour))
	w.WriteRune('') // \uE0B4 (powerline extra symbol)
}

func PrintLine(w TerminalWriter, line int, boxes []*RoundBoxInfo) {
	var leading int

	for i := range boxes {
		switch {
		case i == 0:
			leading = line
		case line == FirstLine && i == len(boxes)-2:
			leading = RoundBoxSpacing + RemainingWidth(line, boxes)
		default:
			leading = RoundBoxSpacing
		}
		boxes[i].Draw(w, leading)
	}
}

func RemainingWidth(line int, boxes []*RoundBoxInfo) int {
	r := Width // how many chars we can print

	// subtract length of lead string
	switch line {
	case FirstLine:
		r -= 1
	case SecondLine:
		r -= 2
	}

	// subtract width of each box (includes the () decorator but not any
	// inter-box space)
	for i := range boxes {
		r -= boxes[i].Width()
	}

	// remove two chars per inter-box space
	r -= RoundBoxSpacing * (len(boxes) - 1)
	return r
}
