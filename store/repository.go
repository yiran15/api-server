package store

import (
	"context"

	// 假设你的错误包路径
	"github.com/yiran15/api-server/base/log"
	"go.uber.org/zap"
)

// repository 是 Repository 接口的 GORM 实现。
type repository[T any] struct {
	DBProviderInterface // 嵌入 DB 提供者接口，以获取 DB 实例
}

// NewRepository 创建一个新的通用 Repository 实例。
// 它接收一个 DBProviderInterface，使得 Repository 能够获取到正确的 GORM DB 实例（可能是主 DB 或事务 DB）。
func NewRepository[T any](dbProvider DBProviderInterface) *repository[T] {
	return &repository[T]{DBProviderInterface: dbProvider}
}

// Create 插入单个对象到数据库。
func (r *repository[T]) Create(ctx context.Context, obj *T) error {
	// r.getDB(ctx, obj) 会根据上下文获取主 DB 或事务 DB，并初始化 Model
	if err := r.getDB(ctx, obj).Create(obj).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to create object", zap.Error(err), zap.Any("obj", obj))
		return err
	}
	return nil
}

// CreateBatch 批量插入多个对象到数据库。
func (r *repository[T]) CreateBatch(ctx context.Context, objs []*T) error {
	if len(objs) == 0 {
		return nil // 如果没有对象，直接返回
	}
	// 这里用第一个对象来初始化 Model，因为批量插入的所有对象类型都相同
	if err := r.getDB(ctx, objs[0]).Create(objs).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to create objects in batch", zap.Error(err), zap.Any("objs_count", len(objs)))
		return err
	}
	return nil
}

// Update 更新对象。
// 可以通过 opts 指定更新条件（如 Where），或者直接使用 obj 的主键进行更新。
func (r *repository[T]) Update(ctx context.Context, obj *T, opts ...Option) error {
	db := r.getDB(ctx, obj, opts...)
	if err := db.Updates(obj).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to update object", zap.Error(err), zap.Any("obj", obj))
		return err
	}
	return nil
}

// Delete 删除对象。
// 可以通过 opts 指定删除条件（如 Where），或者直接使用 obj 的主键进行软删除/硬删除。
func (r *repository[T]) Delete(ctx context.Context, obj *T, opts ...Option) error {
	db := r.getDB(ctx, obj, opts...)
	if err := db.Delete(obj).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to delete object", zap.Error(err), zap.Any("obj", obj))
		return err
	}
	return nil
}

// DeleteBatch 批量删除对象。
// 可以通过 opts 指定删除条件（如 Where），或者直接使用 obj 的主键进行软删除/硬删除。
func (r *repository[T]) DeleteBatch(ctx context.Context, objs []*T, opts ...Option) error {
	// db := r.getDB(ctx, objs[0], opts...) 会在应用 opts 之前，先 db.Model(objs)
	// 如果 opts 中没有 Where 条件，GORM 会根据 obj 的主键来删除。
	// 如果 obj 为零值，且 opts 中有 Where 条件，可以实现条件批量删除。
	db := r.getDB(ctx, objs[0], opts...)
	if err := db.Delete(objs).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to delete objects in batch", zap.Error(err), zap.Any("objs_count", len(objs)))
		return err
	}
	return nil
}

// Query 查询单个对象。
// 如果未找到记录，则返回 apierr.NotFoundErr。
func (r *repository[T]) Query(ctx context.Context, opts ...Option) (*T, error) {
	model := new(T)                    // 创建一个 T 类型的零值实例，用于 GORM 的 First 方法
	db := r.getDB(ctx, model, opts...) // 先应用所有查询选项
	if err := db.First(model).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to query object", zap.Error(err))
		return nil, err
	}
	return model, nil
}

// List 查询对象列表并支持分页。
// 返回总记录数、对象列表和错误。
func (r *repository[T]) List(ctx context.Context, page, pageSize int, colum, oder string, opts ...Option) (total int64, objs []*T, err error) {
	model := new(T)
	// baseDB 应用所有筛选条件和预加载，用于计数和实际查询
	baseDB := r.getDB(ctx, model, opts...)

	// 1. 先计算总数
	// 使用 Clone() 避免后续分页和排序操作影响 Count
	if err = baseDB.Count(&total).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to count objects for list", zap.Error(err))
		return 0, nil, err
	}

	if total == 0 {
		return 0, []*T{}, nil // 如果没有记录，直接返回空列表和总数0
	}

	// 2. 应用分页和排序到新的 DB 实例上
	listDB := baseDB
	if page > 0 && pageSize > 0 {
		if colum != "" && oder != "" {
			listDB = listDB.Order(colum + " " + oder)
		}
		offset := (page - 1) * pageSize
		listDB = listDB.Offset(offset).Limit(pageSize)
	}

	// 3. 执行查询
	if err = listDB.Find(&objs).Error; err != nil {
		log.WithRequestID(ctx).Error("failed to list objects", zap.Error(err))
		return 0, nil, err
	}
	return total, objs, nil
}

func (r *repository[T]) AppendAssociation(ctx context.Context, model *T, objName string, obj any) error {
	if err := r.getDB(ctx, model).Association(objName).Append(obj); err != nil {
		log.WithRequestID(ctx).Error("failed to append association", zap.Error(err), zap.Any("obj", obj))
		return err
	}
	return nil
}

func (r *repository[T]) ReplaceAssociation(ctx context.Context, model *T, objName string, obj any) error {
	if err := r.getDB(ctx, model).Association(objName).Replace(obj); err != nil {
		log.WithRequestID(ctx).Error("failed to replace association", zap.Error(err), zap.Any("obj", obj))
		return err
	}
	return nil
}

func (r *repository[T]) ClearAssociation(ctx context.Context, model *T, objName string) error {
	if err := r.getDB(ctx, model).Association(objName).Clear(); err != nil {
		log.WithRequestID(ctx).Error("failed to clear association", zap.Error(err))
		return err
	}
	return nil
}
