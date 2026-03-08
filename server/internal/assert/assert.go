package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expected string) {
	t.Helper()

	if !strings.Contains(actual, expected) {
		t.Errorf("expected %q to contain %q", actual, expected)
	}
}

func NilError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("got: %v; expected: nil", err)
	}
}
