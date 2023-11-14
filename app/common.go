package app

import (
	"context"

	"github.com/delichik/mfk/config"
)

// AdditionalLoggerModule The module has a costume logger
type AdditionalLoggerModule struct{}

func (m *AdditionalLoggerModule) AdditionalLogger() bool {
	return true
}

// DefaultLoggerModule The module uses the default logger
type DefaultLoggerModule struct{}

func (m *DefaultLoggerModule) AdditionalLogger() bool {
	return false
}

// NoConfigModule The module has no config and uses the default logger
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

func (m *NoConfigModule) SetConfigManager(_ *config.Manager) {}

// ConfigRequiredModule The module has config
type ConfigRequiredModule struct{}

func (m *ConfigRequiredModule) ConfigRequired() bool {
	return true
}

func (m *ConfigRequiredModule) SetConfigManager(_ *config.Manager) {}

// InitializerModule The module only runs during configure changed
type InitializerModule struct{}

func (m *InitializerModule) Run(_ context.Context) error {
	return nil
}

func (m *InitializerModule) Exit() {}
