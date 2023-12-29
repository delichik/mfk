package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/delichik/mfk/sqlbuilder/generator/adapter"
)

func generateGet(source *bytes.Buffer, target Target) {
	modelType := reflect.TypeOf(target.Model)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	source.WriteString("\n")
	source.WriteString(fmt.Sprintf("func (m *%s) Get(db *sql.DB, cond string, params ...any) error {\n", modelType.Name()))
	source.WriteString("if cond == \"\" {\n")
	source.WriteString("cond = \"_id = ?\"")
	source.WriteString("params = []any{}")
	source.WriteString("columnString := \"\"\n")
	source.WriteString("columnString := \"\"\n")
	source.WriteString("columnValueString := \"\"\n\n")

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldName := field.Name
		fieldConfig := parseFieldConfig(field)
		if fieldConfig.Skipped {
			continue
		}

		fieldType := field.Type
		for fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		nilCheck := adapter.NilCheck(fieldType)
		if nilCheck != "" {
			source.WriteString(strings.ReplaceAll("if "+nilCheck+" {\n", "${fieldName}", "m."+fieldName))
		}

		asString := adapter.AsString(fieldType)

		source.WriteString(fmt.Sprintf("columnString += \"%s\"\n", fieldConfig.ColumnName))
		source.WriteString("columnString += \",\"\n")
		source.WriteString(strings.ReplaceAll("columnValueString += "+asString+"\n", "${fieldName}", "m."+fieldName))
		source.WriteString("columnValueString += \",\"\n")
		if nilCheck != "" {
			source.WriteString("}\n")
		}
		source.WriteString("\n")
	}
	source.WriteString(fmt.Sprintf("query := \"update %s (\"\n"+
		"query += columnString\n"+
		"query += \") values (\"\n"+
		"query += columnValueString\n"+
		"query += \")\"\n", target.TableName))

	source.WriteString("_, err := db.Exec(query)\n")
	source.WriteString("if err != nil {\n")
	source.WriteString("return err\n")
	source.WriteString("}\n")
	source.WriteString("return nil\n")
	source.WriteString("}\n")
}

func generateGetWhere(source *bytes.Buffer, target Target) {
	// TODO
}
