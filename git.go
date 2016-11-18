package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func GitBox() *RoundBoxInfo {
	if !InsideGitRepo(Cwd) {
		return RoundBox(" — ")
	}

	branch := GitBranch()
	dirty := strings.HasSuffix(branch, "-dirty")
	if dirty {
		branch = branch[:len(branch)-6]
	}

	var b bytes.Buffer
	SetColour(&b, "1;30")
	b.WriteRune('') // \ue0a0 git branch marker
	SetColour(&b, "0;30;44")
	b.WriteRune(' ')
	b.WriteString(branch)
	if dirty {
		//SetColour(&b, "35")
		//b.WriteRune('') // \ue0ba triangle separator
		SetColour(&b, "35")
		b.WriteRune('\ue0b2') // \ue0ba triangle separator
		SetColour(&b, "30;45")
		b.WriteString("dirty")
	}
	r := RoundBox(b.String())
	r.SetColour(4, 4)
	if dirty {
		r.SetColour(4, 5)
	}
	return r
}
