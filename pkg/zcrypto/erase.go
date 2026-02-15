package zcrypto

import "runtime"

// Erase zeroes out a byte slice to prevent sensitive data from lingering in memory.
//
// This is best-effort: Go's garbage collector may copy memory during compaction,
// so copies of the data may exist elsewhere in the heap. Zeroing the known
// location is still worthwhile — it reduces the window of exposure and clears
// the most obvious location an attacker would look.
func Erase(b []byte) {
	for i := range b {
		b[i] = 0
	}
	// prevent the compiler from optimizing away the zeroing loop
	// by ensuring the slice is considered reachable after the writes
	runtime.KeepAlive(&b)
}
