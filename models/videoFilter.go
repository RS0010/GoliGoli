package models

import "gorm.io/gorm"

type VideoList []*Video

type VideoFilter struct {
	Title     string
	Category  string
	UserID    uint
	PostID    uint
	Page      int
	PageSize  int
	OrderBy   string
	OrderType string
	filtered  bool
}

func (f *VideoFilter) filter(m *gorm.DB) *gorm.DB {
	if f.Title != "" {
		m = m.Where("title LIKE ?", "%"+f.Title+"%")
		f.filtered = true
	}
	if f.Category != "" {
		m = m.Where("category = ?", f.Category)
		f.filtered = true
	}
	if f.PostID > 0 {
		m = m.Where("user_id = ?", f.PostID)
		f.filtered = true
	}
	if videoColumns.Contains(f.OrderBy) {
		if f.OrderType != "desc" {
			f.OrderType = "asc"
		}
		m = m.Order(f.OrderBy + " " + f.OrderType)
	}
	return m
}

func (f *VideoFilter) paginate(m *gorm.DB) *gorm.DB {
	if f.Page > 0 && f.PageSize > 0 {
		return m.Offset((f.Page - 1) * f.PageSize).Limit(f.PageSize)
	}
	return m.Offset(0).Limit(20)
}

func (f *VideoFilter) Query(list *VideoList) error {
	if f.UserID > 0 {
		user := User{ID: f.UserID}
		err := user.QueryPost(*f)
		*list = user.PostList
		return err
	}
	return f.paginate(f.filter(o)).Find(list).Error
}

func (f *VideoFilter) Search() (VideoList, error) {
	var list VideoList
	if err := f.Query(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (f *VideoFilter) Count() (total int64, err error) {
	if f.UserID > 0 {
		user := &User{ID: f.UserID}
		total, err = user.CountPost(*f)
		return total, err
	}
	m := f.filter(o).Count(&total)
	return total, m.Error
}
