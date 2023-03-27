package gormer

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"testing"
	"time"
)

func TestGormer(t *testing.T) {
	db := newManager()

	assert.Nil(t, db.Connect("db1").Err())
	assert.Nil(t, db.Connect("db2").Err())
}

func newManager() *Manager {
	writeDSN := DSN{
		Host: "127.0.0.1",
		Port: 3306,
		User: "gorm",
		Pass: "gorm",
		Db:   "gorm",
	}

	readDSN := DSN{
		Host: "127.0.0.1",
		Port: 3306,
		User: "gorm",
		Pass: "gorm",
		Db:   "gorm",
	}

	return NewManager(&Config{
		Default: "db1",
		Connections: map[string]Connection{
			"db1": &MySQLConfig{
				DSN: writeDSN.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true, // 禁用默认事务
				},
				GormUsees: []gorm.Plugin{
					dbresolver.Register(
						dbresolver.Config{
							Sources:  []gorm.Dialector{mysql.Open(writeDSN.String())},
							Replicas: []gorm.Dialector{mysql.Open(readDSN.String())},
							Policy:   dbresolver.RandomPolicy{},
						}).
						SetConnMaxLifetime(time.Hour).
						SetConnMaxIdleTime(24 * time.Hour).
						SetMaxIdleConns(100).
						SetMaxOpenConns(200),
				},
			},
			"db2": &MySQLConfig{
				DSN: writeDSN.String(),
				GormConfig: &gorm.Config{
					SkipDefaultTransaction: true, // 禁用默认事务
				},
			},
		},
	})
}
