package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	CSI = "\x1B["
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
	fmt.Printf(strings.Repeat(" ", Width))

	statusFlag := CSI + "32m▲" + CSI + "m"
	if *ExitCode != "0" {
		statusFlag = CSI + "31m▼" + CSI + "m"
	}

	prefix := fmt.Sprintf("┌(%s)─(%s@%s)─(",
		time.Now().Format(time.RFC822),
		User.Username, Hostname)
	suffix := fmt.Sprintf("─(%.2f %s)",
		LoadAvg, statusFlag)
	reqLen := PrintableLength(prefix) + PrintableLength(suffix) + 1
	path := FitPath(Cwd, Width-reqLen)

	reqLen += utf8.RuneCountInString(path)
	if reqLen < Width {
		reqLen = Width - reqLen - 1
	} else {
		reqLen = 0
	}
	fmt.Printf("%s%s)%s%s\n", prefix, path,
		strings.Repeat("─", reqLen), suffix)

	branch := "n/a"
	if InsideGitRepo(Cwd) {
		branch = GitBranch()
	}

	fmt.Printf("└─(%s)─(%s) \\$ ",
		branch, filepath.Base(Cwd))
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
		strings.Repeat("·", Width-1))
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
