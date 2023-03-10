package logger

import (
	"errors"
	"fmt"
	"os"

	"github.com/delichik/my-go-pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ModuleLogger = "logger"
)

var defaultLoggerConfig = &LoggerConfig{
	Level:     "info",
	Format:    "text",
	LogPath:   "module.log",
	MaxSize:   50,
	MaxBackup: 1,
	MaxAge:    30,
	Compress:  true,
	logLevel:  zapcore.InfoLevel,
}

func init() {
	config.RegisterModuleConfig(ModuleLogger, defaultLoggerConfig)
}

type LoggerConfig struct {
	Level     string `yaml:"level" comment:"Min log output level"`
	Format    string `yaml:"format" comment:"Log output format: json|text"`
	LogPath   string `yaml:"log_path" comment:"Path to write log, use \"stdout\" to write to console"`
	MaxSize   int    `yaml:"max_size" comment:"Maximum size (MB) of a log file"`
	MaxBackup int    `yaml:"max_backup" comment:"Maximum count of log backup"`
	MaxAge    int    `yaml:"max_age" comment:"Maximum saving days of a log backup"`
	Compress  bool   `yaml:"compress" comment:"Compress the backups"`

	logLevel zapcore.Level `yaml:"-"`
}

func (c *LoggerConfig) Check() error {
	if c.Level != "" {
		err := c.logLevel.UnmarshalText([]byte(c.Level))
		if err != nil {
			return fmt.Errorf("level: %w", err)
		}
	}

	if c.Format != "" && c.Format != "json" && c.Format != "text" {
		return errors.New(`format: must be "json" or "text"`)
	}

	if len(c.LogPath) == 0 {
		return errors.New(`log_path: required`)
	}

	if c.LogPath != "stdout" {
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

func (c *LoggerConfig) Clone() config.ModuleConfig {
	newObj := *c
	newObj.logLevel = c.logLevel
	return &newObj
}

func (c *LoggerConfig) Compare(config.ModuleConfig) bool {
	return true
}

var defaultLogger *zap.Logger

func Init(c config.ConfigSet) {
	var loggerConfig *LoggerConfig
	t := c.GetModuleConfig(ModuleLogger)
	if t == nil || t.(*LoggerConfig) == nil {
		loggerConfig = defaultLoggerConfig
	} else {
		loggerConfig = t.(*LoggerConfig)
	}
	var writeSyncer zapcore.WriteSyncer
	if loggerConfig.LogPath == "stdout" {
		writeSyncer = os.Stdout
	} else {
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   loggerConfig.LogPath,
			MaxSize:    loggerConfig.MaxSize,
			MaxBackups: loggerConfig.MaxBackup,
			MaxAge:     loggerConfig.MaxAge,
			Compress:   loggerConfig.Compress,
		})
	}
	var core zapcore.Core
	if loggerConfig.Format == "json" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), writeSyncer, loggerConfig.logLevel)
	} else if loggerConfig.Format == "text" {
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), writeSyncer, loggerConfig.logLevel)
	}
	defaultLogger = zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller())
}

func NewLogger(moduleName string) *zap.Logger {
	return defaultLogger.With(zap.String("module", moduleName))
}
