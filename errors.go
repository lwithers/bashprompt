package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrLog     *os.File
	ErrLogName string
)

func CaptureError(report error) {
	if report == nil {
		return
	}

	var err error

	if ErrLog == nil {
		ErrLogName = filepath.Join(os.Getenv("HOME"), ".bprompt_errors")
		os.Rename(ErrLogName, ErrLogName+".old")
		if ErrLog, err = os.Create(ErrLogName); err != nil {
			return
		}
	}

	pc := make([]uintptr, 1)
	runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	fmt.Fprintf(ErrLog, "%s (%s:%d):\n *** %v\n\n",
		frame.Function, frame.File, frame.Line, report)
}
