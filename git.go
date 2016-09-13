package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
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
