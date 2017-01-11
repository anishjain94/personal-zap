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

package multierror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorSliceString(t *testing.T) {
	tests := []struct {
		errs     errSlice
		expected string
	}{
		{nil, ""},
		{errSlice{}, ""},
		{errSlice{errors.New("foo")}, "foo"},
		{errSlice{errors.New("foo"), errors.New("bar")}, "foo; bar"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.errs.Error(), "Unexpected output from Error method.")
	}
}

func TestMultiErrorAsError(t *testing.T) {
	assert.Nil(t, (*Error)(nil).AsError(), "Expected calling AsError on nil to return nil.")
	assert.Nil(t, (&Error{}).AsError(), "Expected calling AsError with no accumulated errors to return nil.")

	e := errors.New("foo")
	assert.Equal(
		t,
		e,
		(&Error{errSlice{e}}).AsError(),
		"Expected AsError with single error to return the original error.",
	)

	m := &Error{errSlice{errors.New("foo"), errors.New("bar")}}
	assert.Equal(t, m.errs, m.AsError(), "Unexpected AsError output with multiple errors.")
}

func TestErrorAppend(t *testing.T) {
	foo := errors.New("foo")
	bar := errors.New("bar")
	for _, base := range []*Error{nil, {}} {
		base = base.Append(nil).Append(foo).Append(nil).Append(bar)
		assert.Equal(t, errSlice{foo, bar}, base.errs, "Collected errors don't match expectations.")
	}
}
