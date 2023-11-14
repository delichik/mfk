package generator

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"golang.org/x/tools/imports"
)

func Add(pkgPath string, filename string, modelTargets ...Target) {
	source := bytes.NewBuffer([]byte{})
	for i, target := range modelTargets {
		modelType := reflect.TypeOf(target.Model)
		for modelType.Kind() == reflect.Pointer {
			modelType = modelType.Elem()
		}
		if i == 0 {
			pkgParts := strings.Split(modelType.PkgPath(), "/")
			source.WriteString(fmt.Sprintf("package %s\n", pkgParts[len(pkgParts)-1]))
		}
		generateCreate(source, target)
		generateUpdate(source, target)
		generateScan(source, target)
		generateGet(source, target)
		generateDelete(source, target)
		generateUpdateWhere(source, target)
		generateGetWhere(source, target)
		generateDeleteWhere(source, target)
	}
	data, err := imports.Process("", source.Bytes(), &imports.Options{
		TabIndent: true,
	})
	if err != nil {
		fmt.Println("do imports failed: ", err.Error())
		os.WriteFile(pkgPath+"/"+filename, source.Bytes(), 0777)
		return
	}
	os.WriteFile(pkgPath+"/"+filename, data, 0777)
}
