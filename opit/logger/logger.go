/*
Package logger implements leveled logging utilities. It provides log levels that
have syslog compatibilities.
*/
package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Level defines type of level.
type Level int

const (
	// PanicLevel represents the process caused panic.
	PanicLevel Level = iota

	// FatalLevel represents an un-handleable error. The process exits immediately
	// with exit code 1.
	FatalLevel

	// ErrorLevel represents an error but handleable.
	ErrorLevel

	// WarningLevel represents a warning. Warning level is not critical, but it
	// can potentially cause oddities.
	WarningLevel

	// NoticeLevel is defined for syslog compatibilities. Don't use this level.
	NoticeLevel

	// InfoLevel is used for general information.
	InfoLevel

	// DebugLevel is used for debugging. The message will not print without
	// verbosity option.
	DebugLevel
)

type Output interface {
	WriteEntry(*Entry)
}

// Entry represents an logging entry. It contains created time, log level, and
// message.
type Entry struct {
	Time    time.Time
	Level   Level
	Message string
}

type Logger struct {
	mu   sync.Mutex
	outs []Output
}

// New creates a new Logger.
func New() *Logger {
	return &Logger{
		mu:   sync.Mutex{},
		outs: []Output{StdoutOutput},
	}
}

func (l *Logger) output(lv Level, m string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := &Entry{
		Time:    time.Now(),
		Level:   lv,
		Message: m,
	}

	for _, out := range l.outs {
		out.WriteEntry(e)
	}
}

// AddOutput adds the output destination.
func (l *Logger) AddOutput(out Output) {
	l.outs = append(l.outs, out)
}

// RemoveOutput removes the registered output destination.
func (l *Logger) RemoveOutput(out Output) {
	for i, o := range l.outs {
		if o == out {
			if i == 0 {
				l.outs = l.outs[1:]
			} else if i < (len(l.outs) - 1) {
				l.outs = append(l.outs[:i-1], l.outs[i+1:]...)
			} else {
				l.outs = l.outs[:i-1]
			}
			break
		}
	}
}

// Debug writes message at debug level. The arguments are handled in the manner
// of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	l.output(DebugLevel, fmt.Sprint(v...))
}

// Debugf writes message at debug level. The arguments are handled in the manner
// of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(DebugLevel, fmt.Sprintf(format, v...))
}

// Debugln writes message at debug level. The arguments are handled in the
// manner of fmt.Println.
func (l *Logger) Debugln(v ...interface{}) {
	l.output(DebugLevel, fmt.Sprintln(v...))
}

// Error writes message at error level. The arguments are handled in the manner
// of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.output(ErrorLevel, fmt.Sprint(v...))
}

// Errorf writes message at error level. The arguments are handled in the manner
// of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(ErrorLevel, fmt.Sprintf(format, v...))
}

// Errorln writes message at error level. The arguments are handled in the
// manner of fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	l.output(ErrorLevel, fmt.Sprintln(v...))
}

// Fatal writes message at fatal level. The arguments are handled in the manner
// of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.output(FatalLevel, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf writes message at fatal level. The arguments are handled in the manner
// of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.output(FatalLevel, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln writes message at fatal level. The arguments are handled in the
// manner of fmt.Println.
func (l *Logger) Fatalln(v ...interface{}) {
	l.output(FatalLevel, fmt.Sprintln(v...))
	os.Exit(1)
}

// Info writes message at info level. The arguments are handled in the manner
// of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.output(InfoLevel, fmt.Sprint(v...))
}

// Infof writes message at info level. The arguments are handled in the manner
// of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(InfoLevel, fmt.Sprintf(format, v...))
}

// Infoln writes message at info level. The arguments are handled in the manner
// of fmt.Println.
func (l *Logger) Infoln(v ...interface{}) {
	l.output(InfoLevel, fmt.Sprintln(v...))
}

// Warning writes message at warning level. The arguments are handled in the
// manner of fmt.Print.
func (l *Logger) Warning(v ...interface{}) {
	l.output(WarningLevel, fmt.Sprint(v...))
}

// Warningf writes message at warning level. The arguments are handled in the
// manner of fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.output(WarningLevel, fmt.Sprintf(format, v...))
}

// Warningln writes message at warning level. The arguments are handled in the
// manner of fmt.Println.
func (l *Logger) Warningln(v ...interface{}) {
	l.output(WarningLevel, fmt.Sprintln(v...))
}

// std is the standard logger.
var std = New()

// AddOutput adds the output destination for the standard logger.
func AddOutput(out Output) {
	std.AddOutput(out)
}

// RemoveOutput adds the output destination from the standard logger.
func RemoveOutput(out Output) {
	std.RemoveOutput(out)
}

// Debug writes message to standard logger at debug level. The arguments are
// handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugf writes message to standard logger at debug level. The arguments are
// handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// Debugln writes message to standard logger at debug level. The arguments are
// handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	std.Debugln(v...)
}

// Error writes message to standard logger at error level. The arguments are
// handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorf writes message to standard logger at error level. The arguments are
// handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Errorln writes message to standard logger at error level. The arguments are
// handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	std.Errorln(v...)
}

// Fatal writes message to standard logger at fatal level. The arguments are
// handled in the manner of fmt.Print.
func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

// Fatalf writes message to standard logger at fatal level. The arguments are
// handled in the manner of fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

// Fatalln writes message to standard logger at fatal level. The arguments are
// handled in the manner of fmt.Println.
func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

// Info writes message to standard logger at info level. The arguments are
// handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infof writes message to standard logger at info level. The arguments are
// handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Infoln writes message to standard logger at info level. The arguments are
// handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	std.Infoln(v...)
}

// Warning writes message to standard logger at warning level. The arguments are
// handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	std.Warning(v...)
}

// Warningf writes message to standard logger at warning level. The arguments
// are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

// Warningln writes message to standard logger at warning level. The arguments
// are handled in the manner of fmt.Println.
func Warningln(v ...interface{}) {
	std.Warningln(v...)
}
