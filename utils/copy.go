package utils

import (
	"bytes"
	"encoding/gob"
)

// DeepCopy 使用 gob 将 src 中的信息深度拷贝到 dst 中。
// 速度很慢，只适合操作次数很少的情况
func DeepCopy[T any](dst T, src T) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
