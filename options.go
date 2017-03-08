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

package zap

import "go.uber.org/zap/zapcore"

// Option is used to set options for the logger.
type Option interface {
	apply(*Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

// WrapCore wraps or replaces the Logger's underlying zapcore.Core.
func WrapCore(f func(zapcore.Core) zapcore.Core) Option {
	return optionFunc(func(log *Logger) {
		log.core = f(log.core)
	})
}

// Hooks registers functions which will be called each time the Logger writes
// out an Entry. Repeated use of Hooks is additive.
//
// See zapcore.RegisterHooks for details.
func Hooks(hooks ...func(zapcore.Entry) error) Option {
	return optionFunc(func(log *Logger) {
		log.core = zapcore.RegisterHooks(log.core, hooks...)
	})
}

// Fields adds fields to the Logger.
func Fields(fs ...zapcore.Field) Option {
	return optionFunc(func(log *Logger) {
		log.core = log.core.With(fs)
	})
}

// ErrorOutput sets the destination for errors generated by the logger. The
// supplied Pusher is automatically wrapped with a mutex, so it need not be
// safe for concurrent use.
func ErrorOutput(p zapcore.Pusher) Option {
	return optionFunc(func(log *Logger) {
		log.errorOutput = zapcore.Lock(p)
	})
}

// Development puts the logger in development mode, which alters the behavior
// of the DPanic method.
func Development() Option {
	return optionFunc(func(log *Logger) {
		log.development = true
	})
}

// AddCaller configures the Logger to annotate each message with the filename
// and line number of zap's caller.
func AddCaller() Option {
	return optionFunc(func(log *Logger) {
		log.addCaller = true
	})
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option).
func AddCallerSkip(skip int) Option {
	return optionFunc(func(log *Logger) {
		log.callerSkip += skip
	})
}

// AddStacktrace configures the Logger to record a stack trace for all messages at
// or above a given level.
func AddStacktrace(lvl zapcore.LevelEnabler) Option {
	return optionFunc(func(log *Logger) {
		log.addStack = lvl
	})
}
