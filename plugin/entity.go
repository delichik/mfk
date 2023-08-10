package plugin

import (
	"bufio"
	"encoding/json"
	"io"
	"os/exec"

	"go.uber.org/zap"

	"github.com/delichik/mfk/logger"
)

type Entity struct {
	cmd *exec.Cmd

	pluginOutput   io.ReadCloser
	pluginInput    io.WriteCloser
	stdoutBuffered *bufio.Reader
}

func newEntity(cmd *exec.Cmd) *Entity {
	e := &Entity{
		cmd: cmd,
	}
	return e
}

func (e *Entity) Start() error {
	var err error
	e.pluginOutput, err = e.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	e.pluginInput, err = e.cmd.StdinPipe()
	if err != nil {
		return err
	}
	e.stdoutBuffered = bufio.NewReader(e.pluginOutput)

	err = e.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		for {
			call, data, err := read(e.stdoutBuffered)
			if err != nil {
				return
			}
			switch call {
			case _CALL_LOGGER:
				log(e.cmd.Path, data)
			// case _CALL_REGISTER_METHOD:
			// 	registerMethod(e.cmd.Path, data)
			default:
				logger.Info("plugin order", zap.String("call", call), zap.ByteString("data", data))
			}
		}
	}()

	return nil
}

func (e *Entity) Stop() error {
	return e.cmd.Process.Kill()
}

func (e *Entity) Call(call string, data []byte) error {
	err := send(e.pluginInput, call, data)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) CallWithResponse(call string, data []byte) ([]byte, error) {
	err := send(e.pluginInput, call, data)
	if err != nil {
		return nil, err
	}
	return []byte(""), nil
}

type logContent struct {
	Level   string `json:"level"`
	Caller  string `json:"caller"`
	Message string `json:"msg"`
}

func log(pluginName string, data []byte) {
	c := logContent{}
	err := json.Unmarshal(data, &c)
	if err != nil {
		return
	}

	fieldMap := map[string]interface{}{}
	_ = json.Unmarshal(data, &fieldMap)
	fields := []zap.Field{}
	fields = append(fields, zap.String("plugin_name", pluginName), zap.String("plugin_caller", c.Caller))
	for k, v := range fieldMap {
		if k == "level" ||
			k == "ts" ||
			k == "caller" ||
			k == "msg" ||
			k == "plugin_name" ||
			k == "plugin_caller" {
			continue
		}
		fields = append(fields, zap.Any(k, v))
	}
	c.Message = "[plugin] " + c.Message
	switch c.Level {
	case "debug":
		logger.Debug(c.Message, fields...)
	case "info":
		logger.Info(c.Message, fields...)
	case "warn":
		logger.Warn(c.Message, fields...)
	case "error":
		logger.Error(c.Message, fields...)
	}
}
