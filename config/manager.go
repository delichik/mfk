package config

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type ReloadCallback func(name string, config ModuleConfig)

type Manager struct {
	path string
	ctx  context.Context
	wg   sync.WaitGroup

	locker         sync.RWMutex
	c              *Config
	lastLoadC      *Config
	reloadCallback ReloadCallback
}

func NewManager(ctx context.Context, path string) *Manager {
	return &Manager{
		path: path,
		ctx:  ctx,
	}
}

func (m *Manager) Init() error {
	newConfig, err := Load(m.path)
	if err != nil {
		return err
	}

	m.locker.Lock()
	m.lastLoadC = newConfig
	m.c = &Config{moduleConfigs: map[string]ModuleConfig{}}
	for moduleName, moduleConfig := range newConfig.moduleConfigs {
		m.c.moduleConfigs[moduleName] = moduleConfig.Clone()
	}
	m.locker.Unlock()

	err = Save(m.path, m.c)
	if err != nil {
		return err
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.watchConfig()
	}()
	return nil
}

func (m *Manager) watchConfig() {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			m.locker.Lock()
			newConfig, err := Load(m.path)
			if err != nil {
				log.Println("Config file load failed:", err.Error(), ", abort to load")
				m.locker.Unlock()
				continue
			}
			changed := false
			for name, lmc := range m.lastLoadC.moduleConfigs {
				nmc := newConfig.GetModuleConfig(name)
				if !nmc.Compare(lmc) {
					changed = true
					if m.reloadCallback != nil {
						m.reloadCallback(name, nmc)
					}
				}
			}
			if changed {
				m.c = &Config{moduleConfigs: map[string]ModuleConfig{}}
				for moduleName, moduleConfig := range newConfig.moduleConfigs {
					m.c.moduleConfigs[moduleName] = moduleConfig.Clone()
				}
			}

			m.locker.Unlock()
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *Manager) GetModuleConfig(moduleName string) ModuleConfig {
	return m.lastLoadC.moduleConfigs[moduleName]
}

func (m *Manager) ModifyModuleConfig(moduleName string, call func(ModuleConfig)) error {
	mc, ok := m.c.moduleConfigs[moduleName]
	if !ok {
		return errors.New("module config not found")
	}
	m.locker.Lock()
	defer m.locker.Unlock()
	call(mc)
	err := Save(m.path, m.c)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) SetReloadCallback(callback ReloadCallback) {
	m.reloadCallback = callback
}
