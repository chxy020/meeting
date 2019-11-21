package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

// NewProcessRegexp compile process name with possible prefix paths
func NewProcessRegexp(format string, v ...interface{}) *regexp.Regexp {
	if strings.HasPrefix(format, "/") {
		return regexp.MustCompile(fmt.Sprintf("^"+format+"$", v...))
	}
	return regexp.MustCompile(fmt.Sprintf("^(/usr/local/bin/|/usr/bin/|/bin/|)"+format+"$", v...))
}

// NewProcessNameRegexp init regexp with process name
func NewProcessNameRegexp(name string) *regexp.Regexp {
	return NewProcessRegexp("%s(\\s+.*|)", name)
}

// RunCommandWay can be set to execute os.exec
type RunCommandWay func(timeout time.Duration, command interface{}, args ...string) (bytes.Buffer, RunError)

func realRun(timeout time.Duration, command interface{}, args ...string) (out bytes.Buffer, err RunError) {
	if c, ok := command.(string); ok {
		var errout bytes.Buffer
		cmd := exec.Command(c, args...)
		cmd.Stdout = &out
		cmd.Stderr = &errout
		isTimeout, e := runWithTimeout(cmd, timeout)
		if e != nil {
			err = TryGetExitError(e, errout)
			return out, err
		}
		if isTimeout {
			err = NewRunError(int(syscall.ETIMEDOUT), "process timeout")
		}
	} else if cmd, ok := command.(*exec.Cmd); ok {
		isTimeout, e := runWithTimeout(cmd, timeout)
		if e != nil {
			err = TryGetExitError(e, *bytes.NewBuffer([]byte{}))
			return out, err
		}
		if isTimeout {
			err = NewRunError(int(syscall.ETIMEDOUT), "process timeout")
		}
	} else {
		err = NewRunError(-1, "no command found")
		return out, err
	}
	return
}

func runWithTimeout(cmd *exec.Cmd, timeout time.Duration) (isTimeout bool, err error) {

	if timeout <= 0 {
		// no need to wait timeout
		return false, cmd.Run()
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		go func() {
			<-done // release wait goroutine
		}()
		return true, err
	case err = <-done:
		return false, err
	}
}

// RunCommand is actual variable will be called to execute the need
var RunCommand RunCommandWay = realRun

// SetRunCommand can mock the os.exec
func SetRunCommand(c RunCommandWay) {
	if c == nil {
		RunCommand = realRun
	} else {
		RunCommand = c
	}
}

// RunError is the wrapper of os.exec error, it export error code
type RunError interface {
	ExitCode() int
	Error() string
	ExitCodeEquals(syscall.Errno) bool
	ExitCodeIn(...syscall.Errno) bool
}

// NewRunError is used to get the RunError
func NewRunError(code int, msg string) (err RunError) {
	err = &runErrorImpl{code, msg}
	return
}

type runErrorImpl struct {
	code int
	msg  string
}

func (e runErrorImpl) Error() string {
	return e.msg
}

func (e runErrorImpl) ExitCode() int {
	return e.code
}

func (e runErrorImpl) ExitCodeEquals(errno syscall.Errno) bool {
	return e.code == int(errno)
}

func (e runErrorImpl) ExitCodeIn(errnos ...syscall.Errno) bool {
	for _, errno := range errnos {
		if e.code == int(errno) {
			return true
		}
	}
	return false
}

// TryGetExitError convert cmd error into run error
func TryGetExitError(err error, out bytes.Buffer) RunError {
	if err == nil {
		return nil
	}
	if msg, ok := err.(*exec.ExitError); ok {
		return &runErrorImpl{int(msg.Sys().(syscall.WaitStatus).ExitStatus()), fmt.Sprintf("%s\n%s", err, &out)}
	}
	return &runErrorImpl{-1, fmt.Sprintf("%s\n%s", err, &out)}
}

// Func defines run func type
type Func func() error

// ErrIntervalContinue return ErrIntervalContinue if you want to continue the interval
var ErrIntervalContinue = errors.New("continue")

// RunFuncWithIntervalTimeout run func with interval and wait for timeout
// return ErrIntervalContinue if you want to retry,
func RunFuncWithIntervalTimeout(f Func, interval time.Duration, timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done(): // time out
			return context.DeadlineExceeded
		default:
			err = f()
			if err != ErrIntervalContinue {
				return
			}
			time.Sleep(interval)
		}
	}
}

// RunFuncWithInterval infinitely run func with interval
// return ErrIntervalContinue if you want to retry,
func RunFuncWithInterval(f Func, interval time.Duration) (err error) {
	for range time.Tick(interval) {
		if err = f(); err == ErrIntervalContinue {
			continue
		}
		return
	}
	return
}
