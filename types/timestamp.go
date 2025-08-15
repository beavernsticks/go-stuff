package bsgostuff_types

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Timestamp struct {
	value time.Time
	set   bool
	null  bool
}

func Now() Timestamp {
	return Timestamp{value: time.Now().UTC(), set: true}
}

func (ts *Timestamp) Set(value any) error {
	ts.set = true
	ts.null = false

	switch v := value.(type) {
	case time.Time:
		ts.value = v.UTC()
	case pgtype.Timestamptz:
		ts.null = !v.Valid
		ts.value = v.Time.UTC()
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (ts Timestamp) ToPgx() pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: ts.value, Valid: ts.set && !ts.null}
}
