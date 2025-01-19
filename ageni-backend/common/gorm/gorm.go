package gorm

import (
	"time"

	"github.com/readonme/open-studio/common/log"
	"github.com/readonme/open-studio/common/xtime"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	xgorm "gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Config struct {
	DebugMode    bool           `toml:"debug_mode"`
	DriverName   string         `toml:"driver_name"`
	DSN          string         `toml:"dsn"`
	ReadDSN      []string       `toml:"read_dsn"`
	MaxIdleConns int            `toml:"max_idle_conns"`
	MaxOpenConns int            `toml:"max_open_conns"`
	MaxLifetime  xtime.Duration `toml:"max_lifetime"`
	MaxIdleTime  xtime.Duration `toml:"max_idle_time"`
}

func NewORM(c *Config) (orm *xgorm.DB) {
	orm, err := Open(c)
	if err != nil {
		log.Errorf("gorm open mysql error(%v)", err)
		panic(err)
	}
	return
}

func Open(c *Config) (*xgorm.DB, error) {
	orm, err := xgorm.Open(mysql.New(mysql.Config{
		DriverName:                c.DriverName,
		DSN:                       c.DSN, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,   // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,  // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    false, // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // smart configure based on used version
	}))
	if err != nil {
		log.Errorf("xgorm.Open error %v", err)
		return nil, err
	}
	if len(c.ReadDSN) > 0 {
		replicas := []gorm.Dialector{}
		for k, v := range c.ReadDSN {
			cfg := mysql.Config{
				DSN: v,
			}
			log.Infof("read-write-%d-%s", k, v)
			replicas = append(replicas, mysql.New(cfg))
		}
		err = orm.Use(
			dbresolver.Register(dbresolver.Config{
				Sources: []gorm.Dialector{mysql.New(mysql.Config{
					DSN: c.DSN,
				})},
				Replicas: replicas,
				Policy:   dbresolver.RandomPolicy{},
			}).
				SetMaxIdleConns(c.MaxIdleConns).
				SetConnMaxLifetime(time.Duration(c.MaxLifetime)).
				SetMaxOpenConns(c.MaxOpenConns).
				SetConnMaxIdleTime(time.Duration(c.MaxIdleTime)),
		)
		if err != nil {
			log.Errorf("orm.Use error %v", err)
		}
	} else {
		db, err := orm.DB()
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(c.MaxIdleConns)
		db.SetMaxOpenConns(c.MaxOpenConns)
		db.SetConnMaxLifetime(time.Duration(c.MaxLifetime))
		db.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime))
	}
	return orm, nil
}
