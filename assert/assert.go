package assert

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func Equal(t *testing.T, expected, actual interface{}, msg string, args ...interface{}) {
	if equal(expected, actual) {
		return
	}

	t.Errorf("%s: expected: %#v\n actual: %#v", fmt.Sprintf(msg, args...), expected, actual)
}

func NoError(t *testing.T, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}

	t.Errorf("%s: unexpected error: %v", fmt.Sprintf(msg, args...), err)
}

func ErrorIs(t *testing.T, err, target error, msg string, args ...interface{}) {
	if errors.Is(err, target) {
		return
	}

	t.Errorf("%s: error: %v is not error: %v", fmt.Sprintf(msg, args...), err, target)
}

func equal(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
