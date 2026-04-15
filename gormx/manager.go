package gormx

import (
	"gorm.io/gorm"
)

// 说明:
// 1. 创建适用于一对一、一对多、多对多关系管理器接口
//  - 提供了标准化的API来管理模型间的各种关联关系
//  - 支持一对一、一对多、多对多关系操作
// 2. 关联数据操作
//  - Add: 添加关联关系
//  - Remove: 删除关联关系
//  - Clear: 清空所有关联
//  - Set: 设置关联关系（替换现有关系）
//  - Count: 统计关联数量
//  - All: 获取所有关联数据

// AssociationManager 建立适用于一对一、一对多、多对多关系管理器接口
type AssociationManager[T any, R any] interface {
	Add(relations ...R) error
	Remove(relations ...R) error
	Clear() error
	Set(relations []R) error
	Count() (int64, error)
	All() ([]R, error)
}

// GenericAssociationManager 通用关系管理器
type GenericAssociationManager[T any, R any] struct {
	owner     T
	fieldName string
	db        *gorm.DB
}

// NewAssociationManager 创建关系管理器，可通过 opts 注入请求级 DB 或事务（例如 UseDB(tx)）
func NewAssociationManager[T any, R any](owner T, fieldName string, opts ...DBOption) AssociationManager[T, R] {
	return &GenericAssociationManager[T, R]{
		owner:     owner,
		fieldName: fieldName,
		db:        GetDb(opts...),
	}
}

// Add 添加关联关系
func (m *GenericAssociationManager[T, R]) Add(relations ...R) error {
	assoc := m.db.Model(&m.owner).Association(m.fieldName)
	if err := assoc.Error; err != nil {
		return err
	}
	return assoc.Append(relations)
}

// Remove 删除关联关系
func (m *GenericAssociationManager[T, R]) Remove(relations ...R) error {
	assoc := m.db.Model(&m.owner).Association(m.fieldName)
	if err := assoc.Error; err != nil {
		return err
	}
	return assoc.Delete(relations)
}

// Clear 清空所有关联
func (m *GenericAssociationManager[T, R]) Clear() error {
	assoc := m.db.Model(&m.owner).Association(m.fieldName)
	if err := assoc.Error; err != nil {
		return err
	}
	return assoc.Clear()
}

// Set 设置关联关系（替换现有关系）
func (m *GenericAssociationManager[T, R]) Set(relations []R) error {
	assoc := m.db.Model(&m.owner).Association(m.fieldName)
	if err := assoc.Error; err != nil {
		return err
	}
	return assoc.Replace(relations)
}

// Count 统计关联数量
func (m *GenericAssociationManager[T, R]) Count() (int64, error) {
	assoc := m.db.Model(&m.owner).Distinct().Association(m.fieldName)
	if err := assoc.Error; err != nil {
		return 0, err
	}
	return assoc.Count(), nil
}

// All 获取所有关联数据
func (m *GenericAssociationManager[T, R]) All() ([]R, error) {
	assoc := m.db.Model(&m.owner).Distinct().Association(m.fieldName)
	if err := assoc.Error; err != nil {
		return nil, err
	}

	var relations []R
	if err := assoc.Find(&relations); err != nil {
		return nil, err
	}
	return relations, nil
}
