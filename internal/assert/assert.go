package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	// marks this function as Test Helper Function
	t.Helper()

	if actual != expected {
		t.Errorf("'%v' should be '%v'", actual, expected)
	}
}

// vim: ts=4 sw=4 fdm=indent
