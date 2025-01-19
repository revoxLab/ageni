package dal

import (
	xgorm "github.com/readonme/open-studio/common/gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var StudioDB *gorm.DB

func Init(c *xgorm.Config) {

	StudioDB = xgorm.NewORM(c)
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
	})
	StudioDB.Logger = newLogger
	if StudioDB == nil {
		panic("init db failed")
	}
}
