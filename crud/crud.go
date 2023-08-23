package crud

import (
	"gorm.io/gorm"
)

type Query[Model interface{}] struct {
	m       *Model
	db      *gorm.DB
	page    int
	size    int
	maxSize int
}

func (q *Query[Model]) List(scopes ...func(*gorm.DB) *gorm.DB) ([]*Model, error) {
	offset := (q.page - 1) * q.size
	limit := q.size
	if limit > q.maxSize {
		limit = q.maxSize
	}
	list := make([]*Model, 0, limit)

	return list, q.db.Model(q.m).Scopes(scopes...).Offset(offset).Limit(limit).Find(&list).Error
}

func (q *Query[Model]) Count(scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var cnt int64
	return cnt, q.db.Model(q.m).Scopes(scopes...).Count(&cnt).Error
}

func NewQuery[Model interface{}](m *Model, db *gorm.DB, page, size int) *Query[Model] {
	return &Query[Model]{m: m, db: db, maxSize: 50, page: page, size: size}
}
