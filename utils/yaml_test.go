package utils

import (
	"testing"
)

func TestYaml_MarshallWithComments(t *testing.T) {
	type C struct {
		D []string `yaml:"d" comment:"ddd"`
		H string   `yaml:"h" comment:"hhh"`
	}

	type T struct {
		A string      `yaml:"a" comment:"aaa"`
		B C           `yaml:"b" comment:"bbb"`
		E interface{} `yaml:"e" comment:"eee"`
		F interface{} `yaml:"f" comment:"fff"`
		G []C         `yaml:"g" comment:"ggg"`
		I C           `yaml:"i"`
	}

	u := &T{
		E: &C{},
		G: []C{
			{
				D: []string{"123", "321"},
				H: "111222333",
			},
			{},
		},
	}

	r, err := YamlMarshallWithComments(u)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(string(r))
}
