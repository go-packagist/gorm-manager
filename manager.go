package gormer

import (
	"gorm.io/gorm"
	"sync"
)

type Config struct {
	Default     string
	Connections map[string]ConnectionFunc
}

type ConnectionFunc func() *gorm.DB

type Manager struct {
	config  *Config
	reloved map[string]*gorm.DB
	rw      sync.RWMutex
}

func New(config *Config) *Manager {
	return NewManager(config)
}

func NewManager(config *Config) *Manager {
	return &Manager{
		config:  config,
		reloved: make(map[string]*gorm.DB),
		rw:      sync.RWMutex{},
	}
}

func (m *Manager) Connection(name ...string) *gorm.DB {
	if len(name) > 0 {
		return m.resolve(name[0])
	}

	return m.resolve(m.config.Default)
}

func (m *Manager) resolve(name string) *gorm.DB {
	if db, ok := m.reloved[name]; ok {
		return db
	}

	if _, ok := m.config.Connections[name]; !ok {
		panic("connection " + name + " is not defined")
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	m.reloved[name] = m.config.Connections[name]()

	return m.reloved[name]
}
