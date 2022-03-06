package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Comment struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:"" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:""`
	Time      time.Time      `gorm:"" json:"-"`
	Replies   []*Comment     `gorm:"foreignkey:Parent"`
	Parent    int            `gorm:""`
	State     int            `gorm:"" json:"-"`
	UserID    uint           `gorm:"" json:"-"`
	VideoID   uint           `gorm:"" json:"-"`
}

func (c *Comment) Check() (exist bool, err error) {
	var count int64
	err = o.Model(&Comment{}).Where("id", c.ID).Count(&count).Error
	return count > 0, err
}

func (c *Comment) Query() error {
	return o.First(c, c.ID).Error
}

func (c *Comment) Update() error {
	return o.Save(c).Error
}

func (c *Comment) Delete() error {
	return o.Select(clause.Associations).Delete(c).Error
}
