package store

import (
	"context"

	"gorm.io/gorm"
)

// txStruct 是一个私有的空结构体，用作 context.WithValue 的键，以存储事务 DB 实例。
type txStruct struct{}

// GetTX 从上下文中获取 GORM 事务 DB 实例。
// 如果上下文中没有事务 DB，则返回 nil。
func GetTX(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(&txStruct{}).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

// DBProviderInterface 数据库连接提供接口。
type DBProviderInterface interface {
	// getDB 根据上下文和可选的 GORM 选项获取 GORM DB 实例。
	// 如果上下文中存在事务，则返回事务 DB，否则返回主 DB。
	// model 参数用于初始化 db.Model()。
	getDB(ctx context.Context, model any, opts ...Option) *gorm.DB
}

// DBProvider 是 DBProviderInterface 的实现。
type DBProvider struct {
	db *gorm.DB // 主数据库连接，通常是全局单例
}

// NewDBProvider 创建一个新的 DBProvider 实例。
func NewDBProvider(db *gorm.DB) *DBProvider {
	return &DBProvider{db: db}
}

// getDB 获取数据库实例，优先使用事务中的 DB。
func (d *DBProvider) getDB(ctx context.Context, model any, opts ...Option) (db *gorm.DB) {
	// 尝试从上下文中获取事务 DB
	if db = GetTX(ctx); db == nil {
		// 如果没有事务，使用主 DB
		db = d.db
	}

	// 如果提供了模型，则初始化 Model
	if model != nil {
		db = db.Model(model)
	}

	// 应用所有 GORM 选项
	for _, opt := range opts {
		db = opt(db)
	}
	return db.WithContext(ctx) // 将上下文绑定到 DB 实例，确保上下文取消时数据库操作也能终止
}
