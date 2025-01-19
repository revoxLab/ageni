package gorm

import (
	"testing"
	"time"

	"github.com/readonme/open-studio/common/xtime"
	"github.com/stretchr/testify/assert"
)

func TestOpenGorm(t *testing.T) {
	db := NewORM(&Config{
		DebugMode:    true,
		DriverName:   "mysql",
		DSN:          "test:testpasw@(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local",
		MaxIdleConns: 16,
		MaxOpenConns: 128,
		MaxLifetime:  xtime.Duration(time.Second * 10),
		MaxIdleTime:  xtime.Duration(time.Second * 15),
	})
	sqlDb, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	err = sqlDb.Ping()
	assert.NoError(t, err, "create db connection error!", err)
}
