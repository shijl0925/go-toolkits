package gormx

import "errors"

// 使用说明:
//type UserService struct {
//	*gormx.ServiceImplement[repo.IUserRepo, model.User]
//}
//
//type IUserService interface {
//	gormx.IBaseService[model.User]
//	// 扩展用户相关业务方法
//	// ...
//}
//
//func NewIUserService() IUserService {
//	return &UserService{
//		gormx.ServiceImplement: gormx.NewServiceImplement[repo.IUserRepo, model.User](userRepo),
//	}
//}

// IBaseService 定义泛型基础接口
type IBaseService[T any] interface {
	Page(page, pageSize int, opts ...DBOption) ([]T, int64, error)
	ListByOpts(opts ...DBOption) ([]T, error)
	ListByIds(ids []int, opts ...DBOption) ([]T, error)
	ListByMap(columnMap map[string]interface{}, opts ...DBOption) ([]T, error)
	GetOneByOpts(opts ...DBOption) (*T, error)
	GetOneById(id int, opts ...DBOption) (*T, error)
	GetOneByMap(columnMap map[string]interface{}, opts ...DBOption) (*T, error)
	Exists(opts ...DBOption) bool
	Count(opts ...DBOption) (int64, error)
	Save(item *T, opts ...DBOption) error
	SaveOrUpdate(item *T, opts ...DBOption) error
	SaveBatch(items []*T, opts ...DBOption) error
	SaveInBatches(items []*T, batchSize int, opts ...DBOption) error
	Update(item *T, opts ...DBOption) error
	UpdateById(id int, vars map[string]interface{}, opts ...DBOption) error
	UpdateBatch(vars map[string]interface{}, opts ...DBOption) error
	Remove(item *T, opts ...DBOption) error
	RemoveBatch(opts ...DBOption) error
	RemoveById(id int, opts ...DBOption) error
	RemoveByIds(ids []int, opts ...DBOption) error
	RemoveByMap(columnMap map[string]interface{}, opts ...DBOption) error
	Upsert(item *T, vars map[string]interface{}, opts ...DBOption) error
	GetOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}, opts ...DBOption) (*T, error)
	UpdateOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}, opts ...DBOption) (*T, error)
	Sum(field string, opts ...DBOption) (float64, error)
	Max(field string, opts ...DBOption) (interface{}, error)
	Min(field string, opts ...DBOption) (interface{}, error)
	Avg(field string, opts ...DBOption) (float64, error)
	Increment(id int, field string, value interface{}, opts ...DBOption) error
	Decrement(id int, field string, value interface{}, opts ...DBOption) error
}

// ServiceImplement 通用服务实现
type ServiceImplement[Entity any, T any] struct {
	repo Entity
}

func NewServiceImplement[Entity any, T any](repo Entity) *ServiceImplement[Entity, T] {
	return &ServiceImplement[Entity, T]{repo: repo}
}

// Page 分页查询
func (s *ServiceImplement[Entity, T]) Page(page, pageSize int, opts ...DBOption) ([]T, int64, error) {
	if r, ok := any(s.repo).(interface {
		SelectPage(page, pageSize int, opts ...DBOption) ([]T, int64, error)
	}); ok {
		return r.SelectPage(page, pageSize, opts...)
	}
	return nil, 0, errors.New("repo does not implement SelectPage method")
}

// ListByOpts 根据条件查询
func (s *ServiceImplement[Entity, T]) ListByOpts(opts ...DBOption) ([]T, error) {
	if r, ok := any(s.repo).(interface {
		SelectListByOpts(opts ...DBOption) ([]T, error)
	}); ok {
		return r.SelectListByOpts(opts...)
	}
	return nil, errors.New("repo does not implement SelectListByOpts method")
}

// GetOneByOpts 根据条件查询单个
func (s *ServiceImplement[Entity, T]) GetOneByOpts(opts ...DBOption) (*T, error) {
	if r, ok := any(s.repo).(interface {
		SelectOneByOpts(opts ...DBOption) (T, error)
	}); ok {
		item, err := r.SelectOneByOpts(opts...)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, errors.New("repo does not implement SelectOneByOpts method")
}

// GetOneById 根据ID查询
func (s *ServiceImplement[Entity, T]) GetOneById(id int, opts ...DBOption) (*T, error) {
	if r, ok := any(s.repo).(interface {
		SelectOneById(id int, opts ...DBOption) (T, error)
	}); ok {
		item, err := r.SelectOneById(id, opts...)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, errors.New("repo does not implement SelectOneById method")
}

// GetOneByMap 根据条件查询单个
func (s *ServiceImplement[Entity, T]) GetOneByMap(columnMap map[string]interface{}, opts ...DBOption) (*T, error) {
	if r, ok := any(s.repo).(interface {
		SelectOneByMap(columnMap map[string]interface{}, opts ...DBOption) (T, error)
	}); ok {
		item, err := r.SelectOneByMap(columnMap, opts...)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, errors.New("repo does not implement SelectOneByMap method")
}

// Exists 判断记录是否存在
func (s *ServiceImplement[Entity, T]) Exists(opts ...DBOption) bool {
	if r, ok := any(s.repo).(interface {
		Exists(opts ...DBOption) bool
	}); ok {
		return r.Exists(opts...)
	}
	return false
}

// Save 保存
func (s *ServiceImplement[Entity, T]) Save(item *T, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		Insert(item *T, opts ...DBOption) error
	}); ok {
		return r.Insert(item, opts...)
	}
	return errors.New("repo does not implement Insert method")
}

// SaveOrUpdate 保存或更新
func (s *ServiceImplement[Entity, T]) SaveOrUpdate(item *T, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		InsertOrUpdate(item *T, opts ...DBOption) error
	}); ok {
		return r.InsertOrUpdate(item, opts...)
	}
	return errors.New("repo does not implement InsertOrUpdate method")
}

// SaveBatch 批量保存
func (s *ServiceImplement[Entity, T]) SaveBatch(items []*T, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		InsertBatch(items []*T, opts ...DBOption) error
	}); ok {
		return r.InsertBatch(items, opts...)
	}
	return errors.New("repo does not implement InsertBatch method")
}

// SaveInBatches 批量保存, 按批次处理
func (s *ServiceImplement[Entity, T]) SaveInBatches(items []*T, batchSize int, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		InsertInBatches(items []*T, batchSize int, opts ...DBOption) error
	}); ok {
		return r.InsertInBatches(items, batchSize, opts...)
	}
	return errors.New("repo does not implement InsertInBatches method")
}

// Update 更新
func (s *ServiceImplement[Entity, T]) Update(item *T, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		Update(item *T, opts ...DBOption) error
	}); ok {
		return r.Update(item, opts...)
	}
	return errors.New("repo does not implement Update method")
}

// UpdateById 根据ID更新
func (s *ServiceImplement[Entity, T]) UpdateById(id int, vars map[string]interface{}, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		UpdateById(id int, vars map[string]interface{}, opts ...DBOption) error
	}); ok {
		return r.UpdateById(id, vars, opts...)
	}
	return errors.New("repo does not implement UpdateById method")
}

// UpdateBatch 批量更新
func (s *ServiceImplement[Entity, T]) UpdateBatch(vars map[string]interface{}, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		UpdateByOpts(vars map[string]interface{}, opts ...DBOption) error
	}); ok {
		return r.UpdateByOpts(vars, opts...)
	}
	return errors.New("repo does not implement UpdateByOpts method")
}

// Remove 删除
func (s *ServiceImplement[Entity, T]) Remove(item *T, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		Delete(item *T, opts ...DBOption) error
	}); ok {
		return r.Delete(item, opts...)
	}
	return errors.New("repo does not implement Delete method")
}

// RemoveBatch 批量删除
func (s *ServiceImplement[Entity, T]) RemoveBatch(opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		DeleteByOpts(opts ...DBOption) error
	}); ok {
		return r.DeleteByOpts(opts...)
	}
	return errors.New("repo does not implement DeleteByOpts method")
}

// RemoveById 根据ID删除
func (s *ServiceImplement[Entity, T]) RemoveById(id int, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		DeleteById(id int, opts ...DBOption) error
	}); ok {
		return r.DeleteById(id, opts...)
	}
	return errors.New("repo does not implement DeleteById method")
}

// RemoveByIds 根据ID批量删除
func (s *ServiceImplement[Entity, T]) RemoveByIds(ids []int, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		DeleteBatchIds(ids []int, opts ...DBOption) error
	}); ok {
		return r.DeleteBatchIds(ids, opts...)
	}
	return errors.New("repo does not implement DeleteBatchIds method")
}

// RemoveByMap 根据条件删除
func (s *ServiceImplement[Entity, T]) RemoveByMap(columnMap map[string]interface{}, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		DeleteByMap(columnMap map[string]interface{}, opts ...DBOption) error
	}); ok {
		return r.DeleteByMap(columnMap, opts...)
	}
	return errors.New("repo does not implement DeleteByMap method")
}

// Upsert 插入或更新
func (s *ServiceImplement[Entity, T]) Upsert(item *T, vars map[string]interface{}, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		Upsert(item *T, vars map[string]interface{}, opts ...DBOption) error
	}); ok {
		return r.Upsert(item, vars, opts...)
	}
	return errors.New("repo does not implement Upsert method")
}

// GetOrCreate 获取或创建
func (s *ServiceImplement[Entity, T]) GetOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}, opts ...DBOption) (*T, error) {
	if r, ok := any(s.repo).(interface {
		GetOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}, opts ...DBOption) (*T, error)
	}); ok {
		return r.GetOrCreate(whereCond, assignAttrs, opts...)
	}
	return nil, errors.New("repo does not implement GetOrCreate method")
}

// UpdateOrCreate 不存在则创建，存在则更新
func (s *ServiceImplement[Entity, T]) UpdateOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}, opts ...DBOption) (*T, error) {
	if r, ok := any(s.repo).(interface {
		UpdateOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}, opts ...DBOption) (*T, error)
	}); ok {
		return r.UpdateOrCreate(whereCond, assignAttrs, opts...)
	}
	return nil, errors.New("repo does not implement UpdateOrCreate method")
}

// ListByIds 根据ID批量查询
func (s *ServiceImplement[Entity, T]) ListByIds(ids []int, opts ...DBOption) ([]T, error) {
	if r, ok := any(s.repo).(interface {
		SelectListByIds(ids []int, opts ...DBOption) ([]T, error)
	}); ok {
		return r.SelectListByIds(ids, opts...)
	}
	return nil, errors.New("repo does not implement SelectListByIds method")
}

// ListByMap 根据条件查询
func (s *ServiceImplement[Entity, T]) ListByMap(columnMap map[string]interface{}, opts ...DBOption) ([]T, error) {
	if r, ok := any(s.repo).(interface {
		SelectListByMap(columnMap map[string]interface{}, opts ...DBOption) ([]T, error)
	}); ok {
		return r.SelectListByMap(columnMap, opts...)
	}
	return nil, errors.New("repo does not implement SelectListByMap method")
}

// Count 统计数量
func (s *ServiceImplement[Entity, T]) Count(opts ...DBOption) (int64, error) {
	if r, ok := any(s.repo).(interface {
		SelectCount(opts ...DBOption) (int64, error)
	}); ok {
		return r.SelectCount(opts...)
	}
	return 0, errors.New("repo does not implement SelectCount method")
}

// Sum 求和特定字段
func (s *ServiceImplement[Entity, T]) Sum(field string, opts ...DBOption) (float64, error) {
	if r, ok := any(s.repo).(interface {
		Sum(field string, opts ...DBOption) (float64, error)
	}); ok {
		return r.Sum(field, opts...)
	}
	return 0, errors.New("repo does not implement Sum method")
}

// Max 获取字段最大值
func (s *ServiceImplement[Entity, T]) Max(field string, opts ...DBOption) (interface{}, error) {
	if r, ok := any(s.repo).(interface {
		Max(field string, opts ...DBOption) (interface{}, error)
	}); ok {
		return r.Max(field, opts...)
	}
	return nil, errors.New("repo does not implement Max method")
}

// Min 获取字段最小值
func (s *ServiceImplement[Entity, T]) Min(field string, opts ...DBOption) (interface{}, error) {
	if r, ok := any(s.repo).(interface {
		Min(field string, opts ...DBOption) (interface{}, error)
	}); ok {
		return r.Min(field, opts...)
	}
	return nil, errors.New("repo does not implement Min method")
}

// Avg 获取字段平均值
func (s *ServiceImplement[Entity, T]) Avg(field string, opts ...DBOption) (float64, error) {
	if r, ok := any(s.repo).(interface {
		Avg(field string, opts ...DBOption) (float64, error)
	}); ok {
		return r.Avg(field, opts...)
	}
	return 0, errors.New("repo does not implement Avg method")
}

// Increment 字段自增
func (s *ServiceImplement[Entity, T]) Increment(id int, field string, value interface{}, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		Increment(id int, field string, value interface{}, opts ...DBOption) error
	}); ok {
		return r.Increment(id, field, value, opts...)
	}
	return errors.New("repo does not implement Increment method")
}

// Decrement 字段自减
func (s *ServiceImplement[Entity, T]) Decrement(id int, field string, value interface{}, opts ...DBOption) error {
	if r, ok := any(s.repo).(interface {
		Decrement(id int, field string, value interface{}, opts ...DBOption) error
	}); ok {
		return r.Decrement(id, field, value, opts...)
	}
	return errors.New("repo does not implement Decrement method")
}
