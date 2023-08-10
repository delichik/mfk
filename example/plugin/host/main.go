package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/delichik/mfk/logger"
	"github.com/delichik/mfk/plugin"
)

func main() {
	logger.InitDefault(&logger.Config{
		Level:   "debug",
		Format:  "json",
		LogPath: "stdout",
	})
	h := plugin.NewHost("example-plugin-host", "0.0.1")
	err := h.Load("./example/plugin")
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		h.Call("hello", []byte("hello example-plugin-plugin"))
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	signal.Stop(signalChan)
}
