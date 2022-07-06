package main

import "testing"

func testGetArea(t *testing.T) {

	got := getArea(3, 2)
	want := 6

	if got != want {
		t.Errorf("Expected %q, received %q", got, want)
	}
}
