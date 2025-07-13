package store

import (
	"context"

	"gorm.io/gorm"
)

// TxManagerInterface 事务管理器接口。
type TxManagerInterface interface {
	// Transaction 在一个事务中执行给定的函数。
	// 函数内部可以通过 GetTX(ctx) 获取事务 DB 实例，从而确保所有操作都在同一事务内。
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// TxManager 是 TxManagerInterface 的实现。
type TxManager struct {
	db *gorm.DB // 主数据库连接
}

// NewTxManager 创建一个新的 TxManager 实例。
func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

// Transaction 执行一个数据库事务。
// 如果 fn 返回错误，事务将回滚；否则，事务将提交。
func (s *TxManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 将事务 DB 实例放入新的上下文，传递给回调函数
		return fn(context.WithValue(ctx, &txStruct{}, tx))
	})
}
