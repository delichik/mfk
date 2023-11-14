package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/delichik/mfk/sqlbuilder/generator/adapter"
)

func generateUpdate(source *bytes.Buffer, target Target) {
	modelType := reflect.TypeOf(target.Model)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	source.WriteString("\n")
	source.WriteString(fmt.Sprintf("func (m *%s) Update(db *sql.DB) error {\n", modelType.Name()))
	source.WriteString("columnString := \"\"\n")
	source.WriteString("columnValueString := \"\"\n\n")

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}

		fieldName := field.Name
		fieldJsonName := field.Name

		if tag != "" {
			parts := strings.Split(tag, ",")
			fieldJsonName = parts[0]
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

		source.WriteString(fmt.Sprintf("columnString += \"%s\"\n", fieldJsonName))
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

func generateUpdateWhere(source *bytes.Buffer, target Target) {
	modelType := reflect.TypeOf(target.Model)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	source.WriteString("\n")
	source.WriteString(fmt.Sprintf("func (m *%s) UpdateWhere(db *sql.DB, cond string, params ...interface{}) error {\n", modelType.Name()))
	source.WriteString("columnString := \"\"\n")
	source.WriteString("columnValueString := \"\"\n\n")

mainLoop:
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName := field.Name
		fieldJsonName := field.Name

		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			fieldJsonName = parts[0]
		}

		dymTag := field.Tag.Get("dym")
		if dymTag != "-" {
			parts := strings.Split(dymTag, ";")
			for _, part := range parts {
				switch part {
				case "primaryKey":
					continue mainLoop
				}
			}
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

		source.WriteString(fmt.Sprintf("columnString += \"%s\"\n", fieldJsonName))
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
	source.WriteString("if cond != \"\" {\n")
	source.WriteString("query += \" where \" + cond\n")
	source.WriteString("}\n")
	source.WriteString("_, err := db.Exec(query, params...)\n")
	source.WriteString("if err != nil {\n")
	source.WriteString("return err\n")
	source.WriteString("}\n")
	source.WriteString("return nil\n")
	source.WriteString("}\n")
}
