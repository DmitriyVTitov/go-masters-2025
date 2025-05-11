package main

import (
	"os"
	"testing"
)

func TestWriteObject(t *testing.T) {
	var s String = "ABC"
	var a = Album{Title: "Nevermind", Year: 1991}

	w := os.Stdout

	err := WriteObject(w, s, a)
	if err != nil {
		t.Fatal(err)
	}
}
