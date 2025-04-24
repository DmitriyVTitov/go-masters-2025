package concurrency

import (
	"testing"
)

func Test_nonBlockingWrite(t *testing.T) {
	nonBlockingWrite()
}

func Test_waitGroup(t *testing.T) {
	waitGroup()
}

func Test_sema(t *testing.T) {
	sema()
}
