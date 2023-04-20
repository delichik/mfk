package app

import (
	"context"

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

type InitializerModule struct{}

func (m *InitializerModule) Run(_ context.Context) error {
	return nil
}

func (m *InitializerModule) Exit() {}
