package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type Level int

const (
	MsgSpam = Level(iota)
	MsgNoPrefix
	MsgInfo
	MsgDebug
	MsgWarning
	MsgError
	MsgPanic
	MsgFatal
	MsgNeverDisplay
)

func GetPrefix(level Level) string {
	now := time.Now()
	var file string
	var line int
	var ok bool
	_, file, line, ok = runtime.Caller(3)
	if !ok {
		file = "???"
		line = 0
	}

	var out = ""

	out += now.Local().Format("2006/01/02 15:04:05")
	out += " "
	out += filepath.Base(file)
	out += ":"
	out += strconv.Itoa(line)
	out += ": "
	switch level {
	case MsgSpam:
		return "[SPAM] " + out
	case MsgInfo:
		return "[INFO] " + out
	case MsgDebug:
		return "[DEBUG] " + out
	case MsgWarning:
		return "[WARNING] " + out
	case MsgError:
		return "[ERROR] " + out
	case MsgFatal:
		return "[FATAL] " + out
	case MsgPanic:
		return "[PANIC] " + out
	default:
		return ""
	}
}

type Logger struct {
	logger       *log.Logger
	GetPrefix    func(Level) string
	displayLevel Level
}

var std = CreateLogger(os.Stderr)

func CreateLogger(writer io.Writer) Logger {
	var out Logger = Logger{logger: log.New(writer, "[INFO] ", 0), GetPrefix: GetPrefix, displayLevel: MsgFatal}
	return out
}

func (l *Logger) Output(level Level, calldepth int, s string) error {
	if level < l.displayLevel {
		return nil
	}
	l.logger.SetPrefix(l.GetPrefix(level))
	return l.logger.Output(calldepth+1, s)
}

func (l *Logger) Print(level Level, v ...any) {
	if level < l.displayLevel {
		return
	}
	l.logger.SetPrefix(l.GetPrefix(level))
	if level == MsgFatal {
		l.logger.Fatal(v...)
	}
	if level == MsgPanic {
		l.logger.Panic(v...)
	}
	l.logger.Print(v...)
}

func (l *Logger) Fatal(v ...any) {
	l.Print(MsgFatal, v...)
}

func (l *Logger) Panic(v ...any) {
	l.Print(MsgPanic, v...)
}

func (l *Logger) Printf(level Level, format string, v ...any) {
	if level < l.displayLevel {
		return
	}
	l.logger.SetPrefix(l.GetPrefix(level))
	if level == MsgFatal {
		l.logger.Fatalf(format, v...)
	}
	if level == MsgPanic {
		l.logger.Panicf(format, v...)
	}
	l.logger.Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.Printf(MsgFatal, format, v...)
}

func (l *Logger) Panicf(format string, v ...any) {
	l.Printf(MsgPanic, format, v...)
}

func (l *Logger) Println(level Level, v ...any) {
	if level < l.displayLevel {
		return
	}
	l.logger.SetPrefix(l.GetPrefix(level))
	if level == MsgFatal {
		l.logger.Fatalln(v...)
	}
	if level == MsgPanic {
		l.logger.Panicln(v...)
	}
	l.logger.Println(v...)
}

func (l *Logger) Fatalln(v ...any) {
	l.Println(MsgFatal, v...)
}

func (l *Logger) Panicln(v ...any) {
	l.Println(MsgPanic, v...)
}

func (l *Logger) Writer() io.Writer {
	return l.logger.Writer()
}

func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

func (l *Logger) DisplayLevel() Level {
	return l.displayLevel
}

func (l *Logger) SetDisplay(level Level) {
	l.displayLevel = level
}

//-------------------------------------------------------

func Output(level Level, calldepth int, s string) error {
	return std.Output(level, calldepth+1, s) // +1 for this frame.
}

func Print(level Level, v ...any) {
	std.Print(level, v...)
}

func Fatal(v ...any) {
	std.Print(MsgFatal, v...)
}

func Panic(v ...any) {
	std.Print(MsgPanic, v...)
}

func Printf(level Level, format string, v ...any) {
	std.Printf(level, format, v...)
}

func Fatalf(format string, v ...any) {
	std.Printf(MsgFatal, format, v...)
}

func Panicf(format string, v ...any) {
	std.Printf(MsgPanic, format, v...)
}

func Println(level Level, v ...any) {
	std.Println(level, v...)
}

func Fatalln(v ...any) {
	std.Println(MsgFatal, v...)
}

func Panicln(v ...any) {
	std.Println(MsgPanic, v...)
}

func Writer() io.Writer {
	return std.Writer()
}

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func DisplayLevel() Level {
	return std.displayLevel
}

func SetDisplay(level Level) {
	std.displayLevel = level
}
