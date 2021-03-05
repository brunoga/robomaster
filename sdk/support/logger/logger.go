package logger

import (
	"io"
	"log"
	"sync"
)

// Logger is a multi-level logger, supporting trace, info, warning and error
// levels that can be enabled/disabled independently.
type Logger struct {
	m       sync.Mutex
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

// New returns a new logger instance configured with the given destination for
// each level. Setting any destination to nil will disable the logging (log will
// be sent to ioutil.Discard). Common destinations are os.Stdout and os.Stderr.
func New(traceDest, infoDest, warningDest, errorDest io.Writer) *Logger {
	if traceDest == nil {
		traceDest = io.Discard
	}
	if infoDest == nil {
		infoDest = io.Discard
	}
	if warningDest == nil {
		warningDest = io.Discard
	}
	if errorDest == nil {
		errorDest = io.Discard
	}

	return &Logger{
		sync.Mutex{},
		createSdkLogger(traceDest, "TRACE: "),
		createSdkLogger(infoDest, "INFO: "),
		createSdkLogger(warningDest, "WARNING: "),
		createSdkLogger(errorDest, "ERROR: "),
	}
}

func (l *Logger) SetTraceDest(dest io.Writer) {
	l.m.Lock()
	l.trace = createSdkLogger(dest, "TRACE: ")
	l.m.Unlock()
}

func (l *Logger) SetInfoDest(dest io.Writer) {
	l.m.Lock()
	l.trace = createSdkLogger(dest, "INFO: ")
	l.m.Unlock()
}

func (l *Logger) SetWarningDest(dest io.Writer) {
	l.m.Lock()
	l.trace = createSdkLogger(dest, "WARNING: ")
	l.m.Unlock()
}

func (l *Logger) SetErrorDest(dest io.Writer) {
	l.m.Lock()
	l.trace = createSdkLogger(dest, "ERROR: ")
	l.m.Unlock()
}

// TRACE logs trace messages, used mostly for debugging.
func (l *Logger) TRACE(format string, a ...interface{}) {
	l.m.Lock()
	l.trace.Printf(format, a...)
	l.m.Unlock()
}

// INFO logs informational messages, used for reporting back informational
// messages that are generally useful.
func (l *Logger) INFO(format string, a ...interface{}) {
	l.m.Lock()
	l.info.Printf(format, a...)
	l.m.Unlock()
}

// WARNING logs warning messages, used to report anything that might be
// problematic but is not considered an error.
func (l *Logger) WARNING(format string, a ...interface{}) {
	l.m.Lock()
	l.warning.Printf(format, a...)
	l.m.Unlock()
}

// ERROR logs error messages, used to report actual issues that might need to
// be fixed.
func (l *Logger) ERROR(format string, a ...interface{}) {
	l.m.Lock()
	l.error.Printf(format, a...)
	l.m.Unlock()
}

func createSdkLogger(dest io.Writer, prefix string) *log.Logger {
	return log.New(dest, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}
