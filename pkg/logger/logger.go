package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/vandi37/vanerrors/vanstack"
)

type Logger struct {
	out        io.Writer
	datePrefix string
	prefixes   [5]string
}

var stdPrefixes = [5]string{
	"..debug>>",
	"^^info>>",
	"!!Warn>>",
	"##Error>>",
	"--Fatal>>",
}

var stdDate = "02.01.06 3:04:05"

func New(out io.Writer) *Logger {
	return NewWithSettings(out, stdDate, stdPrefixes)
}

func NewWithSettings(out io.Writer, date string, prefixes [5]string) *Logger {
	return &Logger{
		out:        out,
		datePrefix: date,
		prefixes:   prefixes,
	}
}

func (l *Logger) writeln(lvl int, a []any) (n int, err error) {
	return fmt.Fprintln(l.out, append([]any{time.Now().Format(l.datePrefix), l.prefixes[lvl]}, a...)...)
}

func (l *Logger) writef(lvl int, format string, a []any) (n int, err error) {
	format = "%s %s " + format + "\n"
	return fmt.Fprintf(l.out, format, append([]any{time.Now().Format(l.datePrefix), l.prefixes[lvl]}, a...)...)
}

func (l *Logger) Debugln(a ...any) (n int, err error) {
	return l.writeln(0, a)
}

func (l *Logger) Debugf(format string, a ...any) (n int, err error) {
	return l.writef(0, format, a)

}

func (l *Logger) Println(a ...any) (n int, err error) {
	return l.writeln(1, a)
}

func (l *Logger) Printf(format string, a ...any) (n int, err error) {
	return l.writef(1, format, a)
}

func (l *Logger) Warnln(a ...any) (n int, err error) {
	return l.writeln(2, a)
}

func (l *Logger) Warnf(format string, a ...any) (n int, err error) {
	return l.writef(2, format, a)
}

func (l *Logger) Errorln(a ...any) (n int, err error) {
	return l.writeln(3, a)
}

func (l *Logger) Errorf(format string, a ...any) (n int, err error) {
	return l.writef(3, format, a)
}

func (l *Logger) Fatalln(a ...any) (n int, err error) {
	l.writeln(4, a)
	stack := vanstack.NewStack()
	stack.Fill("", 20)
	n, err = fmt.Fprintln(os.Stderr, stack)
	os.Exit(http.StatusTeapot)
	return
}

func (l *Logger) Fatalf(format string, a ...any) (n int, err error) {
	l.writef(4, format, a)
	stack := vanstack.NewStack()
	stack.Fill("", 20)
	n, err = fmt.Fprintln(os.Stderr, stack)
	os.Exit(http.StatusTeapot)
	return
}
