package zcrypto

import "runtime"

// eraser is an indirect function variable so the compiler cannot prove the
// clear call is a dead store — it can't inline across the variable indirection.
var eraser = func(b []byte) {
	clear(b)
}

// Erase zeroes out a byte slice to prevent sensitive data from lingering in memory.
//
// This is best-effort: Go's garbage collector may copy memory during compaction,
// so copies of the data may exist elsewhere in the heap. Zeroing the known
// location is still worthwhile — it reduces the window of exposure and clears
// the most obvious location an attacker would look.
func Erase(b []byte) {
	eraser(b)
	runtime.KeepAlive(&b)
}
