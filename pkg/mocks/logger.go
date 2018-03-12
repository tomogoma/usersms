package mocks

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tomogoma/usersms/pkg/logging"
)

type Entry struct {
	Level string
	Fmt   *string // nil if logging without formatting e.g. Error() instead of Errorf()
	Args  []interface{}
}

type Logger struct {
	Fields   map[string]interface{}
	Logs     []Entry
	Spinoffs []*Logger
	HTTPReq  *http.Request
}

const (
	LevelInfo  = "Info"
	LevelWarn  = "Warn"
	LevelError = "Error"
	LevelFatal = "Fatal"
)

func (lg *Logger) WithFields(f map[string]interface{}) logging.Logger {
	newLG := lg.copy()
	for k, v := range f {
		newLG.Fields[k] = v
	}
	return newLG
}

func (lg *Logger) WithField(k string, v interface{}) logging.Logger {
	newLG := lg.copy()
	newLG.Fields[k] = v
	return newLG
}

func (lg *Logger) WithHTTPRequest(r *http.Request) logging.Logger {
	newLG := lg.copy()
	newLG.HTTPReq = r
	return newLG
}

func (lg *Logger) Infof(fmt string, args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelInfo, Fmt: &fmt, Args: args})
}

func (lg *Logger) Warnf(fmt string, args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelWarn, Fmt: &fmt, Args: args})
}

func (lg *Logger) Errorf(fmt string, args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelError, Fmt: &fmt, Args: args})
}

func (lg *Logger) Info(args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelInfo, Args: args})
}

func (lg *Logger) Warn(args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelWarn, Args: args})
}

func (lg *Logger) Error(args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelError, Args: args})
}

func (lg *Logger) Fatal(args ...interface{}) {
	lg.prep()
	lg.Logs = append(lg.Logs, Entry{Level: LevelFatal, Args: args})
}

func (lg *Logger) copy() *Logger {
	lg.prep()
	newLG := &Logger{Fields: lg.Fields}
	lg.Spinoffs = append(lg.Spinoffs, newLG)
	return newLG
}

func (lg *Logger) prep() {
	if lg.Logs == nil {
		lg.Logs = make([]Entry, 0)
	}
	if lg.Fields == nil {
		lg.Fields = make(map[string]interface{})
	}
	if lg.Spinoffs == nil {
		lg.Spinoffs = make([]*Logger, 0)
	}
}

func (lg *Logger) PrintLogs(t *testing.T) {
	for _, so := range lg.Spinoffs {
		so.PrintLogs(t)
	}
	for _, e := range lg.Logs {
		var payload string
		if e.Fmt != nil {
			payload = fmt.Sprintf(*e.Fmt, e.Args...)
		} else {
			payload = fmt.Sprint(e.Args...)
		}
		t.Logf("%s: %s\n", e.Level, payload)
	}
}
