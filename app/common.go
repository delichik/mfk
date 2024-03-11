package app

import (
	"context"

	"github.com/delichik/daf/config"
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

func (m *NoConfigModule) DefaultConfig() config.ModuleConfig {
	return nil
}

func (m *NoConfigModule) ApplyConfig(_ config.ModuleConfig) error {
	return nil
}

func (m *NoConfigModule) SetConfigManager(_ *config.Manager) {}

type ConfigRequiredModule struct{}

func (m *ConfigRequiredModule) ConfigRequired() bool {
	return true
}

func (m *ConfigRequiredModule) SetConfigManager(_ *config.Manager) {}

type InitializerModule struct{}

func (m *InitializerModule) OnRun(_ context.Context) error {
	return nil
}

func (m *InitializerModule) OnExit() {}
