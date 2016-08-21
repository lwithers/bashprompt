package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
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
	User          *user.User
	Hostname, Cwd string
	LoadAvg       string
	Width         int

	ExitCode    = flag.String("exitCode", "0", "exit code of last process")
	ScreenWidth = flag.String("screenWidth", "80", "screen width in chars")
)

func main() {
	var err error
	flag.Parse()

	Width, _ = strconv.Atoi(*ScreenWidth)
	NewlineIfNecessary()

	User, err = user.Current()
	_ = err // TODO

	Hostname, err = os.Hostname()
	_ = err // TODO

	Cwd, err = os.Getwd()
	_ = err // TODO

	LoadAvg = "LOAD"
	statusFlag := "▲"
	if *ExitCode != "0" {
		statusFlag = "▼"
	}

	prefix := fmt.Sprintf("┌(%s)─(%s@%s)─(",
		time.Now().Format(time.RFC822),
		User.Username, Hostname)
	suffix := fmt.Sprintf("─(%s %s)",
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

// LoadAverage returns the 1-minute load average string, colour coded.
func LoadAverage() string {
	b, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return CSI + "41mERR" + CSI + "m"
	}

	p := bytes.IndexByte(b, ' ')
	if p == -1 {
	}
}

func InsideGitRepo(dir string) bool {
	for retry := 0; retry < 10; retry++ {
		if dir == "/" {
			return false
		}
		_, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil {
			return true
		}
		dir = filepath.Dir(dir)
	}
	return false
}

func GitBranch() string {
	c := exec.Command("/usr/bin/git", "describe", "--all",
		"--dirty") // TODO dirty code/colour
	b, err := c.Output()
	if err != nil {
		return err.Error() // TODO colour
	}
	for {
		p := bytes.IndexRune(b, '/')
		if p == -1 {
			break
		}
		b = b[p+1:]
	}
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return string(b)
}
