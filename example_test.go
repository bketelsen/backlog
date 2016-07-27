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
	"encoding/json"
	"fmt"
	"time"

	"github.com/bketelsen/backlog"
)

func Example() {
	// Log in JSON, using backlog's reflection-free JSON encoder.
	// The default options will log any Info or higher logs to standard out.
	logger := backlog.NewJSON()
	// For repeatable tests, pretend that it's always 1970.
	logger.StubTime()

	logger.Warn("Log without structured data...")
	logger.Warn(
		"Or use strongly-typed wrappers to add structured context.",
		backlog.String("library", "backlog"),
		backlog.Duration("latency", time.Nanosecond),
	)

	// Avoid re-serializing the same data repeatedly by creating a child logger
	// with some attached context. That context is added to all the child's
	// log output, but doesn't affect the parent.
	child := logger.With(backlog.String("user", "jane@test.com"), backlog.Int("visits", 42))
	child.Error("Oh no!")

	// Output:
	// {"msg":"Log without structured data...","level":"warn","ts":0,"fields":{}}
	// {"msg":"Or use strongly-typed wrappers to add structured context.","level":"warn","ts":0,"fields":{"library":"backlog","latency":1}}
	// {"msg":"Oh no!","level":"error","ts":0,"fields":{"user":"jane@test.com","visits":42}}
}

func ExampleNest() {
	logger := backlog.NewJSON()
	// Stub the current time in tests.
	logger.StubTime()

	// We'd like the logging context to be {"outer":{"inner":42}}
	nest := backlog.Nest("outer", backlog.Int("inner", 42))
	logger.Info("Logging a nested field.", nest)

	// Output:
	// {"msg":"Logging a nested field.","level":"info","ts":0,"fields":{"outer":{"inner":42}}}
}

func ExampleNewJSON() {
	// The default logger outputs to standard out and only writes logs that are
	// Info level or higher.
	logger := backlog.NewJSON()
	// Stub the current time in tests.
	logger.StubTime()

	// The default logger does not print Debug logs.
	logger.Debug("This won't be printed.")
	logger.Info("This is an info log.")

	// Output:
	// {"msg":"This is an info log.","level":"info","ts":0,"fields":{}}
}

func ExampleNewJSON_options() {
	// We can pass multiple options to the NewJSON method to configure
	// the logging level, output location, or even the initial context.
	logger := backlog.NewJSON(
		backlog.DebugLevel,
		backlog.Fields(backlog.Int("count", 1)),
	)
	// Stub the current time in tests.
	logger.StubTime()

	logger.Debug("This is a debug log.")
	logger.Info("This is an info log.")

	// Output:
	// {"msg":"This is a debug log.","level":"debug","ts":0,"fields":{"count":1}}
	// {"msg":"This is an info log.","level":"info","ts":0,"fields":{"count":1}}
}

func ExampleCheckedMessage() {
	logger := backlog.NewJSON()
	// Stub the current time in tests.
	logger.StubTime()

	// By default, the debug logging level is disabled. However, calls to
	// logger.Debug will still allocate a slice to hold any passed fields.
	// Particularly performance-sensitive applications can avoid paying this
	// penalty by using checked messages.
	if cm := logger.Check(backlog.DebugLevel, "This is a debug log."); cm.OK() {
		// Debug-level logging is disabled, so we won't get here.
		cm.Write(backlog.Int("foo", 42), backlog.Stack())
	}

	if cm := logger.Check(backlog.InfoLevel, "This is an info log."); cm.OK() {
		// Since info-level logging is enabled, we expect to write out this message.
		cm.Write()
	}

	// Output:
	// {"msg":"This is an info log.","level":"info","ts":0,"fields":{}}
}

func ExampleLevel_MarshalText() {
	level := backlog.ErrorLevel
	s := struct {
		Level *backlog.Level `json:"level"`
	}{&level}
	bytes, _ := json.Marshal(s)
	fmt.Println(string(bytes))

	// Output:
	// {"level":"error"}
}

func ExampleLevel_UnmarshalText() {
	var s struct {
		Level backlog.Level `json:"level"`
	}
	// The zero value for a backlog.Level is backlog.InfoLevel.
	fmt.Println(s.Level)

	json.Unmarshal([]byte(`{"level":"error"}`), &s)
	fmt.Println(s.Level)

	// Output:
	// info
	// error
}
