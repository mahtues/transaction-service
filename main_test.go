package main

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t.Run("successful test", func(t *testing.T) {
		expected := 5

		actual := foo(2, 3)

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("failed test", func(t *testing.T) {
		t.Errorf("forced failure")
	})
}
