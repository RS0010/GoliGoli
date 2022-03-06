package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BarrageList []*Comment

type BarrageFilter struct {
	ID        uint
	VideoID   uint
	UserID    uint
	Page      int
	PageSize  int
	OrderBy   string
	OrderType string
	filtered  bool
}

func (f *BarrageFilter) filter(db *gorm.DB) *gorm.DB {
	if f.ID != 0 {
		return db.Where("id = ?", f.ID).Preload(clause.Associations)
	}
	db = db.Where("video_id = ?", f.VideoID)
	if f.UserID != 0 {
		db = db.Where("user_id = ?", f.UserID)
	}
	db = db.Preload(clause.Associations)
	if commentColumns.Contains(f.OrderBy) {
		if f.OrderType != "desc" {
			f.OrderType = "asc"
		}
		db = db.Order(f.OrderBy + " " + f.OrderType)
	}
	return db
}

func (f *BarrageFilter) paginate(db *gorm.DB) *gorm.DB {
	if f.Page > 0 && f.PageSize > 0 {
		return db.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize)
	}
	return db.Offset(0).Limit(20)
}

func (f *BarrageFilter) Query(list *CommentList) error {
	return f.paginate(f.filter(o)).Find(list).Error
}

func (f *BarrageFilter) Search() (CommentList, error) {
	var list CommentList
	if err := f.Query(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (f *BarrageFilter) Count() (total int64, err error) {
	m := f.filter(o).Count(&total)
	return total, m.Error
}
