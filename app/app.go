package app

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/delichik/my-go-pkg/config"
	"github.com/delichik/my-go-pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	modules        map[string]Module
	orderedModules []Module

	l      *zap.Logger
	cm     *config.Manager
	ctx    context.Context
	cancel context.CancelFunc
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

func (a *App) Run() {
	a.l.Info("App init")
	err := a.cm.Init()
	if err != nil {
		log.Println("Init config failed, exit")
		return
	}
	a.cm.SetReloadCallback(a.ReloadConfig)
	a.l = logger.NewLogger("app")

	wg := sync.WaitGroup{}

	a.l.Info("Loading app modules")
	for _, module := range a.orderedModules {
		cfg := a.cm.GetModuleConfig(module.Name())
		if cfg == nil {
			a.l.Debug("Skip module", zap.String("name", module.Name()))
			continue
		}
		a.l.Debug("Applying module config", zap.String("name", module.Name()))
		err := module.ApplyConfig(cfg)
		if err != nil {
			if module.Critical() {
				a.l.Fatal("Apply critical module config failed, exit",
					zap.String("name", module.Name()),
					zap.Error(err))
			} else {
				a.l.Error("Apply module config failed",
					zap.String("name", module.Name()),
					zap.Error(err))
			}
		}
		wg.Add(1)
		go func(module Module) {
			defer wg.Done()
			module.WaitExit()
		}(module)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	signal.Stop(signalChan)
	a.l.Info("App shutdown")
	wg.Wait()
}

func (a *App) ReloadConfig(name string, cfg config.ModuleConfig) {
	module, ok := a.modules[name]
	if !ok {
		return
	}
	a.l.Info("Reloading module config", zap.String("name", name))
	err := module.ApplyConfig(cfg)
	if err != nil {
		if module.Critical() {
			a.l.Fatal("Apply critical module config failed, exit",
				zap.String("name", name),
				zap.Error(err))
		} else {
			a.l.Error("Apply module config failed",
				zap.String("name", name),
				zap.Error(err))
		}
	}
}
