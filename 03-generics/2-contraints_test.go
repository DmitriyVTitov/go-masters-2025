package generics

import (
	"reflect"
	"testing"
)

func Test_add(t *testing.T) {
	got := addNumbers(2, 3)
	if got != 5 {
		t.Fatalf("addNumbers(2, 3) = %v, want %v", got, 5)
	}

	got2 := addNumbers(2.0, 3.0)
	if got2 != 5.0 {
		t.Fatalf("addNumbers(2.0, 3.0) = %v, want %v", got2, 5.0)
	}

	got3 := addIntegers(int8(2), int8(3))
	if got3 != int8(5) {
		t.Fatalf("addIntegers(int8(2), int8(3)) = %v, want %v", got3, int8(5))
	}

	got4 := addIntegers(int64(2), int64(3))
	if got4 != int64(5) {
		t.Fatalf("addIntegers(int64(2), int64(3)) = %v, want %v", got4, int64(5))
	}
	if reflect.TypeOf(got4).Kind() != reflect.Int64 {
		t.Fatalf("addIntegers(int64(2), int64(3)) type = %v, want %v", reflect.TypeOf(got4).Kind(), reflect.Int64)
	}
}

func Test_useMixed(t *testing.T) {
	var t1 T1
	var t2 T2
	var mi myInt

	_, _, _ = t1, t2, mi

	//useMixed(t1)
	//useMixed(t2)
	useMixed(mi)
}
