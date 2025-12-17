package gormx

import (
	"errors"
	"fmt"
	"github.com/shijl0925/go-toolkits/stringx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

// 使用说明
//type UserRepo struct {
//	gormx.BaseRepo[model.User] // 继承基础仓库功能
//}
//
//type IUserRepo interface {
//	gormx.IBaseRepo[model.User] // 嵌入基础接口
//
//	// 扩展用户相关业务方法
//	// ...
//}
//
//func NewIUserRepo() IUserRepo {
//	return &UserRepo{}
//}

// IBaseRepo 定义基础接口
type IBaseRepo[T any] interface {
	SelectOneById(id int) (T, error)
	SelectOneByOpts(opts ...DBOption) (T, error)
	SelectOneByMap(columnMap map[string]interface{}) (T, error)
	SelectListByIds(ids []int) ([]T, error)
	SelectListByOpts(opts ...DBOption) ([]T, error)
	SelectListByMap(columnMap map[string]interface{}) ([]T, error)
	SelectPage(page, pageSize int, opts ...DBOption) ([]T, int64, error)
	SelectCount(opts ...DBOption) (int64, error)
	Exists(opts ...DBOption) bool

	Insert(item *T) error
	InsertBatch(items []*T) error
	InsertInBatches(items []*T, batchSize int) error
	InsertOrUpdate(item *T) error

	Update(item *T) error
	UpdateById(id int, vars map[string]interface{}) error
	UpdateByOpts(vars map[string]interface{}, opts ...DBOption) error

	Upsert(item *T, vars map[string]interface{}) error
	GetOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}) (*T, error)
	UpdateOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}) (*T, error)

	Delete(item *T) error
	DeleteById(id int) error
	DeleteBatchIds(ids []int) error
	DeleteByOpts(opts ...DBOption) error
	DeleteByMap(columnMap map[string]interface{}) error

	Sum(field string, opts ...DBOption) (float64, error)
	Max(field string, opts ...DBOption) (interface{}, error)
	Min(field string, opts ...DBOption) (interface{}, error)
	Avg(field string, opts ...DBOption) (float64, error)

	Increment(id int, field string, value interface{}) error
	Decrement(id int, field string, value interface{}) error
}

type BaseRepo[T any] struct{}

// SelectOneByOpts 根据条件查询单条记录
func (r *BaseRepo[T]) SelectOneByOpts(opts ...DBOption) (T, error) {
	var items []T
	db := GetDb(opts...).Model(new(T))

	// 使用 Limit(2) 来检测是否有多条记录
	err := db.Limit(2).Find(&items).Error

	if err != nil {
		var zero T
		return zero, err
	}

	size := len(items)
	if size == 0 {
		var zero T
		return zero, gorm.ErrRecordNotFound
	} else if size > 1 {
		var zero T

		errMessage := fmt.Sprintf("expected one result (or null) to be returned by SelectOneByOpts(), but found:  %d", size)
		return zero, errors.New(errMessage)
	}

	return items[0], nil
}

// SelectOneById 根据ID查询单条记录
func (r *BaseRepo[T]) SelectOneById(id int) (T, error) {
	var item T
	err := globalDb.Model(new(T)).Where("id = ?", id).First(&item).Error

	return item, err
}

func (r *BaseRepo[T]) SelectOneByMap(columnMap map[string]interface{}) (T, error) {
	var items []T
	db := GetDb().Model(new(T))

	if len(columnMap) > 0 {
		db = db.Where(columnMap)
	}

	// 使用 Limit(2) 来检测是否有多条记录
	err := db.Limit(2).Find(&items).Error

	if err != nil {
		var zero T
		return zero, err
	}

	size := len(items)
	if size == 0 {
		var zero T
		return zero, gorm.ErrRecordNotFound
	} else if size > 1 {
		var zero T

		errMessage := fmt.Sprintf("expected one result (or null) to be returned by SelectOneByMap(), but found:  %d", size)
		return zero, errors.New(errMessage)
	}

	return items[0], nil
}

// SelectListByIds 根据ID批量查询
func (r *BaseRepo[T]) SelectListByIds(ids []int) ([]T, error) {
	var items []T
	err := globalDb.Model(new(T)).Where("id in ?", ids).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// SelectListByOpts 根据条件查询所有记录
func (r *BaseRepo[T]) SelectListByOpts(opts ...DBOption) ([]T, error) {
	var items []T

	db := GetDb(opts...).Model(new(T))

	err := db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// SelectPage 分页查询
func (r *BaseRepo[T]) SelectPage(page, pageSize int, opts ...DBOption) ([]T, int64, error) {
	var items []T

	db := GetDb(opts...).Model(new(T))
	query := db.Limit(pageSize).Offset((page - 1) * pageSize)

	// 应用条件查询获取分页数据
	err := query.Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// SelectListByMap 根据条件查询
func (r *BaseRepo[T]) SelectListByMap(columnMap map[string]interface{}) ([]T, error) {
	var items []T
	db := GetDb().Model(new(T))

	if len(columnMap) > 0 {
		db = db.Where(columnMap)
	}

	//// 将map转换为结构体
	//if len(columnMap) > 0 {
	//	structInstance, err := utils.MapToStruct[T](columnMap)
	//	if err != nil {
	//		return nil, err
	//	}
	//	db = db.Where(structInstance)
	//}

	err := db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// SelectCount 统计数量
func (r *BaseRepo[T]) SelectCount(opts ...DBOption) (int64, error) {
	var count int64
	db := GetDb(opts...).Model(new(T))
	return count, db.Count(&count).Error
}

// Exists 判断记录是否存在
func (r *BaseRepo[T]) Exists(opts ...DBOption) bool {
	var item T
	db := GetDb(opts...).Model(new(T))
	err := db.Limit(1).First(&item).Error
	return err == nil
}

// Insert 创建单条记录
func (r *BaseRepo[T]) Insert(item *T) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	return globalDb.Create(item).Error
}

// InsertBatch 批量插入
func (r *BaseRepo[T]) InsertBatch(items []*T) error {
	if len(items) == 0 {
		return errors.New("items cannot be empty")
	}

	return globalDb.Create(items).Error
}

// InsertInBatches 批量插入, 按批次插入
func (r *BaseRepo[T]) InsertInBatches(items []*T, batchSize int) error {
	if len(items) == 0 {
		return errors.New("items cannot be empty")
	}

	if batchSize <= 0 {
		return errors.New("batchSize must be greater than 0")
	}

	return globalDb.CreateInBatches(items, batchSize).Error
}

// InsertOrUpdate 创建或更新, 若主键有值, 则更新, 否则创建
func (r *BaseRepo[T]) InsertOrUpdate(item *T) error {
	return globalDb.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(item).Error
}

// getAssociationFields 通过反射获取模型中的关联字段
func (r *BaseRepo[T]) getAssociationFields() map[string]bool {
	restrictedFields := make(map[string]bool)

	// 创建模型实例
	var instance T
	modelType := reflect.TypeOf(instance)

	// 处理指针类型
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// 递归处理所有字段（包括嵌套结构体）
	r.processFields(modelType, restrictedFields)

	return restrictedFields
}

// processFields 递归处理字段（包括嵌套结构体）
func (r *BaseRepo[T]) processFields(structType reflect.Type, restrictedFields map[string]bool) {
	// 遍历所有字段
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// 跳过未导出字段
		if !field.IsExported() {
			continue
		}

		// 处理匿名嵌套结构体
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			// 递归处理匿名结构体的字段
			r.processFields(field.Type, restrictedFields)
			continue
		}

		// 检查是否为关联字段
		if r.isAssociationField(field) {
			fieldName := r.getFieldName(field)
			restrictedFields[fieldName] = true
		}
	}
}

// isAssociationField 判断是否为关联字段
func (r *BaseRepo[T]) isAssociationField(field reflect.StructField) bool {
	// 检查字段类型是否为切片或结构体（通常是关联）
	fieldType := field.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// 关联字段通常为切片类型（一对多、多对多）或结构体类型（一对一）
	if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Struct {
		// 检查是否有 GORM 关联标签
		gormTag := field.Tag.Get("gorm")
		if gormTag != "" {
			// 检查是否为忽略迁移的字段
			if strings.Contains(gormTag, "-:migration") {
				return true
			}

			// 检查是否包含关联相关的关键字
			associationKeywords := []string{"foreignkey", "many2many", "belongsto", "hasone", "hasmany"}
			for _, keyword := range associationKeywords {
				if strings.Contains(strings.ToLower(gormTag), keyword) {
					return true
				}
			}
		}
	}

	return false
}

// getFieldName 获取字段名（优先使用 JSON 标签，其次使用字段名）
func (r *BaseRepo[T]) getFieldName(field reflect.StructField) string {
	// 优先使用 JSON 标签中的名称
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		// 如果是忽略标签，则使用数据库列名或蛇形命名
		if jsonTag == "-" {
			return stringx.ToSnake(field.Name)
		}

		if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
			return jsonTag[:commaIndex]
		}
		return jsonTag
	}

	// 默认使用字段名的蛇形命名
	return stringx.ToSnake(field.Name)
}

// filterUpdateFields 过滤不允许更新的字段
func (r *BaseRepo[T]) filterUpdateFields(vars map[string]interface{}) map[string]interface{} {
	// 定义不允许更新的字段
	restrictedFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
	}

	// 动态获取模型的关联字段
	associationFields := r.getAssociationFields()
	for field := range associationFields {
		restrictedFields[field] = true
	}

	// 过滤变量
	filtered := make(map[string]interface{}, len(vars))
	for key, value := range vars {
		if !restrictedFields[key] {
			filtered[key] = value
		}
	}

	return filtered
}

// Update 更新单条记录
func (r *BaseRepo[T]) Update(item *T) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}
	// 使用Updates 而不是Save, 避免意外插入新纪录
	return globalDb.Model(item).Updates(item).Error
}

// UpdateById 根据ID更新记录
func (r *BaseRepo[T]) UpdateById(id int, vars map[string]interface{}) error {
	// 过滤不允许更新的字段
	filteredVars := r.filterUpdateFields(vars)
	return globalDb.Model(new(T)).Where("id = ?", id).Updates(filteredVars).Error
}

// UpdateByOpts 根据条件更新记录
func (r *BaseRepo[T]) UpdateByOpts(vars map[string]interface{}, opts ...DBOption) error {
	// 检查是否提供了更新条件
	if len(opts) == 0 {
		return errors.New("cannot update records without conditions, please provide where conditions")
	}

	db := GetDb(opts...).Model(new(T))
	// 过滤不允许更新的字段
	filteredVars := r.filterUpdateFields(vars)
	return db.Updates(filteredVars).Error
}

// Upsert 插入或更新（根据唯一约束）
func (r *BaseRepo[T]) Upsert(item *T, vars map[string]interface{}) error {
	return globalDb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(vars),
	}).Create(item).Error
}

// GetOrCreate 查找第一条匹配的记录，否则根据条件创建新记录
func (r *BaseRepo[T]) GetOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}) (*T, error) {
	var item T
	db := GetDb().Model(new(T))

	// 构造查询条件
	if len(whereCond) > 0 {
		db = db.Where(whereCond)
	}

	result := db.FirstOrCreate(&item, assignAttrs)
	if result.Error != nil {
		return nil, result.Error
	}

	return &item, nil
}

// UpdateOrCreate 不存在则创建，存在则更新
func (r *BaseRepo[T]) UpdateOrCreate(whereCond map[string]interface{}, assignAttrs map[string]interface{}) (*T, error) {
	var item T
	db := GetDb().Model(new(T))

	// FirstOrCreate 会先查找匹配 whereCond 的记录
	// 如果找不到，则使用 whereCond + assignAttrs 创建新记录
	// 如果找到了，则使用 assignAttrs 更新现有记录
	result := db.Where(whereCond).Assign(assignAttrs).FirstOrCreate(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	return &item, nil
}

// Delete 删除单条记录
func (r *BaseRepo[T]) Delete(item *T) error {
	return globalDb.Delete(item).Error
}

// DeleteById 根据ID删除记录
func (r *BaseRepo[T]) DeleteById(id int) error {
	return globalDb.Where("id = ?", id).Delete(new(T)).Error
}

// DeleteBatchIds 根据ID批量删除记录
func (r *BaseRepo[T]) DeleteBatchIds(ids []int) error {
	return globalDb.Where("id in ?", ids).Delete(new(T)).Error
}

// DeleteByOpts 根据条件删除记录
func (r *BaseRepo[T]) DeleteByOpts(opts ...DBOption) error {
	db := GetDb(opts...)

	// 如果opts为空，不执行删除操作
	if len(opts) == 0 {
		return errors.New("delete operation requires where conditions")
	}

	return db.Delete(new(T)).Error
}

// DeleteByMap 根据条件删除记录
func (r *BaseRepo[T]) DeleteByMap(columnMap map[string]interface{}) error {
	db := GetDb().Model(new(T))

	// 如果columnMap为空，不执行删除操作
	if len(columnMap) == 0 {
		return errors.New("cannot delete records without conditions")
	}

	return db.Where(columnMap).Delete(new(T)).Error
}

// Sum 求和特定字段
func (r *BaseRepo[T]) Sum(field string, opts ...DBOption) (float64, error) {
	// 验证字段名合法性
	if field == "" {
		return 0, errors.New("field name cannot be empty")
	}

	var result float64
	db := GetDb(opts...).Model(new(T))
	err := db.Select("SUM(" + field + ")").Scan(&result).Error
	return result, err
}

// Max 获取字段最大值
func (r *BaseRepo[T]) Max(field string, opts ...DBOption) (interface{}, error) {
	// 验证字段名合法性
	if field == "" {
		return nil, errors.New("field name cannot be empty")
	}

	var result interface{}
	db := GetDb(opts...).Model(new(T))

	// 使用原生 SQL 方式处理聚合查询
	row := db.Select("MAX(" + field + ")").Row()
	if row == nil {
		return nil, errors.New("failed to execute query")
	}

	err := row.Scan(&result)

	return result, err
}

// Min 获取字段最小值
func (r *BaseRepo[T]) Min(field string, opts ...DBOption) (interface{}, error) {
	// 验证字段名合法性
	if field == "" {
		return nil, errors.New("field name cannot be empty")
	}

	var result interface{}
	db := GetDb(opts...).Model(new(T))

	// 使用原生 SQL 方式处理聚合查询
	row := db.Select("MIN(" + field + ")").Row()
	if row == nil {
		return nil, errors.New("failed to execute query")
	}

	err := row.Scan(&result)

	return result, err
}

// Avg 获取字段平均值
func (r *BaseRepo[T]) Avg(field string, opts ...DBOption) (float64, error) {
	// 验证字段名合法性
	if field == "" {
		return 0, errors.New("field name cannot be empty")
	}

	var result float64
	db := GetDb(opts...).Model(new(T))
	err := db.Select("AVG(" + field + ")").Scan(&result).Error
	return result, err
}

// Increment 字段自增
func (r *BaseRepo[T]) Increment(id int, field string, value interface{}) error {
	return globalDb.Model(new(T)).Where("id = ?", id).UpdateColumn(field, gorm.Expr(field+" + ?", value)).Error
}

// Decrement 字段自减
func (r *BaseRepo[T]) Decrement(id int, field string, value interface{}) error {
	return globalDb.Model(new(T)).Where("id = ?", id).UpdateColumn(field, gorm.Expr(field+" - ?", value)).Error
}
