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

package backlog_test

import (
	"errors"
	"testing"
	"time"

	"github.com/bketelsen/backlog"
)

type user struct {
	Name      string
	Email     string
	CreatedAt time.Time
}

func (u *user) MarshalLog(kv backlog.KeyValue) error {
	kv.AddString("name", u.Name)
	kv.AddString("email", u.Email)
	kv.AddInt64("created_at", u.CreatedAt.UnixNano())
	return nil
}

var _jane = &user{
	Name:      "Jane Doe",
	Email:     "jane@test.com",
	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
}

func withBenchedLogger(b *testing.B, f func(backlog.Logger)) {
	logger := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func BenchmarkNoContext(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("No context.")
	})
}

func BenchmarkBoolField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Boolean.", backlog.Bool("foo", true))
	})
}

func BenchmarkFloat64Field(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Floating point.", backlog.Float64("foo", 3.14))
	})
}

func BenchmarkIntField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Integer.", backlog.Int("foo", 42))
	})
}

func BenchmarkInt64Field(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("64-bit integer.", backlog.Int64("foo", 42))
	})
}

func BenchmarkStringField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Strings.", backlog.String("foo", "bar"))
	})
}

func BenchmarkStringerField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Level.", backlog.Stringer("foo", backlog.InfoLevel))
	})
}

func BenchmarkTimeField(b *testing.B) {
	t := time.Unix(0, 0)
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Time.", backlog.Time("foo", t))
	})
}

func BenchmarkDurationField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Duration", backlog.Duration("foo", time.Second))
	})
}

func BenchmarkErrorField(b *testing.B) {
	err := errors.New("egad!")
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Error.", backlog.Error(err))
	})
}

func BenchmarkStackField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Error.", backlog.Stack())
	})
}

func BenchmarkMarshalerField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Arbitrary backlog.LogMarshaler.", backlog.Marshaler("user", _jane))
	})
}

func BenchmarkObjectField(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Reflection-based serialization.", backlog.Object("user", _jane))
	})
}

func BenchmarkAddCallerHook(b *testing.B) {
	logger := backlog.NewJSON(
		backlog.Output(backlog.Discard),
		backlog.AddCaller(),
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Caller.")
		}
	})
}

func Benchmark10Fields(b *testing.B) {
	withBenchedLogger(b, func(log backlog.Logger) {
		log.Info("Ten fields, passed at the log site.",
			backlog.Int("one", 1),
			backlog.Int("two", 2),
			backlog.Int("three", 3),
			backlog.Int("four", 4),
			backlog.Int("five", 5),
			backlog.Int("six", 6),
			backlog.Int("seven", 7),
			backlog.Int("eight", 8),
			backlog.Int("nine", 9),
			backlog.Int("ten", 10),
		)
	})
}

func Benchmark100Fields(b *testing.B) {
	const batchSize = 50
	logger := backlog.NewJSON(backlog.DebugLevel, backlog.Output(backlog.Discard))

	// Don't include allocating these helper slices in the benchmark. Since
	// access to them isn't synchronized, we can't run the benchmark in
	// parallel.
	first := make([]backlog.Field, batchSize)
	second := make([]backlog.Field, batchSize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for i := 0; i < batchSize; i++ {
			// We're duplicating keys, but that doesn't affect performance.
			first[i] = backlog.Int("foo", i)
			second[i] = backlog.Int("foo", i+batchSize)
		}
		logger.With(first...).Info("Child loggers with lots of context.", second...)
	}
}
