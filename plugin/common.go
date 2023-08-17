package plugin

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap/zapcore"
)

const (
	_ORDER_PREFIX = "__mfk_plugin_order__"
	_SPLITER      = "|/|"
	_CALL_LOGGER  = "logger"
	_CALL_REPLY   = "_reply"
)

var _SAFE_CALL_LOGGER = base64.StdEncoding.EncodeToString([]byte(_CALL_LOGGER))

var locker sync.Mutex

func LogDriver() zapcore.WriteSyncer {
	return &logWriter{writer: os.Stdout}
}

type logWriter struct {
	writer io.Writer
}

func (w *logWriter) Write(data []byte) (n int, err error) {
	safeData := base64.StdEncoding.EncodeToString(data)
	locker.Lock()
	defer locker.Unlock()
	return fmt.Fprintln(w.writer, _ORDER_PREFIX+_SAFE_CALL_LOGGER+_SPLITER+safeData)
}

func (w *logWriter) Sync() error {
	return nil
}

type sendObject struct {
	err     error
	id      uint64
	call    string
	content []byte
}

type sendObject2 struct {
	id      uint64
	call    string
	content []byte
}

func send(writer io.Writer, r *sendObject2) error {
	safeCall := base64.StdEncoding.EncodeToString([]byte(r.call))
	safeData := base64.StdEncoding.EncodeToString(r.content)
	locker.Lock()
	defer locker.Unlock()
	_, err := fmt.Fprintln(writer, _ORDER_PREFIX+strconv.FormatUint(r.id, 10)+_SPLITER+safeCall+_SPLITER+safeData)
	return err
}

func read(reader *bufio.Reader) (*sendObject, error) {
	line := ""
	for {
		p, prefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line += string(p)
		if prefix {
			continue
		}
		t := line
		line = ""
		if strings.HasPrefix(t, _ORDER_PREFIX) {
			parts := strings.Split(t[len(_ORDER_PREFIX):], _SPLITER)
			if len(parts) != 3 {
				continue
			}

			rsp := &sendObject{}
			rsp.id, err = strconv.ParseUint(parts[0], 10, 0)
			if err != nil {
				rsp.err = fmt.Errorf("decode id failed: %w", err)
				return rsp, nil
			}

			call, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				rsp.err = fmt.Errorf("decode call failed: %w", err)
				return rsp, nil
			}
			rsp.call = string(call)
			rsp.content, err = base64.StdEncoding.DecodeString(parts[2])
			if err != nil {
				rsp.err = fmt.Errorf("decode content failed: %w", err)
				return rsp, nil
			}
			return rsp, nil
		}
	}
}

type HandshakeInfo struct {
	Name    string
	Version string
}

func checkHandshake(handshake string, options *Options) bool {
	data, err := base64.StdEncoding.DecodeString(handshake)
	if err != nil {
		return false
	}
	info := &HandshakeInfo{}
	err = msgpack.Unmarshal(data, info)
	if err != nil {
		return false
	}
	if info.Name != options.HostName {
		return false
	}

	if info.Version < options.HostMinimalVersion {
		return false
	}

	return true
}
