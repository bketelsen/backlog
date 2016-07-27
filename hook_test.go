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

package backlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
// todo fix regex
func TestHookAddCaller(t *testing.T) {
	buf := newTestBuffer()
	logger := NewJSON(DebugLevel, Output(buf), AddCaller())
	logger.Info("Callers.")

	re := regexp.MustCompile(`"msg":"hook_test.go:[\d]+ TestHookAddCaller(): Callers\.[\S]+"`)
	assert.Regexp(t, re, buf.String(), "Expected to find package name and file name in output.")
}
*/

func TestHookAddCallerFail(t *testing.T) {
	buf := newTestBuffer()
	errBuf := newTestBuffer()

	originalSkip := _callerSkip
	_callerSkip = 1e3
	defer func() { _callerSkip = originalSkip }()

	logger := NewJSON(DebugLevel, Output(buf), ErrorOutput(errBuf), AddCaller())
	logger.Info("Failure.")
	assert.Equal(t, "failed to get caller\n", errBuf.String(), "Didn't find expected failure message.")
	assert.Contains(t, buf.String(), `"msg":"Failure."`, "Expected original message to survive failures in runtime.Caller.")
}

func TestHookAddStacks(t *testing.T) {
	buf := newTestBuffer()
	logger := NewJSON(DebugLevel, Output(buf), AddStacks(InfoLevel))

	logger.Info("Stacks.")
	output := buf.String()
	require.Contains(t, output, "backlog.TestHookAddStacks", "Expected to find test function in stacktrace.")
	assert.Contains(t, output, `"stacktrace":`, "Stacktrace added under an unexpected key.")

	buf.Reset()
	logger.Warn("Stacks.")
	assert.Contains(t, buf.String(), `"stacktrace":`, "Expected to include stacktrace at Warn level.")

	buf.Reset()
	logger.Debug("No stacks.")
	assert.NotContains(t, buf.String(), "Unexpected stacktrace at Debug level.")
}
