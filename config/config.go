package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"gopkg.in/yaml.v3"

	myyaml "github.com/delichik/my-go-pkg/yaml"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ConfigSet interface {
	// GetModuleConfig 获取 Config。
	GetModuleConfig(moduleName string) ModuleConfig
}

type ModuleConfig interface {
	// Check 预检查错误
	Check() error
	// Clone 深度拷贝内容并返回
	Clone() ModuleConfig
	// Compare 对比两个同类型配置是否一样, 一样则返回 true
	Compare(ModuleConfig) bool
}

var moduleDefaultConfigs = map[string]ModuleConfig{}

type Config struct {
	moduleConfigs map[string]ModuleConfig
}

func (c *Config) GetModuleConfig(moduleName string) ModuleConfig {
	return c.moduleConfigs[moduleName]
}

func (c *Config) String() string {
	b, err := myyaml.MarshallWithComments(c.moduleConfigs)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func RegisterModuleConfig(moduleName string, defaultConfig ModuleConfig) {
	_, ok := moduleDefaultConfigs[moduleName]
	if ok {
		panic(errors.New("module name existed"))
	}

	moduleDefaultConfigs[moduleName] = defaultConfig
}

func Load(path string) (*Config, error) {
	config := &Config{
		moduleConfigs: map[string]ModuleConfig{},
	}
	for moduleName, moduleConfig := range moduleDefaultConfigs {
		config.moduleConfigs[moduleName] = moduleConfig.Clone()
	}

	b, err := ioutil.ReadFile(path)
	if err == nil {
		t := map[string]yaml.Node{}
		err = yaml.Unmarshal(b, &t)
		if err != nil {
			return nil, err
		}

		for k, v := range t {
			c, ok := config.moduleConfigs[k]
			if !ok {
				continue
			}
			_ = v.Decode(c)
		}
	}

	err = Save(path, config)
	if err != nil {
		return nil, err
	}

	for moduleName, moduleConfig := range config.moduleConfigs {
		err = moduleConfig.Check()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", moduleName, err)
		}
	}

	return config, nil
}

func Save(path string, config *Config) error {
	b, err := myyaml.MarshallWithComments(config.moduleConfigs)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0655)
	if err != nil {
		return err
	}
	return nil
}
