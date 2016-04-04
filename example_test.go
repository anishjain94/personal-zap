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

package zap_test

import (
	"os"
	"time"

	"github.com/uber-common/zap"
)

func Example() {
	// Log in JSON, using zap's reflection-free JSON encoder.
	logger := zap.NewJSON(zap.Info, os.Stdout)
	// For repeatable tests, pretend that it's always 1970.
	logger.StubTime()

	logger.Warn("Log without structured data...")
	logger.Warn(
		"Or use strongly-typed wrappers to add structured context.",
		zap.String("library", "zap"),
		zap.Duration("latency", time.Nanosecond),
	)

	// Avoid re-serializing the same data repeatedly by creating a child logger
	// with some attached context. That context is added to all the child's
	// log output, but doesn't affect the parent.
	child := logger.With(zap.String("user", "jane@test.com"), zap.Int("visits", 42))
	child.Error("Oh no!")

	// To reduce allocations, fields are returned to a sync.Pool immediately
	// after use. To safely re-use fields, pass the Keep option when
	// constructing them.
	fields := []zap.Field{zap.Int("one", 1, zap.Keep), zap.Int("two", 2, zap.Keep)}
	logger.Info("Using fields once is always safe.", fields...)
	logger.Info("Because we passed Keep, it's safe to re-use our fields.", fields...)

	// Output:
	// {"msg":"Log without structured data...","level":"warn","ts":0,"fields":{}}
	// {"msg":"Or use strongly-typed wrappers to add structured context.","level":"warn","ts":0,"fields":{"library":"zap","latency":1}}
	// {"msg":"Oh no!","level":"error","ts":0,"fields":{"user":"jane@test.com","visits":42}}
	// {"msg":"Using fields once is always safe.","level":"info","ts":0,"fields":{"one":1,"two":2}}
	// {"msg":"Because we passed Keep, it's safe to re-use our fields.","level":"info","ts":0,"fields":{"one":1,"two":2}}
}

func ExampleNest() {
	logger := zap.NewJSON(zap.Info, os.Stdout)
	// Stub the current time in tests.
	logger.StubTime()

	// We'd like the logging context to be {"outer":{"inner":42}}
	logger.Debug("Nesting context.", zap.Nest("outer",
		zap.Int("inner", 42),
	))

	// If we want to stop a field from being returned to a sync.Pool on use,
	// use Keep.
	nest := zap.Nest("outer", zap.Int("inner", 42))
	zap.Keep(nest)
	logger.Info("The first use is always safe.", nest)
	logger.Info("Since we called Keep, re-use is safe.", nest)

	// Output:
	// {"msg":"The first use is always safe.","level":"info","ts":0,"fields":{"outer":{"inner":42}}}
	// {"msg":"Since we called Keep, re-use is safe.","level":"info","ts":0,"fields":{"outer":{"inner":42}}}
}
