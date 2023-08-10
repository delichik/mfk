package logger

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/delichik/mfk/config"
)

const (
	ModuleName = "logger"
)

var defaultLoggerConfig = &Config{
	Level:     "info",
	Format:    "text",
	LogPath:   "logs/logger.log",
	MaxSize:   50,
	MaxBackup: 1,
	MaxAge:    30,
	Compress:  true,
	logLevel:  zapcore.InfoLevel,
}

func init() {
	config.RegisterModuleConfig(ModuleName, defaultLoggerConfig)
}

type Config struct {
	Level     string `yaml:"level" comment:"Min log output level"`
	Format    string `yaml:"format" comment:"Log output format: json|text"`
	LogPath   string `yaml:"log_path" comment:"Path to write log, use \"stdout\" to write to console"`
	MaxSize   int    `yaml:"max_size" comment:"Maximum size (MB) of a log file"`
	MaxBackup int    `yaml:"max_backup" comment:"Maximum count of log backup"`
	MaxAge    int    `yaml:"max_age" comment:"Maximum saving days of a log backup"`
	Compress  bool   `yaml:"compress" comment:"Compress the backups"`

	LogDriver zapcore.WriteSyncer `yaml:"-"`
	logLevel  zapcore.Level       `yaml:"-"`
}

func (c *Config) Check() error {
	if c.Level != "" {
		err := c.logLevel.UnmarshalText([]byte(c.Level))
		if err != nil {
			return fmt.Errorf("level: %w", err)
		}
	}

	if c.Format != "" && c.Format != "json" && c.Format != "text" {
		return errors.New(`format: must be "json" or "text"`)
	}

	if len(c.LogPath) == 0 && c.LogDriver == nil {
		return errors.New(`log_path: required`)
	}

	if c.LogPath != "stdout" && c.LogDriver == nil {
		if c.MaxSize <= 0 {
			return errors.New(`max_size: must grater than 0`)
		}

		if c.MaxBackup < 0 {
			return errors.New(`max_backup: must grater than or equal to 0`)
		}

		if c.MaxAge <= 0 {
			return errors.New(`max_backup: must grater than 0`)
		}
	}
	return nil
}

func (c *Config) Clone() config.ModuleConfig {
	newObj := *c
	newObj.logLevel = c.logLevel
	return &newObj
}

func (c *Config) Compare(config.ModuleConfig) bool {
	return true
}

var defaultLogger *zap.Logger

var loggers = make(map[string]*zap.Logger)

func InitDefault(c *Config) {
	defaultLogger = Init(ModuleName, c)
}

func Init(name string, c *Config) *zap.Logger {
	var writeSyncer zapcore.WriteSyncer
	if c.LogDriver != nil {
		writeSyncer = c.LogDriver
	} else if c.LogPath == "stdout" {
		writeSyncer = os.Stdout
	} else {
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   c.LogPath,
			MaxSize:    c.MaxSize,
			MaxBackups: c.MaxBackup,
			MaxAge:     c.MaxAge,
			Compress:   c.Compress,
		})
	}
	var core zapcore.Core
	if c.Format == "json" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), writeSyncer, c.logLevel)
	} else if c.Format == "text" {
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), writeSyncer, c.logLevel)
	}
	l := zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
	loggers[name] = l
	return l
}

func InitDefaultByConfigSet(c config.ConfigSet) {
	var loggerConfig *Config
	t := c.GetModuleConfig(ModuleName)
	if t == nil || t.(*Config) == nil {
		loggerConfig = defaultLoggerConfig
	} else {
		loggerConfig = t.(*Config)
	}
	defaultLogger = Init(ModuleName, loggerConfig)
}

func GetDefaultConfig() *Config {
	return defaultLoggerConfig.Clone().(*Config)
}

func InitByConfigSet(name string, c config.ConfigSet) *zap.Logger {
	var loggerConfig *Config
	t := c.GetModuleConfig(name)
	if t == nil || t.(*Config) == nil {
		loggerConfig = defaultLoggerConfig
	} else {
		loggerConfig = t.(*Config)
	}

	return Init(name, loggerConfig)
}
