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

package benchmarks

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bketelsen/backlog"
	"github.com/bketelsen/backlog/zwrap"
)

var errExample = errors.New("fail")

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u user) MarshalLog(kv backlog.KeyValue) error {
	kv.AddString("name", u.Name)
	kv.AddString("email", u.Email)
	kv.AddInt64("created_at", u.CreatedAt.UnixNano())
	return nil
}

var _jane = user{
	Name:      "Jane Doe",
	Email:     "jane@test.com",
	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
}

func fakeFields() []backlog.Field {
	return []backlog.Field{
		backlog.Int("int", 1),
		backlog.Int64("int64", 2),
		backlog.Float64("float", 3.0),
		backlog.String("string", "four!"),
		backlog.Bool("bool", true),
		backlog.Time("time", time.Unix(0, 0)),
		backlog.Error(errExample),
		backlog.Duration("duration", time.Second),
		backlog.Marshaler("user-defined type", _jane),
		backlog.String("another string", "done!"),
	}
}

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func BenchmarkbacklogDisabledLevelsWithoutFields(b *testing.B) {
	logger := backlog.NewJSON(backlog.ErrorLevel, backlog.Output(backlog.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Should be discarded.")
		}
	})
}

func BenchmarkbacklogDisabledLevelsAccumulatedContext(b *testing.B) {
	context := fakeFields()
	logger := backlog.NewJSON(backlog.ErrorLevel, backlog.Output(backlog.Discard), backlog.Fields(context...))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Should be discarded.")
		}
	})
}

func BenchmarkbacklogDisabledLevelsAddingFields(b *testing.B) {
	logger := backlog.NewJSON(backlog.ErrorLevel, backlog.Output(backlog.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Should be discarded.", fakeFields()...)
		}
	})
}

func BenchmarkbacklogDisabledLevelsCheckAddingFields(b *testing.B) {
	logger := backlog.NewJSON(backlog.ErrorLevel, backlog.Output(backlog.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if m := logger.Check(backlog.InfoLevel, "Should be discarded."); m.OK() {
				m.Write(fakeFields()...)
			}
		}
	})
}

func BenchmarkbacklogAddingFields(b *testing.B) {
	logger := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.", fakeFields()...)
		}
	})
}

func BenchmarkbacklogWithAccumulatedContext(b *testing.B) {
	context := fakeFields()
	logger := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard), backlog.Fields(context...))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go really fast.")
		}
	})
}

func BenchmarkbacklogWithoutFields(b *testing.B) {
	logger := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.")
		}
	})
}

func BenchmarkbacklogSampleWithoutFields(b *testing.B) {
	messages := fakeMessages(1000)
	base := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	logger := zwrap.Sample(base, time.Second, 10, 10000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			logger.Info(messages[i%1000])
		}
	})
}

func BenchmarkbacklogSampleAddingFields(b *testing.B) {
	messages := fakeMessages(1000)
	base := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	logger := zwrap.Sample(base, time.Second, 10, 10000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			logger.Info(messages[i%1000], fakeFields()...)
		}
	})
}

func BenchmarkbacklogSampleCheckWithoutFields(b *testing.B) {
	messages := fakeMessages(1000)
	base := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	logger := zwrap.Sample(base, time.Second, 10, 10000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			if cm := logger.Check(backlog.InfoLevel, messages[i%1000]); cm.OK() {
				cm.Write()
			}
		}
	})
}

func BenchmarkbacklogSampleCheckAddingFields(b *testing.B) {
	messages := fakeMessages(1000)
	base := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	logger := zwrap.Sample(base, time.Second, 10, 10000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			if m := logger.Check(backlog.InfoLevel, messages[i%1000]); m.OK() {
				m.Write(fakeFields()...)
			}
		}
	})
}
