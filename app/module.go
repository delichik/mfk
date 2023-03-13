package app

import (
	"github.com/delichik/mfk/config"
)

type DynamicModule struct{}

func (m *DynamicModule) DynamicConfig() bool {
	return true
}

type StaticModule struct{}

func (m *StaticModule) DynamicConfig() bool {
	return false
}

type CriticalModule struct{}

func (m *CriticalModule) Critical() bool {
	return true
}

type UncriticalModule struct{}

func (m *UncriticalModule) Critical() bool {
	return false
}

type Module interface {
	// Name 返回模块的名字，所有模块名字不能重复
	Name() string

	// ApplyConfig 触发配置应用，当启动和配置发生变化时会被调用
	ApplyConfig(cfg config.ModuleConfig) error

	// WaitExit 等待 Module 完成退出操作
	WaitExit()

	// DynamicConfig 返回当前模块是否允许多次调用 ApplyConfig
	DynamicConfig() bool

	// Critical 返回当前 Module 是否为关键 Module
	Critical() bool
}
