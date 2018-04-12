// Copyright (c) 2017 Uber Technologies, Inc.
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

package zaptest

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger builds a new Logger that logs all messages to the given
// testing.TB.
//
// Use this with a *testing.T or *testing.B to get logs which get printed only
// if a test fails or if you ran go test -v.
func NewLogger(t TestingT) *zap.Logger {
	return NewLoggerAt(t, zapcore.DebugLevel)
}

// NewLoggerAt builds a new Logger that logs messages to the given testing.TB
// if the given LevelEnabler allows it.
//
// Use this with a *testing.T or *testing.B to get logs which get printed only
// if a test fails or if you ran go test -v.
func NewLoggerAt(t TestingT, enab zapcore.LevelEnabler) *zap.Logger {
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			// EncoderConfig copied from zap.NewDevelopmentEncoderConfig.
			// Can't use it directly because that would cause a cyclic import.
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		testingWriter{t},
		enab,
	))
}

// testingWriter is a WriteSyncer that writes to the given testing.TB.
type testingWriter struct{ t TestingT }

func (w testingWriter) Write(p []byte) (n int, err error) {
	s := string(p)

	// Strip trailing newline because t.Log always adds one.
	if s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}

	// Note: t.Log is safe for concurrent use.
	w.t.Logf("%s", s)
	return len(p), nil
}

func (w testingWriter) Sync() error {
	return nil
}
