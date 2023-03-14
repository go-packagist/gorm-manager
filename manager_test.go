package gormer

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestGormer(t *testing.T) {
	db := New(&Config{
		Default: "db1",
		Connections: map[string]ConnectionFunc{
			"db1": func() *gorm.DB {
				return nil
			},
			"db2": func() *gorm.DB {
				return nil
			},
		},
	})

	assert.Equal(t, db.Connection("db1"), db.Connection("db1"))
	assert.Equal(t, db.Connection(), db.Connection("db1"))

	assert.Panics(t, func() {
		db.Connection("db3")
	})
}
