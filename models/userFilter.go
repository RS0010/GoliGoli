package models

import "gorm.io/gorm"

type UserList []*User

type UserFilter struct {
	ID        uint
	Name      string
	Email     string
	VideoID   uint
	Page      int
	PageSize  int
	OrderBy   string
	OrderType string
	filtered  bool
}

func (f *UserFilter) filter(db *gorm.DB) *gorm.DB {
	if f.ID != 0 {
		return db.Where("id = ?", f.ID)
	}
	if f.Name != "" {
		db = db.Where("username like ?", "%"+f.Name+"%")
	}
	if f.Email != "" {
		db = db.Where("email = ?", f.Email)
	}
	if userColumns.Contains(f.OrderBy) {
		if f.OrderType != "desc" {
			f.OrderType = "asc"
		}
		db = db.Order(f.OrderBy + " " + f.OrderType)
	}
	return db
}

func (f *UserFilter) paginate(db *gorm.DB) *gorm.DB {
	if f.Page > 0 && f.PageSize > 0 {
		return db.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize)
	}
	return db.Offset(0).Limit(20)
}

func (f *UserFilter) Query(list *UserList) error {
	if f.VideoID != 0 {
		video := Video{ID: f.VideoID}
		err := video.QueryAuthor(*f)
		*list = video.AuthorList
		return err
	}
	return f.paginate(f.filter(o)).Find(list).Error
}

func (f *UserFilter) Search() (UserList, error) {
	var list UserList
	if err := f.Query(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (f *UserFilter) Count() (total int64, err error) {
	if f.VideoID > 0 {
		video := &Video{ID: f.VideoID}
		total, err = video.CountAuthor(*f)
		return total, err
	}
	m := f.filter(o).Count(&total)
	return total, m.Error
}
