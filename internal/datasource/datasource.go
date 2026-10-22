package datasource

import "context"

// DataSource opens a forward-only stream of records for the fill stage.
//
// A DataSource is a factory: each call to Open returns an independent Cursor.
// Implementations must be safe for concurrent Open calls unless a concrete
// source documents otherwise.
//
// Open must never return (nil, nil). A nil Cursor with a nil error is a
// programming error in the implementation, not a valid empty result; callers
// may treat it as a bug. Open must honor ctx and return ctx.Err() when the
// context is already canceled.
type DataSource interface {
	Open(ctx context.Context) (Cursor, error)
}

// Cursor iterates the records of an opened DataSource.
//
// Lifecycle:
//
//   - Next advances one record and reports whether a record is available. It
//     returns (nil, false, nil) when the stream is exhausted. After Close it
//     returns (nil, false, nil).
//   - Close releases the cursor's resources. It is idempotent: calling Close
//     more than once returns nil after the first successful call.
//
// A Cursor is not safe for concurrent use. The caller owns the cursor and must
// Close it, typically with defer, even when Next fails or the caller returns
// early.
type Cursor interface {
	Next(ctx context.Context) (Record, bool, error)
	Close() error
}

// Record exposes the fields of a single iteration by name.
//
// Field names are case-sensitive and match the declared report field names. A
// missing field returns (nil, false).
//
// A Record is only guaranteed valid until the next Next call on its Cursor;
// callers must copy any values they intend to retain beyond that point.
type Record interface {
	Value(field string) (any, bool)
}
