package app

import (
	"context"

	"github.com/delichik/mfk/config"
)

type Module interface {
	// Name 返回模块的名字，所有模块名字不能重复
	Name() string

	// Run 启动模块，当应用启动时会被调用
	Run(ctx context.Context) error

	// ApplyConfig 触发配置应用，当启动和配置发生变化时会被调用
	ApplyConfig(cfg config.ModuleConfig) error

	// AdditionalLogger 返回是否需要额外的日志记录
	AdditionalLogger() bool

	// Exit 触发模块退出，当应用退出时会被调用
	Exit()
}
