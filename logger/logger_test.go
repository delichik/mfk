package logger

import (
	"os"
	"testing"

	"github.com/delichik/my-go-pkg/config"
)

func TestConfig_NewLogger(t *testing.T) {
	_ = os.Remove("./config.yaml")
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		t.FailNow()
	}
	Init(cfg)
	l := NewLogger("log_test")
	if l == nil {
		t.FailNow()
	}
	l.Info("logger testing")
}
