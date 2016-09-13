package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	CSI        = "\x1B["
	TimeFormat = "Mon Jan _2 15:04:05 MST"
)

var (
	Width int

	ExitCode    = flag.String("exitCode", "0", "exit code of last process")
	ScreenWidth = flag.String("screenWidth", "", "screen width in chars")
)

func main() {
	flag.Parse()

	GetUser()
	GetHost()
	GetLoadAverage()
	GetCwd()

	GetWidth()
	NewlineIfNecessary()

	statusFlag := CSI + "32m▲" + CSI + "m"
	if *ExitCode != "0" {
		statusFlag = CSI + "31m▼" + CSI + "m"
	}

	var secondLine []RoundBox
	firstLine := []RoundBox{
		Time(),
		Who(),
		"", // placeholder for truncated directory
		LoadAverage(),
		RoundBox(statusFlag),
	}

	path := FitPath(Cwd, RemainingWidth(FirstLine, firstLine))
	firstLine[2] = RoundBox(path)

	branch := "n/a"
	if InsideGitRepo(Cwd) {
		branch = GitBranch()
	}
	secondLine = append(secondLine,
		RoundBox(branch),
		RoundBox(filepath.Base(Cwd)),
	)

	buf := bytes.NewBuffer(nil)
	PrintLine(buf, FirstLine, firstLine)
	buf.WriteRune('\n')
	PrintLine(buf, SecondLine, secondLine)
	if IsRoot {
		buf.WriteString(CSI + "31m # " + CSI + "m")
	} else {
		buf.WriteString(CSI + "32m $ " + CSI + "m")
	}

	os.Stdout.Write(buf.Bytes())
}

// GetWidth interprets the screen width command line parameter, saving it in
// Width.
func GetWidth() {
	Width = 80 // failsafe value

	if *ScreenWidth == "" {
		CaptureError(errors.New("-screenWidth not specified"))
	}

	w, err := strconv.Atoi(*ScreenWidth)
	if err != nil {
		CaptureError(err)
		return
	}
	if w < 2 {
		CaptureError(fmt.Errorf("width too small: %d", w))
		return
	}

	Width = w
}

// NewlineIfNecessary determines if a newline is required, and arranges for it
// to be present if so. This is actually done by printing out a line's worth of
// filler characters, then resetting the cursor position to the left-hand side.
// If no newline was required, then the fillers will be overwritten by our
// normal output; otherwise, the fillers on the starting line (with partial
// output from the previous command) will be left in place.
func NewlineIfNecessary() {
	fmt.Printf(CSI+"37m%s"+CSI+"G"+CSI+"m",
		strings.Repeat("·", Width))
}

// Time returns a RoundBox formatted with the system time.
func Time() RoundBox {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(CSI + "30m" + CSI + "1m")
	buf.WriteString(time.Now().Format(TimeFormat))
	return RoundBox(buf.String())
}

// Who returns a RoundBox formatted with the username and, if not the local
// system, hostname.
func Who() RoundBox {
	buf := bytes.NewBuffer(nil)
	if IsRoot {
		buf.WriteString(CSI + "31m")
	} else {
		buf.WriteString(CSI + "32m")
	}
	buf.WriteString(User.Username)
	buf.WriteString(CSI + "34m@")

	if IsLocalhost {
		buf.WriteString(CSI + "32m")
	} else {
		buf.WriteString(CSI + "31m")
	}
	buf.WriteString(Hostname)

	return RoundBox(buf.String())
}

// LoadAverage returns a RoundBox displaying the colour-coded load average.
func LoadAverage() RoundBox {
	buf := bytes.NewBuffer(nil)
	switch {
	case LoadAvg < 0.2: // green
		buf.WriteString(CSI + "32m")
	case LoadAvg < 2.0: // yellow
		buf.WriteString(CSI + "33m")
	default: // red
		buf.WriteString(CSI + "31m")
	}

	fmt.Fprintf(buf, "%.2f", LoadAvg)

	return RoundBox(buf.String())
}

// PrintableLength returns the length of the printable characters in the given
// string. It will strip out ANSI escape codes.
func PrintableLength(str string) int {
	var (
		length   int
		inEscape bool
	)
	for _, r := range str {
		switch {
		case inEscape:
			if r == 'm' {
				inEscape = false
			}
		case r == 27:
			inEscape = true
		default:
			length++
		}
	}
	return length
}

func FitPath(p string, max int) string {
	if utf8.RuneCountInString(p) <= max {
		return p
	}

	components := strings.Split(p, "/")
	for i := 1; i < len(components); i++ {
		components[i-1] = "…/"
		p := filepath.Join(components[i-1:]...)
		if len(p) <= max {
			return p
		}
	}
	return "…"
}
