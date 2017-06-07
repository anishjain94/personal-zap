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

package log

import "github.com/uber-go/zap"

var singletonLogger zap.Logger

func init() {
	ConfigureStandard(zap.NewJSONEncoder())
}

// Standard returns the standard, singleton logger. Although this behavior is discouraged in many settings,
// many projects use it regardless and we'd prefer a reference implementation to each service doing it differently.
func Standard() zap.Logger {
	return singletonLogger
}

// ConfigureStandard configures the singleton logger. Note that this does not provide any
// synchronization guarantees -- it is up to you to configure your logger before calling Standard()
// and using it if you want to customize the options.
func ConfigureStandard(enc zap.Encoder, options ...zap.Option) zap.Logger {
	singletonLogger = zap.New(enc, options...)
	return singletonLogger
}
