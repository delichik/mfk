package plugin

import (
	"encoding/base64"
	"os"
	"os/exec"

	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"

	"github.com/delichik/mfk/logger"
)

type Host struct {
	plugins map[string]*Entity
	name    string
	version string
}

func NewHost(name string, version string) *Host {
	return &Host{
		plugins: make(map[string]*Entity),
		name:    name,
		version: version,
	}
}

func (h *Host) Load(pluginPath string) error {
	handshake, err := msgpack.Marshal(&HandshakeInfo{
		Name:    h.name,
		Version: h.version,
	})
	if err != nil {
		panic(err)
	}

	handshakeStr := base64.StdEncoding.EncodeToString(handshake)

	entries, err := os.ReadDir(pluginPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		logger.Info("load plugin", zap.String("name", entry.Name()))
		cmd := exec.Command(pluginPath+"/"+entry.Name(), "-h", handshakeStr)
		h.plugins[entry.Name()] = newEntity(cmd)
		err = h.plugins[entry.Name()].Start()
		if err != nil {
			logger.Error("fail to load plugin", zap.String("name", entry.Name()), zap.Error(err))
		}
	}

	return nil
}

func (h *Host) Call(call string, data []byte) ([]byte, error) {
	for _, plg := range h.plugins {
		plg.Call(call, data)
		return []byte(""), nil
	}
	return []byte(""), nil
}
