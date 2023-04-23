package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"

	"github.com/delichik/mfk/config"
	"github.com/delichik/mfk/logger"
)

var (
	ctx                 context.Context
	cancel              context.CancelFunc
	cm                  *config.Manager
	beforeRunCall       func()
	afterRunCall        func()
	modules             map[string]Module
	orderedModules      []Module
	autoLoadModuleCount int
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	modules = map[string]Module{}
}

type CommandLineVars struct {
	ConfigPath string
	Help       bool
	Version    bool
}

func parseFlags(version string) CommandLineVars {
	clvs := CommandLineVars{}
	flag.StringVar(&clvs.ConfigPath, "c", "config.yaml", "Config path")
	flag.BoolVar(&clvs.Help, "h", false, "Print help")
	flag.BoolVar(&clvs.Version, "v", false, "Print version")
	flag.Parse()
	if clvs.Help {
		flag.Usage()
		os.Exit(1)
	}
	if clvs.Version {
		fmt.Println("Go version:\t\t", runtime.Version())
		fmt.Println("Binary version:\t", version)
		os.Exit(1)
	}
	return clvs
}

func RegisterAutoLoadModule(module Module) {
	RegisterModule(module)
	autoLoadModuleCount++
}

func RegisterModule(module Module) {
	_, ok := modules[module.Name()]
	if ok {
		panic(errors.New("module " + module.Name() + " already registered"))
	}
	modules[module.Name()] = module
	orderedModules = append(orderedModules, module)
}

func BeforeRun(call func()) {
	beforeRunCall = call
}

func AfterRun(call func()) {
	afterRunCall = call
}

func Run(version string) {
	clvs := parseFlags(version)
	cm = config.NewManager(ctx, clvs.ConfigPath)
	for _, module := range orderedModules {
		module.SetConfigManager(cm)
		if !module.ConfigRequired() {
			continue
		}
		if module.AdditionalLogger() {
			c := logger.GetDefaultConfig()
			c.LogPath = "logs/" + module.Name() + ".log"
			config.RegisterModuleConfig(module.Name()+"-logger", c)
		}
	}

	err := cm.Init()
	if err != nil {
		log.Printf("Init config failed: %s, exit", err.Error())
		return
	}
	cm.SetReloadCallback(ReloadConfig)
	logger.InitDefault(cm)
	for _, module := range orderedModules {
		if module.AdditionalLogger() {
			logger.Init(module.Name()+"-logger", cm)
		}
	}
	logger.Info("App init", zap.String("version", version))

	if beforeRunCall != nil {
		beforeRunCall()
	}
	logger.Info("Loading app modules")
	for i, module := range orderedModules {
		logger.Debug("Prepare module",
			zap.String("name", module.Name()),
			zap.Bool("auto_loaded", i < autoLoadModuleCount))
		if module.ConfigRequired() {
			cfg := cm.GetModuleConfig(module.Name())
			if cfg == nil {
				logger.Warn("Skip module", zap.String("name", module.Name()))
				continue
			}
			logger.Debug("Applying module config", zap.String("name", module.Name()))
			err := module.ApplyConfig(cfg)
			if err != nil {
				logger.Fatal("Apply module config failed, exit",
					zap.String("name", module.Name()),
					zap.Error(err))
			}
		}
		err = module.Run(ctx)
		if err != nil {
			logger.Fatal("Run module failed",
				zap.String("name", module.Name()),
				zap.Error(err))
		}
	}
	logger.Info("App modules loaded")
	if afterRunCall != nil {
		afterRunCall()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	signal.Stop(signalChan)
	logger.Info("App shutdown")

	cancel()
	for _, module := range orderedModules {
		if module.ConfigRequired() {
			cfg := cm.GetModuleConfig(module.Name())
			if cfg == nil {
				continue
			}
		}
		module.Exit()
	}
}

func ReloadConfig(name string, cfg config.ModuleConfig) {
	module, ok := modules[name]
	if !ok {
		return
	}
	logger.Info("Reloading module config", zap.String("name", name))
	err := module.ApplyConfig(cfg)
	if err != nil {
		logger.Error("Apply module config failed",
			zap.String("name", name),
			zap.Error(err))
	}
}
