package fill

import (
	"context"
	"errors"
	"github.com/roksky/gosper/internal/datasource"
	"github.com/roksky/gosper/internal/expression"
	"github.com/roksky/gosper/internal/limits"
	"github.com/roksky/gosper/internal/testutil"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentFillsSharedTemplate(t *testing.T) {
	program, err := expression.Compile("3")
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &Template{Fields: []string{"n"}, Parameters: []ParameterDefault{{Name: "base", Default: program}}}
	source := datasource.FromSlice([]map[string]int{{"n": 1}, {"n": 2}, {"n": 3}})
	const workers = 32
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := Run(context.Background(), tmpl, source, limits.Budget{MaxRecords: 3})
			if err != nil {
				t.Errorf("Run: %v", err)
				return
			}
			if len(result.Records) != 3 || result.Parameters["base"] != expression.Int(3) {
				t.Errorf("result=%+v", result)
				return
			}
			for j, rec := range result.Records {
				if v, _ := rec.Variable("REPORT_COUNT"); v != expression.Int(j+1) {
					t.Errorf("record %d count=%v", j, v)
				}
			}
		}()
	}
	wg.Wait()
}

type blockingRecord map[string]any

func (r blockingRecord) Value(name string) (any, bool) { v, ok := r[name]; return v, ok }

type blockingCursor struct {
	entered       chan struct{}
	release       chan struct{}
	ignoreContext bool
	calls         int
	closes        atomic.Int32
}

func (c *blockingCursor) Next(ctx context.Context) (datasource.Record, bool, error) {
	c.calls++
	if c.calls == 1 {
		return blockingRecord{"n": 1}, true, nil
	}
	if c.calls == 2 {
		close(c.entered)
		if c.ignoreContext {
			<-c.release
		} else {
			select {
			case <-ctx.Done():
				return nil, false, ctx.Err()
			case <-c.release:
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	return nil, false, nil
}
func (c *blockingCursor) Close() error { c.closes.Add(1); return nil }

type blockingSource struct{ cursor *blockingCursor }

func (s blockingSource) Open(context.Context) (datasource.Cursor, error) { return s.cursor, nil }
func TestCancellationClosesCursorAndDrainsGoroutines(t *testing.T) {
	for _, ignore := range []bool{false, true} {
		name := "honors context"
		if ignore {
			name = "ignores context"
		}
		t.Run(name, func(t *testing.T) {
			before := testutil.GoroutineCount()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			cur := &blockingCursor{entered: make(chan struct{}), release: make(chan struct{}), ignoreContext: ignore}
			done := make(chan error, 1)
			go func() {
				_, err := Run(ctx, &Template{Fields: []string{"n"}}, blockingSource{cur}, limits.Budget{})
				done <- err
			}()
			select {
			case <-cur.entered:
			case <-time.After(time.Second):
				close(cur.release)
				t.Fatal("cursor never reached blocking Next")
			}
			cancel()
			close(cur.release)
			select {
			case err := <-done:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("cancel error=%v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("fill did not stop after cancellation")
			}
			if got := cur.closes.Load(); got != 1 {
				t.Fatalf("cursor Close calls=%d", got)
			}
			testutil.AssertNoGoroutineGrowth(t, before, 1)
		})
	}
}
