package testutil

import (
	"runtime"
	"testing"
	"time"
)

// GoroutineCount captures the current number of goroutines for a leak check.
func GoroutineCount() int { return runtime.NumGoroutine() }

// AssertNoGoroutineGrowth waits briefly for cleanup and fails if the count
// remains above before+allowance. Use a small allowance for unrelated runtime
// activity; pair this with explicit resource-close assertions.
func AssertNoGoroutineGrowth(t testing.TB, before, allowance int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for {
		runtime.GC()
		count := runtime.NumGoroutine()
		if count <= before+allowance {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("goroutines grew from %d to %d (allowance %d)", before, count, allowance)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
