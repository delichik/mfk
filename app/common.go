package app

import (
	"context"
	"reflect"

	"github.com/delichik/daf/config"
)

var noConfigIfaceType = reflect.TypeOf((*noConfigIface)(nil)).Elem()

type noConfigIface interface {
	__NoConfig()
}

type NoConfig struct{}

func (c *NoConfig) Check() error {
	return nil
}

func (c *NoConfig) Clone() config.ModuleConfig {
	return c
}

func (c *NoConfig) Compare(_ config.ModuleConfig) bool {
	return true
}

func (c *NoConfig) __NoConfig() {}

type AdditionalLoggerModule struct{}

func (m *AdditionalLoggerModule) SetConfigManager(_ *config.Manager) {}

func (m *AdditionalLoggerModule) AdditionalLogger() bool {
	return true
}

type DefaultLoggerModule struct{}

func (m *DefaultLoggerModule) SetConfigManager(_ *config.Manager) {}

func (m *DefaultLoggerModule) AdditionalLogger() bool {
	return false
}

type InitializerModule struct{}

func (m *InitializerModule) OnRun(_ context.Context) error {
	return nil
}

func (m *InitializerModule) OnExit() {}
