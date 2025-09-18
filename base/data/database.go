package data

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/conf"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type zapWriter struct {
	log *zap.Logger
}

func (w zapWriter) Printf(_ string, args ...interface{}) {
	// gorm: l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	if len(args) == 5 {
		w.log.Error("gorm log", zap.String("err", args[1].(error).Error()), zap.Float64("elapsed", args[2].(float64)), zap.String("rows", "-"), zap.String("sql", args[4].(string)))
		return
	}
	switch args[2].(type) {
	case int64:
		w.log.Info("gorm log", zap.Float64("elapsed", args[1].(float64)), zap.Int64("rows", args[2].(int64)), zap.String("sql", args[3].(string)))
	case string:
		w.log.Info("gorm log", zap.Float64("elapsed", args[1].(float64)), zap.String("rows", args[2].(string)), zap.String("sql", args[3].(string)))
	default:
		w.log.Info("gorm log", zap.Float64("elapsed", args[1].(float64)), zap.String("rows", "-"), zap.String("sql", args[3].(string)))
	}
}

// newGormLogger create a new gorm logger
func newGormLogger(z *zap.Logger) logger.Interface {
	return logger.New(
		zapWriter{log: z}, // 使用 zap writer
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,
		},
	)
}

func NewDB() (*gorm.DB, func(), error) {
	dsn, err := conf.GetMysqlDsn()
	if err != nil {
		return nil, nil, err
	}
	var dbLogger logger.Interface
	// 开启mysql日志
	if viper.GetBool("mysql.debug") || conf.GetLogLevel() == "debug" {
		dbLogger = newGormLogger(zap.L())
		zap.S().Info("enable debug mode on the database")
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
