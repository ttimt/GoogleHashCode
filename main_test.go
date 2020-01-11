package main

import (
	"testing"
)

func TestReturnHello(t *testing.T) {
	want := "hello"
	get := returnHello()

	if want != get {
		t.Error("Want:", want, "Got:", get)
	}
}
