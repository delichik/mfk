package main

import (
	"context"
	"strconv"
	"time"

	"github.com/delichik/daf/app"
	"github.com/delichik/daf/config"
	"github.com/delichik/daf/logger"
	"github.com/delichik/daf/utils"
)

const ModuleName = "demo_module"

type DemoModuleConfig struct {
	ListenAddr string `yaml:"listen-addr" comment:"Address to listen, 0.0.0.0 for default"`
	ListenPort int    `yaml:"listen-port" comment:"Port to listen, 80 for default"`
}

func (c *DemoModuleConfig) Check() error {
	if c.ListenAddr == "" {
		c.ListenAddr = "0.0.0.0"
	}

	if c.ListenPort <= 0 || c.ListenPort > 65534 {
		c.ListenPort = 80
	}
	return nil
}

func (c *DemoModuleConfig) Clone() config.ModuleConfig {
	n := &DemoModuleConfig{}
	_ = utils.DeepCopy(n, c)
	return n
}

func (c *DemoModuleConfig) Compare(moduleConfig config.ModuleConfig) bool {
	return utils.DeepCompare(c, moduleConfig)
}

type DemoModule struct {
	app.ConfigRequiredModule
	app.DefaultLoggerModule
}

func (m *DemoModule) Name() string {
	return ModuleName
}

func (m *DemoModule) ApplyConfig(cfg *DemoModuleConfig) error {
	logger.Info("apply config")
	return nil
}

func (m *DemoModule) OnRun(ctx context.Context) error {
	logger.Info("on run")
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for i := 5; i > 0; i-- {
			logger.Info("shutdown after " + strconv.Itoa(i))
			<-ticker.C
		}
		app.Shutdown()
	}()
	return nil
}

func (m *DemoModule) OnExit() {
	logger.Info("on exit")
}

func main() {
	app.BeforeRun(func() {
		logger.Info("before run")
	})
	app.AfterRun(func() {
		logger.Info("after run")
	})
	app.RegisterModule(&DemoModule{}, &DemoModuleConfig{})
	app.Run("0.0.1")
}
