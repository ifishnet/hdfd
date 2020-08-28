// Copyright (c) 2014 The ifishnet developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdfjson_test

import (
	"testing"

	"github.com/ifishnet/hdfd/hdfjson"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hdfjson.ErrorCode
		want string
	}{
		{hdfjson.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{hdfjson.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{hdfjson.ErrInvalidType, "ErrInvalidType"},
		{hdfjson.ErrEmbeddedType, "ErrEmbeddedType"},
		{hdfjson.ErrUnexportedField, "ErrUnexportedField"},
		{hdfjson.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{hdfjson.ErrNonOptionalField, "ErrNonOptionalField"},
		{hdfjson.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{hdfjson.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{hdfjson.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{hdfjson.ErrNumParams, "ErrNumParams"},
		{hdfjson.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}

	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(hdfjson.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the Error type.
func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   hdfjson.Error
		want string
	}{
		{
			hdfjson.Error{Description: "some error"},
			"some error",
		},
		{
			hdfjson.Error{Description: "human-readable error"},
			"human-readable error",
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
