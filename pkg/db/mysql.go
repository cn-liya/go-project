package db

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"log"
	"project/pkg/logger"
	"time"
)

type Mysql struct {
	Username string
	Password string
	Address  string
	Database string
	MaxOpen  int
	MaxIdle  int
}

func NewMysqlDB(cfg *Mysql) *gorm.DB {
	dsn := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Address + ")/" + cfg.Database +
		"?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: &gormLog{glog.Discard},
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := orm.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return orm
}

type gormLog struct {
	glog.Interface
}

const msg = "gorm"

func (*gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rows int64), err error) {
	sql, rows := fc()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.FromContext(ctx).Error(msg, sql, err)
	} else {
		logger.FromContext(ctx).Trace(msg, sql, rows, begin)
	}
}
