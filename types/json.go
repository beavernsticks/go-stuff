package bsgostuff_types

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

type JSON struct {
	value any
	set   bool
	null  bool
}

func (j *JSON) Set(value any) error {
	j.set = true
	j.null = false

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &j.value)
	case string:
		return json.Unmarshal([]byte(v), &j.value)
	default:
		j.value = v
	}
	return nil
}

func (j JSON) ToPgx() pgtype.Text {
	if !j.set || j.null {
		return pgtype.Text{Valid: false}
	}
	data, _ := json.Marshal(j.value)
	return pgtype.Text{String: string(data), Valid: true}
}
