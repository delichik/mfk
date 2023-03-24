package app

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/delichik/mfk/config"
	"github.com/delichik/mfk/logger"
)

var (
	ctx                 context.Context
	cancel              context.CancelFunc
	cm                  *config.Manager
	afterRunCall        func()
	modules             map[string]Module
	orderedModules      []Module
	autoLoadModuleCount int
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	modules = map[string]Module{}
	clvs := parseFlags()
	cm = config.NewManager(ctx, clvs.ConfigPath)
}

type CommandLineVars struct {
	ConfigPath string
	Help       bool
}

func parseFlags() CommandLineVars {
	clvs := CommandLineVars{}
	flag.StringVar(&clvs.ConfigPath, "c", "config.yaml", "Config path")
	flag.BoolVar(&clvs.Help, "h", false, "Print help")
	flag.Parse()
	if clvs.Help {
		flag.Usage()
		os.Exit(1)
	}
	return clvs
}

func RegisterAutoLoadModule(module Module) {
	modules[module.Name()] = module
	orderedModules = append(orderedModules, module)
	autoLoadModuleCount++
}

func RegisterModule(module Module) {
	modules[module.Name()] = module
	orderedModules = append(orderedModules, module)
}

func AfterRun(call func()) {
	afterRunCall = call
}

func Run() {
	for _, module := range orderedModules {
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
	logger.Info("App init")

	logger.Info("Loading app modules")
	for i, module := range orderedModules {
		cfg := cm.GetModuleConfig(module.Name())
		if cfg == nil {
			logger.Debug("Skip module", zap.String("name", module.Name()))
			continue
		}
		if i < autoLoadModuleCount {
			logger.Debug("Applying auto load module config", zap.String("name", module.Name()))
		} else {
			logger.Debug("Applying module config", zap.String("name", module.Name()))
		}
		err := module.ApplyConfig(cfg)
		if err != nil {
			logger.Fatal("Apply module config failed, exit",
				zap.String("name", module.Name()),
				zap.Error(err))
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

	for _, module := range orderedModules {
		cfg := cm.GetModuleConfig(module.Name())
		if cfg == nil {
			continue
		}
		module.Exit()
	}
	cancel()
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
