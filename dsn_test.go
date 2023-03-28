package gormer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDsn(t *testing.T) {
	dsn := DSN{
		Host: "127.0.0.1",
		Port: 3306,
		User: "root",
		Pass: "password",
		DB:   "gorm",
	}

	t.Run("dsn string", func(t *testing.T) {
		assert.Equal(t, "root:password@tcp(127.0.0.1:3306)/gorm", dsn.String())
	})

	t.Run("dsn string with options", func(t *testing.T) {
		dsn.Options = map[string]string{
			"charset":   "utf8mb4",
			"parseTime": "True",
		}

		assert.Equal(t, "root:password@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True", dsn.String())
	})

	t.Run("dsn string func", func(t *testing.T) {
		assert.Equal(t, "root:password@tcp(127.0.0.1:3306)/gorm", DSN{
			Host: "127.0.0.1",
			Port: 3306,
			User: "root",
			Pass: "password",
			DB:   "gorm",
		}.String())
	})
}
