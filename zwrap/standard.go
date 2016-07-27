// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zwrap

import (
	"errors"
	"fmt"

	"github.com/bketelsen/backlog"
)

// ErrInvalidLevel indicates that the user chose an invalid Level when
// constructing a StandardLogger.
var ErrInvalidLevel = errors.New("StandardLogger's print level must be Debug, Info, Warn, or Error")

// StandardLogger is the interface exposed by the standard library's log.Logger.
type StandardLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
}

// Standardize wraps a Logger to make it compatible with the standard library.
// It takes the Logger itself, and the level to use for the StandardLogger's
// Print family of methods. If the specified Level isn't Debug, Info, Warn, or
// Error, Standardize returns ErrInvalidLevel.
func Standardize(l backlog.Logger, printAt backlog.Level) (StandardLogger, error) {
	s := stdLogger{
		panic: l.Panic,
		fatal: l.Fatal,
	}
	switch printAt {
	case backlog.DebugLevel:
		s.write = l.Debug
	case backlog.InfoLevel:
		s.write = l.Info
	case backlog.WarnLevel:
		s.write = l.Warn
	case backlog.ErrorLevel:
		s.write = l.Error
	default:
		return nil, ErrInvalidLevel
	}
	return &s, nil
}

type stdLogger struct {
	write func(string, ...backlog.Field)
	panic func(string, ...backlog.Field)
	fatal func(string, ...backlog.Field)
}

func (s *stdLogger) Print(args ...interface{}) {
	s.write(fmt.Sprint(args...))
}

func (s *stdLogger) Printf(format string, args ...interface{}) {
	s.write(fmt.Sprintf(format, args...))
}

func (s *stdLogger) Println(args ...interface{}) {
	// Don't use fmt.Sprintln, since the Logger will be wrapping this
	// message in an envelope.
	s.write(fmt.Sprint(args...))
}

func (s *stdLogger) Panic(args ...interface{}) {
	s.panic(fmt.Sprint(args...))
}

func (s *stdLogger) Panicf(format string, args ...interface{}) {
	s.panic(fmt.Sprintf(format, args...))
}

func (s *stdLogger) Panicln(args ...interface{}) {
	// See Println.
	s.panic(fmt.Sprint(args...))
}

func (s *stdLogger) Fatal(args ...interface{}) {
	s.fatal(fmt.Sprint(args...))
}

func (s *stdLogger) Fatalf(format string, args ...interface{}) {
	s.fatal(fmt.Sprintf(format, args...))
}

func (s *stdLogger) Fatalln(args ...interface{}) {
	// See Println.
	s.fatal(fmt.Sprint(args...))
}
