package validation_test

import (
	"fmt"
	"io"

	charmlog "github.com/charmbracelet/log"
	"github.com/codeready-toolchain/argocd-checker/pkg/validation"
	"github.com/sanity-io/litter"
)

type TestLogger struct {
	*charmlog.Logger
	records map[charmlog.Level][]LogRecord
}

func NewTestLogger(w io.Writer, opts charmlog.Options) *TestLogger {
	return &TestLogger{
		Logger:  charmlog.NewWithOptions(w, opts),
		records: map[charmlog.Level][]LogRecord{},
	}
}

type LogRecord struct {
	Msg     any
	KeyVals []any
}

var _ fmt.Stringer = LogRecord{}

func (r LogRecord) String() string {
	return litter.Sdump(r)
}

func (l *TestLogger) Fatals() []LogRecord {
	return l.records[charmlog.FatalLevel]
}

func (l *TestLogger) Errors() []LogRecord {
	return l.records[charmlog.ErrorLevel]
}

func (l *TestLogger) Warnings() []LogRecord {
	return l.records[charmlog.WarnLevel]
}

var _ validation.Logger = &TestLogger{}

// Debug implements Logger.
func (l *TestLogger) Debug(msg any, keyvals ...any) {
	l.records[charmlog.DebugLevel] = append(l.records[charmlog.DebugLevel], LogRecord{
		Msg:     msg,
		KeyVals: keyvals,
	})
	l.Logger.Debug(msg, keyvals...)
}

// Info implements Logger.
func (l *TestLogger) Info(msg any, keyvals ...any) {
	l.records[charmlog.InfoLevel] = append(l.records[charmlog.InfoLevel], LogRecord{
		Msg:     msg,
		KeyVals: keyvals,
	})
	l.Logger.Info(msg, keyvals...)
}

// Warn implements Logger.
func (l *TestLogger) Warn(msg any, keyvals ...any) {
	l.records[charmlog.WarnLevel] = append(l.records[charmlog.WarnLevel], LogRecord{
		Msg:     msg,
		KeyVals: keyvals,
	})
	l.Logger.Warn(msg, keyvals...)
}

// Error implements Logger.
func (l *TestLogger) Error(msg any, keyvals ...any) {
	l.records[charmlog.ErrorLevel] = append(l.records[charmlog.ErrorLevel], LogRecord{
		Msg:     msg,
		KeyVals: keyvals,
	})
}

// Fatal implements Logger.
func (l *TestLogger) Fatal(msg any, keyvals ...any) {
	l.records[charmlog.FatalLevel] = append(l.records[charmlog.FatalLevel], LogRecord{
		Msg:     msg,
		KeyVals: keyvals,
	})
	l.Logger.Fatal(msg, keyvals...)
}
