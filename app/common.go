package app

import (
	"strconv"
	"sync/atomic"

	"github.com/delichik/mfk/config"
)

type AdditionalLoggerModule struct{}

func (m *AdditionalLoggerModule) AdditionalLogger() bool {
	return true
}

type DefaultLoggerModule struct{}

func (m *DefaultLoggerModule) AdditionalLogger() bool {
	return false
}

type NoConfigModule struct{}

func (m *NoConfigModule) ConfigRequired() bool {
	return false
}

func (m *NoConfigModule) AdditionalLogger() bool {
	return false
}

func (m *NoConfigModule) ApplyConfig(_ config.ModuleConfig) error {
	return nil
}

type ConfigRequiredModule struct{}

func (m *ConfigRequiredModule) ConfigRequired() bool {
	return true
}

var initializer uint32 = 0

type InitializerModule struct {
	id string
}

func (m *InitializerModule) Name() string {
	if m.id == "" {
		m.id = strconv.Itoa(int(atomic.AddUint32(&initializer, 1)))
	}
	return "initializer_" + m.id
}

func (m *InitializerModule) Exit() {}
