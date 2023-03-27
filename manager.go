package gormer

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type DB struct {
	*gorm.DB
	Err error
}

type Manager struct {
	config   *Config
	resolved map[string]*DB
	rw       sync.RWMutex
}

type Option func(*Manager) *Manager

func NewManager(config *Config, opts ...Option) *Manager {
	m := &Manager{
		config:   config,
		resolved: make(map[string]*DB, len(config.Connections)),
		rw:       sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithInstance(m *Manager) *Manager {
	SetInstance(m)

	return m
}

func (m *Manager) Connect(name ...string) *DB {
	if len(name) > 0 {
		return m.resolve(name[0])
	}

	return m.resolve(m.config.Default)
}

func (m *Manager) resolve(name string) *DB {
	m.rw.Lock()
	defer m.rw.Unlock()

	if db, ok := m.resolved[name]; ok {
		return db
	}

	if _, ok := m.config.Connections[name]; !ok {
		return &DB{Err: errors.New("connection " + name + " is not defined")}
	}

	var (
		db  *gorm.DB
		err error
	)

	switch m.config.Connections[name].(type) {
	case *MySQLConfig:
		db, err = m.createMySQLConnection(m.config.Connections[name].(*MySQLConfig))
		break
	default:
		return &DB{Err: errors.New("connection " + name + " is not defined")}
	}

	if err != nil {
		return &DB{Err: err}
	}

	m.resolved[name] = &DB{DB: db}

	return m.resolved[name]
}

func (m *Manager) createMySQLConnection(config *MySQLConfig) (*gorm.DB, error) {
	// open connection
	db, err := gorm.Open(mysql.Open(config.DSN), config.GormConfig)
	if err != nil {
		return nil, err
	}

	// use plugins
	for _, use := range config.GormUsees {
		err := db.Use(use)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
