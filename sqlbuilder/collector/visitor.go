package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type _TypeConfig struct {
	Wired     bool
	TableName string
}

type _Model struct {
	PackageName string
	StructName  string
	Config      *_TypeConfig
}

type _TypeSpecVisitor struct {
	packageName    string
	lastTypeConfig *_TypeConfig
	models         []_Model
}

func (v *_TypeSpecVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		v.packageName = n.Name.Name
		fmt.Println("found package:", v.packageName)
	case *ast.Comment:
		fmt.Println(n.Text)
		text := strings.ToLower(strings.TrimSpace(n.Text))
		if strings.HasPrefix(text, "@dym/") {
			kv := strings.Split(text[5:], ":")
			switch kv[0] {
			case "wired":
				v.lastTypeConfig.Wired = kv[1] == "true"
			case "table":
				v.lastTypeConfig.TableName = kv[1]
			}
		}
	case *ast.TypeSpec:
		fmt.Println("found struct:", n.Name.Name, ", table name:", v.lastTypeConfig.TableName)
		_, ok := n.Type.(*ast.StructType)
		if ok {
			if v.lastTypeConfig.TableName == "" {
				v.lastTypeConfig.TableName = n.Name.Name
			}
			v.models = append(v.models, _Model{
				PackageName: v.packageName,
				StructName:  n.Name.Name,
				Config:      v.lastTypeConfig,
			})
		}
	}
	v.lastTypeConfig = &_TypeConfig{}
	return v
}
