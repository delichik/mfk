package utils

import "testing"

func TestDeepCopy(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	type TestStruct struct {
		A User
		B *User
	}

	s1 := &TestStruct{
		A: User{
			"Able",
			11,
		},
		B: &User{
			Name: "Unable",
			Age:  -11,
		},
	}

	s2 := &TestStruct{}
	err := DeepCopy(s2, s1)
	if err != nil {
		t.FailNow()
	}
	if s1.A.Name != s2.A.Name ||
		s1.A.Age != s2.A.Age {
		t.FailNow()
	}
	if s2.B == nil ||
		s1.B.Name != s2.B.Name ||
		s1.B.Age != s2.B.Age {
		t.FailNow()
	}
}
