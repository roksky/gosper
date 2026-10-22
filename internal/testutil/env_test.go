package testutil

import (
	"reflect"
	"testing"
	"time"
)

func TestFixedClock(t *testing.T) {
	want := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
	clock := FixedClock(want)
	for i := 0; i < 3; i++ {
		if got := clock(); !got.Equal(want) {
			t.Fatalf("call %d = %v, want %v", i, got, want)
		}
	}
}

func TestUTCDefaults(t *testing.T) {
	if DefaultLocation != time.UTC {
		t.Fatal("DefaultLocation must be UTC")
	}
	if got := UTC(); got != time.UTC {
		t.Fatalf("UTC() = %v, want time.UTC", got)
	}
	if DefaultLocale == "" {
		t.Fatal("DefaultLocale must be set")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	want := []string{"a", "b", "c"}
	for i := 0; i < 5; i++ {
		if got := SortedKeys(m); !reflect.DeepEqual(got, want) {
			t.Fatalf("SortedKeys = %v, want %v", got, want)
		}
	}

	if got := SortedKeys[string](nil); len(got) != 0 {
		t.Fatalf("SortedKeys(nil) = %v, want empty", got)
	}
}

func TestSortedKeysConcurrent(t *testing.T) {
	m := map[string]int{"z": 1, "y": 2, "x": 3, "w": 4}
	done := make(chan struct{})
	for i := 0; i < 8; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for j := 0; j < 100; j++ {
				_ = SortedKeys(m)
			}
		}()
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}
