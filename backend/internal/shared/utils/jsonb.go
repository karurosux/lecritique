package utils

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONBField[T any] struct {
	Data T
}

func (j JSONBField[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j *JSONBField[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("{}"), &j.Data)
	}
	return json.Unmarshal(bytes, &j.Data)
}

func MarshalJSONB(v interface{}) (driver.Value, error) {
	return json.Marshal(v)
}

func UnmarshalJSONB(value interface{}, dest interface{}, defaultValue string) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(defaultValue), dest)
	}
	return json.Unmarshal(bytes, dest)
}