package toolkits_test

import (
	toolkits "github.com/shijl0925/go-toolkits"
	"testing"
)

func TestCast(t *testing.T) {
	t.Run("TC01: int number", func(t *testing.T) {
		var n any = 10                     // var n interface{} = 10
		got1, ok1 := toolkits.Cast[int](n) // 10, true
		if !ok1 {
			t.Errorf("expected success, got failure")
		}
		if got1 != 10 {
			t.Errorf("got %v (%T), expected %v (%T)", got1, got1, 10, 10)
		}

		got2, ok2 := toolkits.Cast[string](n) // "", false
		if ok2 {
			t.Errorf("expected failure, got success")
		}
		if got2 != "" {
			t.Errorf("got %v (%T), expected %v (%T)", got2, got2, "", "")
		}
	})

	t.Run("TC02: string number", func(t *testing.T) {
		var s any = "10"
		got1, ok1 := toolkits.Cast[int](s) // 0, false
		if ok1 {
			t.Errorf("expected failure, got success")
		}
		if got1 != 0 {
			t.Errorf("got %v (%T), expected %v (%T)", got1, got1, 0, 0)
		}

		got2, ok2 := toolkits.Cast[string](s) // "10", true
		if !ok2 {
			t.Errorf("expected success, got failure")
		}
		if got2 != "10" {
			t.Errorf("got %v (%T), expected %v (%T)", got2, got2, "", "")
		}
	})

	t.Run("TC03: pointer", func(t *testing.T) {
		var a = 10
		var p interface{} = &a

		got1, ok1 := toolkits.Cast[int](p) // 0, false
		if ok1 {
			t.Errorf("expected failure, got success")
		}
		if got1 != 0 {
			t.Errorf("got %v (%T), expected %v (%T)", got1, got1, 0, 0)
		}

		got2, ok2 := toolkits.Cast[*int](p) // 10, true
		if !ok2 {
			t.Errorf("expected success, got failure")
		}
		if *got2 != 10 {
			t.Errorf("got %v (%T), expected %v (%T)", *got2, *got2, 10, 10)
		}
	})

	t.Run("TC04: nil", func(t *testing.T) {
		var input interface{}
		got, ok := toolkits.Cast[int](input) // 0, false
		if ok {
			t.Errorf("expected failure, got success")
		}
		if got != 0 {
			t.Errorf("got %v (%T), expected %v (%T)", got, got, 0, 0)
		}
	})
}
