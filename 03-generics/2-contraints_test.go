package generics

import "testing"

func Test_useMixed(t *testing.T) {
	var t1 T1
	var t2 T2
	var mi myInt

	_, _, _ = t1, t2, mi

	//useMixed(t1)
	//useMixed(t2)
	useMixed(mi)
}
