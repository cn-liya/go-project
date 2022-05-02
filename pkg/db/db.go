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

type Config struct {
	Address  string
	Username string
	Password string
	Database string
	MaxOpen  int
	MaxIdle  int
	TraceLog bool
}

func NewMysql(cfg *Config) *gorm.DB {
	dsn := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Address + ")/" + cfg.Database +
		"?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
	return newDB(mysql.Open(dsn), cfg)
}

//func NewPostgres(cfg *Config) *gorm.DB {
//	host, port, _ := net.SplitHostPort(cfg.Address)
//	dsn := "host=" + host + " port=" + port + " user=" + cfg.Username + " password=" + cfg.Password +
//		" dbname=" + cfg.Database + " sslmode=disable TimeZone=Asia/Shanghai"
//	return newDB(postgres.Open(dsn), cfg)
//}
//
//func NewSqlserver(cfg *Config) *gorm.DB {
//	dsn := "sqlserver://" + cfg.Username + ":" + cfg.Password +
//		"@" + cfg.Address + "?database=" + cfg.Database
//	return newDB(sqlserver.Open(dsn), cfg)
//}
//
//func NewClickhouse(cfg *Config) *gorm.DB {
//	dsn := "tcp://" + cfg.Address + "?database=" + cfg.Database +
//		"&username=" + cfg.Username + "&password=" + cfg.Password
//	return newDB(clickhouse.Open(dsn), cfg)
//}

func newDB(dial gorm.Dialector, cfg *Config) *gorm.DB {
	opt := &gorm.Config{}
	if cfg.TraceLog {
		opt.Logger = &gormLog{glog.Discard}
	} else {
		opt.Logger = glog.Discard.LogMode(glog.Silent)
	}
	orm, err := gorm.Open(dial, opt)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := orm.DB()
	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	//sqlDB.SetConnMaxLifetime(time.Hour)
	//sqlDB.SetConnMaxIdleTime(time.Minute)
	return orm
}

type gormLog struct {
	glog.Interface
}

func (*gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rows int64), err error) {
	sql, rows := fc()
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.FromContext(ctx).Error("gorm", sql, err)
	} else {
		logger.FromContext(ctx).Trace("gorm", sql, rows, begin)
	}
}
