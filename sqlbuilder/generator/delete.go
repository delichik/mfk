package generator

import (
	"bytes"
	"fmt"
	"reflect"
)

func generateDelete(source *bytes.Buffer, target Target) {
	modelType := reflect.TypeOf(target.Model)
	for modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	source.WriteString("\n")
	source.WriteString(fmt.Sprintf("func (m *%s) Delete(db *sql.DB, cond string, params ...any) error {\n", modelType.Name()))
	source.WriteString(fmt.Sprintf("query := \"delete from %s\"\n", target.TableName))
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

func generateDeleteWhere(source *bytes.Buffer, target Target) {
	// TODO
}
