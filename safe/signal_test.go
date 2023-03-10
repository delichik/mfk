package safe

import (
	"testing"
)

func TestRWMutex_Lock(t *testing.T) {
	runner := NewSignal()
	runner.Pause(SignalDefaultOwner)
	if !runner.paused {
		t.FailNow()
	}
	runner.Pause(SignalDefaultOwner)
	if !runner.paused {
		t.FailNow()
	}
	runner.Resume(SignalDefaultOwner)
	if runner.paused {
		t.FailNow()
	}
	runner.Pause("1")
	if runner.paused {
		t.FailNow()
	}
	runner.Pause(SignalDefaultOwner)
	if !runner.paused {
		t.FailNow()
	}
	runner.Resume("1")
	if runner.paused {
		t.FailNow()
	}
	runner.Pause("1")
	if !runner.paused {
		t.FailNow()
	}
}
