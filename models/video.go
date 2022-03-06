package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Video struct {
	ID            uint           `gorm:"primarykey"`
	CreatedAt     time.Time      `gorm:""`
	UpdatedAt     time.Time      `gorm:""`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Path          string         `gorm:"type:varchar(255);unique_index"`
	Title         string         `gorm:""`
	Info          string         `gorm:""`
	Category      string         `gorm:""`
	ViewCount     int            `gorm:""`
	State         int            `gorm:""`
	LikerList     []*User        `gorm:"many2many:user_likes"`
	CollectorList []*User        `gorm:"many2many:user_collects"`
	SharerList    []*User        `gorm:"many2many:user_shares"`
	AuthorList    []*User        `gorm:"many2many:user_videos"`
	CommentList   []*Comment     `gorm:""`
	BarrageList   []*Barrage     `gorm:""`
	UserID        uint           `gorm:""`
}

func (v Video) Render(endpoint string) Video {
	if v.Path != "" {
		v.Path = endpoint + v.Path
	}
	return v
}

func (v *Video) Check() (exist bool, err error) {
	var count int64
	err = o.Model(v).Where("id = ?", v.ID).Count(&count).Error
	return count > 0, err
}

func (v *Video) Query() error {
	return o.First(v, v.ID).Error
}

func (v *Video) Update() error {
	return o.Save(v).Error
}

func (v *Video) Delete() error {
	return o.Select(clause.Associations).Delete(v).Error
}

func (v *Video) Post(authors ...*User) error {
	if len(authors) > 0 {
		v.UserID = authors[0].ID
	}
	return o.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(v).Error; err != nil {
			return err
		}
		if len(authors) > 0 {
			for _, author := range authors {
				if err := tx.Model(author).Association("PostList").Append(v); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (v *Video) Unpost() error {
	return o.Select(clause.Associations).Delete(v).Error
}

func (v *Video) IsLiker(u *User) (bool, error) {
	var likerList []*User
	err := o.Model(v).Association("LikerList").Find(&likerList, u.ID)
	return len(likerList) > 0, err
}

func (v *Video) QueryLiker() error {
	return o.Model(v).Association("LikerList").Find(&v.LikerList)
}

func (v *Video) CountLiker() int {
	return int(o.Model(v).Association("LikerList").Count())
}

func (v *Video) IsCollector(u *User) (bool, error) {
	var collectorList []*User
	err := o.Model(v).Association("CollectorList").Find(&collectorList, u.ID)
	return len(collectorList) > 0, err
}

func (v *Video) QueryCollector() error {
	return o.Model(v).Association("CollectorList").Find(&v.CollectorList)
}

func (v *Video) CountCollector() int {
	return int(o.Model(v).Association("CollectorList").Count())
}

func (v *Video) IsSharer(u *User) (bool, error) {
	var sharerList []*User
	err := o.Model(v).Association("SharerList").Find(&sharerList, u.ID)
	return len(sharerList) > 0, err
}

func (v *Video) QuerySharer() error {
	return o.Model(v).Association("SharerList").Find(&v.SharerList)
}

func (v *Video) CountSharer() int {
	return int(o.Model(v).Association("SharerList").Count())
}

func (v *Video) IsAuthor(u *User) (bool, error) {
	var authorList []*User
	err := o.Model(v).Association("AuthorList").Find(&authorList, u.ID)
	return len(authorList) > 0, err
}

func (v *Video) QueryAuthor(filters ...UserFilter) error {
	if len(filters) > 0 {
		return o.Where("ID", v.ID).Preload("AuthorList", func(db *gorm.DB) *gorm.DB {
			return filters[0].filter(db)
		}).Find(v).Error
	}
	return o.Model(v).Association("AuthorList").Find(&v.AuthorList)
}

func (v *Video) CountAuthor(filter ...UserFilter) (total int64, err error) {
	if len(filter) > 0 {
		err := o.Where("ID", v.ID).Preload("AuthorList", func(db *gorm.DB) *gorm.DB {
			m := filter[0].filter(db)
			m.Model(&User{}).Count(&total)
			return m
		}).Find(&Video{}).Error
		return total, err
	}
	return o.Model(v).Association("AuthorList").Count(), nil
}

func (v *Video) QueryComment() error {
	return o.Model(v).Association("CommentList").Find(&v.CommentList)
}

func (v *Video) CountComment() int {
	return int(o.Model(v).Association("CommentList").Count())
}

func (v *Video) QueryBarrage() error {
	return o.Model(v).Association("BarrageList").Find(&v.BarrageList)
}

func (v *Video) CountBarrage() int {
	return int(o.Model(v).Association("BarrageList").Count())
}
