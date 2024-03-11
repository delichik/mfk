package app

import (
	"context"
	"reflect"
	"runtime"
	"strconv"

	"github.com/delichik/mfk/config"
)

type Module[T config.ModuleConfig] interface {
	// Name 返回模块的名字，所有模块名字不能重复
	Name() string

	// ConfigRequired 返回是否需要配置
	ConfigRequired() bool

	// ApplyConfig 触发配置应用，当启动和配置发生变化时会被调用
	ApplyConfig(cfg T) error

	// OnRun 启动模块，当应用启动时会被调用
	OnRun(ctx context.Context) error

	// AdditionalLogger 返回是否需要额外的日志记录
	AdditionalLogger() bool

	// OnExit 触发模块退出，当应用退出时会被调用
	OnExit()

	SetConfigManager(cm *config.Manager)
}

type ModuleEntry struct {
	module     any
	registerer string

	Name             func() string
	ConfigRequired   func() bool
	ApplyConfig      func(cfg config.ModuleConfig) error
	OnRun            func(ctx context.Context) error
	AdditionalLogger func() bool
	OnExit           func()
	SetConfigManager func(cm *config.Manager)
}

func nullableError(in any) error {
	switch r := in.(type) {
	case error:
		return r
	default:
		return nil
	}
}

func newModuleEntry(module any) *ModuleEntry {
	_, file, line, _ := runtime.Caller(3)
	moduleEntry := &ModuleEntry{
		module:     module,
		registerer: file + ":" + strconv.Itoa(line),
	}

	rv := reflect.ValueOf(module)
	moduleEntry.Name = func() string {
		fn := rv.MethodByName("Name")
		res := fn.Call([]reflect.Value{})
		return (res[0].Interface()).(string)
	}

	moduleEntry.ConfigRequired = func() bool {
		fn := rv.MethodByName("ConfigRequired")
		res := fn.Call([]reflect.Value{})
		return (res[0].Interface()).(bool)
	}

	moduleEntry.ApplyConfig = func(cfg config.ModuleConfig) error {
		fn := rv.MethodByName("ApplyConfig")
		res := fn.Call([]reflect.Value{reflect.ValueOf(cfg)})
		return nullableError(res[0].Interface())
	}

	moduleEntry.OnRun = func(ctx context.Context) error {
		fn := rv.MethodByName("OnRun")
		res := fn.Call([]reflect.Value{reflect.ValueOf(ctx)})
		return nullableError(res[0].Interface())
	}

	moduleEntry.AdditionalLogger = func() bool {
		fn := rv.MethodByName("AdditionalLogger")
		res := fn.Call([]reflect.Value{})
		return (res[0].Interface()).(bool)
	}

	moduleEntry.OnExit = func() {
		fn := rv.MethodByName("OnExit")
		fn.Call([]reflect.Value{})
	}

	moduleEntry.SetConfigManager = func(cm *config.Manager) {
		fn := rv.MethodByName("SetConfigManager")
		fn.Call([]reflect.Value{reflect.ValueOf(cm)})
	}

	return moduleEntry
}
