package bytepool

import "testing"

var originalData = []byte("1234567890")

func TestArea_Write(t *testing.T) {
	p := New(8, 2)
	a := p.NewArea(12)
	n, err := a.Write(originalData)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if n != len(originalData) {
		t.Error("n != len(originalData)")
		t.FailNow()
	}

}

func TestArea_Read(t *testing.T) {
	p := New(8, 2)
	a := p.NewArea(12)
	_, _ = a.Write(originalData)
	_, _ = a.Write(originalData)
	d := make([]byte, 12)
	n, err := a.Read(d)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if n != 12 {
		t.Error("n != len(originalData)")
		t.FailNow()
	}

	for i, c := range d {
		if c != originalData[i%len(originalData)] {
			t.Errorf("c != originalData[%d]", i%len(originalData))
			t.FailNow()
		}
	}
}

func TestArea_Detach1(t *testing.T) {
	p := New(8, 2)
	a := p.NewArea(12)
	a.Detach()
	if a.data.attached != 0 {
		t.Error("a.data.attached != 0")
		t.FailNow()
	}

	if a.data.pages != nil {
		t.Error("a.data.pages != nil")
		t.FailNow()
	}

	_, err := a.Read([]byte{})
	if err == nil {
		t.Error("a.read err == nil after detached")
		t.FailNow()
	}
}

func TestArea_Detach2(t *testing.T) {
	p := New(8, 2)
	a := p.NewArea(12)
	b := a.ReadOnlyCopy()
	a.Detach()
	if a.data.attached == 0 {
		t.Error("a.data.attached == 0")
		t.FailNow()
	}
	if a.data.pages == nil {
		t.Error("a.data.pages == nil")
		t.FailNow()
	}

	_, err := a.Read([]byte{})
	if err == nil {
		t.Error("a.read err == nil after detached")
		t.FailNow()
	}
	_, err = b.Read([]byte{})
	if err != nil {
		t.Error("b.read err != nil after detached")
		t.FailNow()
	}

	b.Detach()
	if a.data.attached != 0 {
		t.Error("b.data.attached != 0")
		t.FailNow()
	}

	if a.data.pages != nil {
		t.Error("b.data.pages != nil")
		t.FailNow()
	}
	_, err = b.Read([]byte{})
	if err == nil {
		t.Error("b.read err == nil after detached")
		t.FailNow()
	}
}

func TestArea_Detach3(t *testing.T) {
	p := New(8, 2)
	a := p.NewArea(12)
	b := a.WritableCopy()
	a.Detach()
	if a.data.attached != 0 {
		t.Error("a.data.attached != 0")
		t.FailNow()
	}
	if a.data.pages != nil {
		t.Error("a.data.pages != nil")
		t.FailNow()
	}

	_, err := a.Read([]byte{})
	if err == nil {
		t.Error("a.read err == nil after detached")
		t.FailNow()
	}
	_, err = b.Read([]byte{})
	if err != nil {
		t.Error("b.read err != nil after detached")
		t.FailNow()
	}

	b.Detach()
	if b.data.attached != 0 {
		t.Error("b.data.attached != 0")
		t.FailNow()
	}

	if b.data.pages != nil {
		t.Error("b.data.pages != nil")
		t.FailNow()
	}
	_, err = b.Read([]byte{})
	if err == nil {
		t.Error("b.read err == nil after detached")
		t.FailNow()
	}
}
