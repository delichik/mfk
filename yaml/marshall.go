package yaml

import (
	"errors"
	"reflect"

	"gopkg.in/yaml.v3"
)

// MarshallWithComments 编码 yaml 并读取 comment 标签作为注释写入到编码后的 yaml 字符串中
// 速度很慢，适合少量使用的情况
func MarshallWithComments(obj interface{}) ([]byte, error) {
	n := &yaml.Node{}
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		return nil, errors.New("invalid object")
	}
	err := n.Encode(obj)
	if err != nil {
		return nil, err
	}
	addComments(n, v)
	return yaml.Marshal(n)
}

// AddComments 读取 comment 标签作为注释写入到 yaml.Node 中
// 速度很慢，适合少量使用的情况
func AddComments(node *yaml.Node, obj interface{}) error {
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		return errors.New("invalid object")
	}
	addComments(node, v)
	return nil
}

func addComments(node *yaml.Node, v reflect.Value) {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() || !v.IsValid() {
			return
		} else {
			v = v.Elem()
		}
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.Type().NumField(); i++ {
			f := v.Type().Field(i)
			for j, n := range node.Content {
				if j&1 == 0 {
					if f.Tag.Get("yaml") == n.Value {
						switch f.Type.Kind() {
						case reflect.Map:
							fallthrough
						case reflect.Struct:
							fallthrough
						case reflect.Pointer:
							fallthrough
						case reflect.Array:
							fallthrough
						case reflect.Slice:
							fallthrough
						case reflect.Interface:
							addComments(node.Content[j+1], v.Field(i))
							fallthrough
						default:
							c := f.Tag.Get("comment")
							node.Content[j+1].LineComment = c
							n.LineComment = c
						}
						break
					}
				}
			}
		}
	} else if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			addComments(node.Content[i], f)
		}
	} else if v.Kind() == reflect.Map {
		keys := v.MapKeys()
		for _, k := range keys {
			f := v.MapIndex(k)
			for j, n := range node.Content {
				if n.Value == k.String() {
					addComments(node.Content[j+1], f)
					break
				}
			}
		}
	}
}
