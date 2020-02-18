package main

import "testing"

func TestReplaceDotAndUppercase(t *testing.T) {
	got := replaceDotAndUppercase("FOO.COM")
	want := "FOO_COM"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
