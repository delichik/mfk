package app

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/delichik/mfk/config"
	"github.com/delichik/mfk/logger"
	"go.uber.org/zap"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	cm             *config.Manager
	afterRunCall   func()
	modules        map[string]Module
	orderedModules []Module
}

type CommandLineVars struct {
	ConfigPath string
	Help       bool
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	app := &App{
		ctx:     ctx,
		cancel:  cancel,
		modules: map[string]Module{},
	}
	clvs := parseFlags()
	app.cm = config.NewManager(ctx, clvs.ConfigPath)
	return app
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

func (a *App) RegisterModule(module Module) {
	a.modules[module.Name()] = module
	a.orderedModules = append(a.orderedModules, module)
}

func (a *App) AfterRun(call func()) {
	a.afterRunCall = call
}

func (a *App) Run() {
	err := a.cm.Init()
	if err != nil {
		log.Println("Init config failed, exit")
		return
	}
	a.cm.SetReloadCallback(a.ReloadConfig)
	logger.Init(a.cm)
	logger.Info("App init")

	logger.Info("Loading app modules")
	for _, module := range a.orderedModules {
		cfg := a.cm.GetModuleConfig(module.Name())
		if cfg == nil {
			logger.Debug("Skip module", zap.String("name", module.Name()))
			continue
		}
		logger.Debug("Applying module config", zap.String("name", module.Name()))
		err := module.ApplyConfig(cfg)
		if err != nil {
			logger.Fatal("Apply module config failed, exit",
				zap.String("name", module.Name()),
				zap.Error(err))
		}
		err = module.Run(a.ctx)
		if err != nil {
			logger.Fatal("Run module failed",
				zap.String("name", module.Name()),
				zap.Error(err))
		}
	}
	logger.Info("App modules loaded")
	if a.afterRunCall != nil {
		a.afterRunCall()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	signal.Stop(signalChan)
	logger.Info("App shutdown")

	for _, module := range a.orderedModules {
		cfg := a.cm.GetModuleConfig(module.Name())
		if cfg == nil {
			continue
		}
		module.Exit()
	}
	a.cancel()
}

func (a *App) ReloadConfig(name string, cfg config.ModuleConfig) {
	module, ok := a.modules[name]
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
