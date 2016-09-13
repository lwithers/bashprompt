package main

import "strings"

type RoundBox string

type StringWriter interface {
	WriteString(s string) (int, error)
}

const (
	FirstLine  = -1
	SecondLine = -2

	RoundBoxPadding = 2
)

func (r RoundBox) Width() int {
	return PrintableLength(string(r)) + RoundBoxPadding
}

func (r RoundBox) Draw(w StringWriter, leading int) {
	w.WriteString(CSI + "37m")
	switch leading {
	case FirstLine:
		w.WriteString("┌")
	case SecondLine:
		w.WriteString("└─")
	default:
		w.WriteString(strings.Repeat("─", leading))
	}

	w.WriteString("(")
	w.WriteString(CSI + "m")
	w.WriteString(string(r))
	w.WriteString(CSI + "0m" + CSI + "37m" + ")")
}

func PrintLine(w StringWriter, line int, boxes []RoundBox) {
	var leading int

	for i := range boxes {
		switch {
		case i == 0:
			leading = line
		case line == FirstLine && i == len(boxes)-2:
			leading = 1 + RemainingWidth(line, boxes)
		default:
			leading = 1
		}
		boxes[i].Draw(w, leading)
	}
}

func RemainingWidth(line int, boxes []RoundBox) int {
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

	// remove one char per inter-box space
	r -= len(boxes) - 1
	return r
}
