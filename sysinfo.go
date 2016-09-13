package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
)

var (
	User          *user.User
	IsRoot        bool
	Hostname, Cwd string
	IsLocalhost   bool
	LoadAvg       float32
)

// GetUser finds the current user's details, saving them in User.
func GetUser() {
	var err error
	User, err = user.Current()
	if err != nil {
		CaptureError(err)
		User = &user.User{
			Uid:      "(err)",
			Gid:      "(err)",
			Username: "(err)",
			Name:     "(err)",
			HomeDir:  "/tmp",
		}
	}
	IsRoot = User.Uid == "0"
}

// GetHost finds the current hostname, and tests for the SSH_CONNECTION
// environment variable to determine whether or not this is a local or remote
// connection. Sets Hostname and IsLocalhost.
func GetHost() {
	var err error

	Hostname, err = os.Hostname()
	if err != nil {
		CaptureError(err)
		Hostname = "(err)"
	}

	IsLocalhost = (os.Getenv("SSH_CONNECTION") == "")
}

// GetLoadAverage finds the system's 1-minute load average, saving it in
// LoadAvg.
func GetLoadAverage() {
	LoadAvg = -1 // error condition

	b, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		CaptureError(err)
		return
	}

	p := bytes.IndexByte(b, ' ')
	if p == -1 {
		p = len(b)
	}

	var load float64
	if load, err = strconv.ParseFloat(string(b[:p]), 32); err != nil {
		CaptureError(err)
		return
	}

	LoadAvg = float32(load)
}

// GetCwd finds the current working directory, saving it in Cwd.
func GetCwd() {
	var err error
	Cwd, err = os.Getwd()
	if err != nil {
		Cwd = "/"
		CaptureError(err)
	}
}
