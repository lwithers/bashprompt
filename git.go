package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"unicode/utf8"
)

const (
	GitMaxBranchLen = 20
)

func GitBox() *RoundBoxInfo {
	gitDir := FindGitDir(Cwd)
	if gitDir == "" {
		return RoundBox(" — ")
	}

	branch, branchMode, err := FindGitBranch(gitDir)
	if err != nil {
		return gitBoxError(err)
	}

	dirty := CheckGitDirty(gitDir)

	var (
		b           bytes.Buffer
		leftColour  int = 4
		rightColour int = 4
	)

	if branchMode == GitModeDetachedHead {
		leftColour = 3
		SetColour(&b, "1;30;43")
		b.Write([]byte(" ")) // \uf05e crossed-out circle
		SetColour(&b, "0;34;43")
		b.WriteRune('') // \ue0b2 vim powerline separator
		SetColour(&b, "30;44")
		b.WriteString(branch[:7])
	} else {
		SetColour(&b, "1;30;44")
		if branchMode == GitModeTag {
			b.WriteRune('') // \uf005 starred
		} else {
			b.WriteRune('') // \ue0a0 git branch marker
		}
		SetColour(&b, "0;30;44")
		b.WriteRune(' ')

		if utf8.RuneCountInString(branch) > GitMaxBranchLen {
			short := make([]rune, 0, GitMaxBranchLen)
			for i, r := range branch {
				if i == GitMaxBranchLen-1 {
					break
				}
				short = append(short, r)
			}
			short = append(short, '…')
			branch = string(short)
		}
		b.WriteString(branch)
	}

	if dirty {
		rightColour = 5
		SetColour(&b, "35;44")
		b.WriteRune('') // \ue0b2 vim powerline separator
		SetColour(&b, "30;45")
		b.Write([]byte("dirty"))
	}

	r := RoundBox(b.String())
	r.SetColour(leftColour, rightColour)
	return r
}

// FindGitDir finds the path to the ".git" directory if we are inside a git
// repository, returning the path on success or an empty string on failure.
func FindGitDir(dir string) string {
	for depth := 0; depth < 10; depth++ {
		if dir == "/" {
			return ""
		}
		gitDir := filepath.Join(dir, ".git")
		if fi, err := os.Stat(gitDir); err == nil {
			if fi.IsDir() {
				return gitDir
			}
			if b, err := ioutil.ReadFile(gitDir); err == nil {
				if bytes.HasPrefix(b, []byte("gitdir: ")) {
					return filepath.Clean(filepath.Join(dir,
						string(bytes.TrimSpace(b[8:]))))
				}
			}
		}
		dir = filepath.Dir(dir)
	}
	return ""
}

const (
	GitModeUnknown = iota
	GitModeBranch
	GitModeTag
	GitModeDetachedHead
)

// FindGitBranch reads the branch of the given git directory, returning either
// the branch name, the tag name, or the detached commit.
func FindGitBranch(gitDir string) (string, int, error) {
	b, err := ioutil.ReadFile(filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return "", 0, err
	}
	b = bytes.TrimSpace(b)

	switch {
	case bytes.HasPrefix(b, []byte("ref: ")):
		b = b[5:]
		bb := bytes.SplitN(b, []byte("/"), 3)
		if len(bb) != 3 || !bytes.Equal(bb[0], []byte("refs")) {
			return "", 0, errors.New(".git/HEAD bad ref")
		}
		mode := GitModeBranch
		if bytes.Equal(bb[1], []byte("tags")) {
			mode = GitModeTag
		}
		return string(bb[2]), mode, nil

	case len(b) == 40:
		return string(b), GitModeDetachedHead, nil
	default:
		return "", 0, errors.New(".git/HEAD unparseable")
	}
}

// CheckGitDirty tests whether the given git repository is dirty or not.
func CheckGitDirty(gitDir string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Env = append(os.Environ(), "GIT_DIR="+gitDir)
	cmd.Dir = filepath.Dir(gitDir)
	op, _ := cmd.Output()
	return len(op) != 0
}

func gitBoxError(err error) *RoundBoxInfo {
	var b bytes.Buffer
	SetColour(&b, "1;30")
	b.WriteRune('') // \ue0a0 git branch marker
	SetColour(&b, "0;30;41")
	b.WriteRune(' ')
	b.WriteString(err.Error())

	r := RoundBox(b.String())
	r.SetColour(1, 1)
	return r
}
