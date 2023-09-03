// MIT License, (c) Maragu ApS
// https://github.com/maragudk/gomponents/blob/main/internal/assert/assert.go

// Package assert provides testing helpers.
package assert

import (
	"testing"
)

func Equal(t *testing.T, expected any, actual any) {
	t.Helper()

	if expected != actual {
		t.Fatalf(`expected "%v" but got "%v"`, expected, actual)
	}
}

// Error checks for a non-nil error.
func Error(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("error is nil")
	}
}
