package main

import (
	"github.com/delichik/daf/app"
	"github.com/delichik/daf/logger"
)

func main() {
	app.BeforeRun(func() {
		logger.Info("before run")
	})
	app.AfterRun(func() {
		logger.Info("after run")
	})
	app.RegisterAutoLoadModule(&DemoNoConfModule{}, &app.NoConfig{})
	app.RegisterModule(&DemoModule{}, &DemoModuleConfig{})
	app.Run("0.0.1")
}
