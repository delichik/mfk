package plugin

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/delichik/mfk/logger"
)

type Options struct {
	Name               string
	Version            string
	HostName           string
	HostMinimalVersion string
}

var registeredPlugins = make(map[string]Plugin)

func RegisterHandler(name string, plugin Plugin) {
	registeredPlugins[name] = plugin
}

func RunPlugin(options *Options) {
	if options.Name == "" {
		panic("plugin name is required")
	}

	if options.Version == "" {
		panic("plugin version is required")
	}

	if options.HostName == "" {
		panic("plugin host name is required")
	}

	if options.HostMinimalVersion == "" {
		panic("plugin host minimal version is required")
	}

	handshake := ""
	flag.StringVar(&handshake, "h", "", "")
	flag.Parse()
	if !checkHandshake(handshake, options) {
		fmt.Printf("This executable binary is a plugin for %s %s+, Do not run it alone\n",
			options.HostName, options.HostMinimalVersion)
		os.Exit(1)
	}

	for name, plugin := range registeredPlugins {
		fmt.Printf("Plugin %s is starting...\n", name)
		plugin.Init()
		fmt.Printf("Plugin %s is started\n", name)
	}

	for name := range registeredPlugins {
		send(os.Stdout, _CALL_REGISTER_METHOD, []byte(name))
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		buf := bufio.NewReader(os.Stdin)
		for {
			call, data, err := read(buf)
			if err != nil {
				return
			}
			if ctx.Err() != nil {
				return
			}
			plugin, ok := registeredPlugins[call]
			if !ok {
				send(os.Stdout, call+"_rsp", []byte(""))
				continue
			}
			rsp, err := plugin.Handle(data)
			if err != nil {
				logger.Error("plugin handle failed", zap.String("call", call), zap.Error(err))
				send(os.Stdout, call+"_rsp", []byte(err.Error()))
				continue
			}
			send(os.Stdout, call+"_rsp", rsp)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	cancel()
	signal.Stop(signalChan)

	for _, plugin := range registeredPlugins {
		plugin.UnInit()
	}
}
