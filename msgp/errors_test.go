package msgp

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestWrapVanillaErrorWithNoAdditionalContext(t *testing.T) {
	err := errors.New("test")
	w := WrapError(err)
	if w == err {
		t.Fatal()
	}
	if w.Error() != err.Error() {
		t.Fatal()
	}
	if w.(errWrapped).Resumable() {
		t.Fatal()
	}
}

func TestWrapVanillaErrorWithAdditionalContext(t *testing.T) {
	err := errors.New("test")
	w := WrapError(err, "foo", "bar")
	if w == err {
		t.Fatal()
	}
	if w.Error() == err.Error() {
		t.Fatal()
	}
	if w.(Error).Resumable() {
		t.Fatal()
	}
	if !strings.HasPrefix(w.Error(), err.Error()) {
		t.Fatal()
	}
	rest := w.Error()[len(err.Error()):]
	if rest != " at foo/bar" {
		t.Fatal()
	}
}

func TestWrapResumableError(t *testing.T) {
	err := ArrayError{}
	w := WrapError(err)
	if !w.(Error).Resumable() {
		t.Fatal()
	}
}

func TestWrapMultiple(t *testing.T) {
	err := &TypeError{}
	w := WrapError(WrapError(err, "b"), "a")
	expected := `msgp: attempted to decode type "<invalid>" with method for "<invalid>" at a/b`
	if expected != w.Error() {
		t.Fatal()
	}
}

func TestCause(t *testing.T) {
	for idx, err := range []error{
		errors.New("test"),
		ArrayError{},
		&ErrUnsupportedType{},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			cerr := WrapError(err, "test")
			if cerr == err {
				t.Fatal()
			}
			if Cause(err) != err {
				t.Fatal()
			}
		})
	}
}

func TestCauseShortByte(t *testing.T) {
	err := ErrShortBytes
	cerr := WrapError(err, "test")
	if cerr != err {
		t.Fatal()
	}
	if Cause(err) != err {
		t.Fatal()
	}
}

func TestUnwrap(t *testing.T) {
	// check errors that get wrapped
	for idx, err := range []error{
		errors.New("test"),
		io.EOF,
	} {
		t.Run(fmt.Sprintf("wrapped_%d", idx), func(t *testing.T) {
			cerr := WrapError(err, "test")
			if cerr == err {
				t.Fatal()
			}
			uwerr := errors.Unwrap(cerr)
			if uwerr != err {
				t.Fatal()
			}
			if !errors.Is(cerr, err) {
				t.Fatal()
			}
		})
	}

	// check errors where only context is applied
	for idx, err := range []error{
		ArrayError{},
		&ErrUnsupportedType{},
	} {
		t.Run(fmt.Sprintf("ctx_only_%d", idx), func(t *testing.T) {
			cerr := WrapError(err, "test")
			if cerr == err {
				t.Fatal()
			}
			if errors.Unwrap(cerr) != nil {
				t.Fatal()
			}
		})
	}
}

func TestSimpleQuoteStr(t *testing.T) {
	// check some cases for simpleQuoteStr
	type tcase struct {
		in  string
		out string
	}
	tcaseList := []tcase{
		{
			in:  ``,
			out: `""`,
		},
		{
			in:  `abc`,
			out: `"abc"`,
		},
		{
			in:  `"`,
			out: `"\""`,
		},
		{
			in:  `'`,
			out: `"'"`,
		},
		{
			in:  `on🔥!`,
			out: `"on\xf0\x9f\x94\xa5!"`,
		},
		{
			in:  "line\r\nbr",
			out: `"line\r\nbr"`,
		},
		{
			in:  "\x00",
			out: `"\x00"`,
		},
		{ // invalid UTF-8 should make no difference but check it regardless
			in:  "not\x80valid",
			out: `"not\x80valid"`,
		},
	}

	for i, tc := range tcaseList {
		tc := tc
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			out := simpleQuoteStr(tc.in)
			if out != tc.out {
				t.Errorf("input %q; expected: %s; but got: %s", tc.in, tc.out, out)
			}
		})
	}
}
