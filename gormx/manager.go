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
//  - Set: 设置关联关系（事务内清空并重建）
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

// Set 设置关联关系（事务内清空并重建）。
//
// GORM 的 Association.Replace 在 many2many join 表缺少组合唯一约束、
// 且 join 表包含软删除字段时，可能无法保证重复调用的幂等性。这里显式
// Clear 后再 Append，使 Set 的行为符合“整体替换”语义。
func (m *GenericAssociationManager[T, R]) Set(relations []R) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		assoc := tx.Model(&m.owner).Association(m.fieldName)
		if err := assoc.Error; err != nil {
			return err
		}
		if err := assoc.Clear(); err != nil {
			return err
		}
		if len(relations) == 0 {
			return nil
		}
		return assoc.Append(relations)
	})
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
