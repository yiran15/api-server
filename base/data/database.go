package data

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/conf"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, func(), error) {
	dsn, err := conf.GetMysqlDsn()
	if err != nil {
		return nil, nil, err
	}
	var dbLogger logger.Interface
	// 开启mysql日志
	if viper.GetBool("mysql.debug") || conf.GetLogLevel() == "debug" {
		zap.S().Debug("enable debug mode on the database")
		dbLogger = logger.Default.LogMode(logger.Info)
	}

	dbInstance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   dbLogger,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("exception in initializing mysql database, %w", err)
	}

	// 确保数据库连接已建立
	sqlDB, err := dbInstance.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to obtain database connection, %w", err)
	}

	// 尝试Ping数据库以确保连接有效
	err = sqlDB.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to obtain database connection, %w", err)
	}

	sqlDB.SetMaxOpenConns(conf.GetMysqlMaxOpenConns())
	sqlDB.SetMaxIdleConns(conf.GetMysqlMaxIdleConns())
	sqlDB.SetConnMaxLifetime(conf.GetMysqlMaxLifetime())

	zap.S().Info("db connect success")
	return dbInstance, func() { _ = sqlDB.Close() }, nil
}
