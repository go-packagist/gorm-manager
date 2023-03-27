package gormer

import "gorm.io/gorm"

type Connection interface{}

type MySQLConfig struct {
	DSN        string
	GormConfig *gorm.Config
	GormUsees  []gorm.Plugin
}

var _ Connection = (*MySQLConfig)(nil)

type Config struct {
	Default     string
	Connections map[string]Connection
}
