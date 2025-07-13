package store

import (
	"fmt"

	"gorm.io/gorm"
)

// Option 函数用于配置 GORM 查询。
// 它接收一个 *gorm.DB 实例并返回修改后的 *gorm.DB 实例，
// 从而允许链式调用 GORM 方法。
//
// @example:
//
//	options := []store.Option{
//	    store.Where("id", req.ID),
//	    store.Where("status", model.UserStatusEnable),
//	    store.Preload("Profile"),
//	    store.Order("created_at desc"),
//	}
type Option func(db *gorm.DB) *gorm.DB

// Where 用于添加 WHERE 条件。
func Where(colum string, value any) Option {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s = ?", colum), value)
	}
}

func In(colum string, values any) Option {
	return func(query *gorm.DB) *gorm.DB {
		return query.Where(fmt.Sprintf("%s in (?)", colum), values)
	}
}

// Preload 用于预加载关联模型。
func Preload(modelName string, args ...any) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(modelName, args...)
	}
}

// Select 用于指定查询的列。
func Select(columns ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(columns)
	}
}

// Order 用于指定排序条件。
func Order(value string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}

// Limit 用于限制查询结果的数量。
func Limit(limit int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(limit)
	}
}

// Offset 用于设置查询的偏移量。
func Offset(offset int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset)
	}
}

// Scopes 用于应用 GORM Scopes。
func Scopes(funcs ...func(*gorm.DB) *gorm.DB) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(funcs...)
	}
}
