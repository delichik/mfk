package generator

import (
	"reflect"
	"strings"
)

type Target struct {
	Model     any
	TableName string
}

type FieldConfig struct {
	Skipped    bool
	ColumnName string
	PrimaryKey bool
}

func parseFieldConfig(field reflect.StructField) FieldConfig {
	tagRaw := field.Tag.Get("gen")
	if tagRaw == "-" {
		return FieldConfig{
			Skipped: true,
		}
	}

	f := FieldConfig{
		ColumnName: field.Name,
	}
	parts := strings.Split(tagRaw, ";")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		switch kv[0] {
		case "primaryKey":
			f.PrimaryKey = true
		case "column":
			f.ColumnName = kv[1]
		}
	}
	return f
}
