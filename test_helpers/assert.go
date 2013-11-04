package test_helpers

import (
	"testing"
)

func AssertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected <%s>, but got <%s>", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected, actual interface{}) {
	if expected == actual {
		t.Errorf("not expected <%s>, but got <%s>", expected, actual)
	}
}
