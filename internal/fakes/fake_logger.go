package fakes

import (
	"fmt"
	"io"

	"github.com/apex/log"
)

type fakeLog struct {
	log.Logger
	w io.Writer
}

// NewFakeLogger create a logger to capture output for testing purposes.
func NewFakeLogger(w io.Writer) *fakeLog {
	f := &fakeLog{
		w: w,
	}
	f.Logger.Handler = f
	f.Logger.Level = log.DebugLevel
	return f
}

func (f *fakeLog) HandleLog(e *log.Entry) error {
	if e.Level == log.ErrorLevel {
		_, _ = fmt.Fprintf(f.w, "ERROR: %s\n", e.Message)
		return nil
	}
	_, _ = fmt.Fprintln(f.w, e.Message)
	return nil
}

func (f *fakeLog) Writer() io.Writer {
	return f.w
}
