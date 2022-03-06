package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Barrage struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:""`
	UpdatedAt time.Time      `gorm:"" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:""`
	PostTime  time.Time      `gorm:"" json:"-"`
	PlayTime  time.Time      `gorm:""`
	State     int            `gorm:"" json:"-"`
	UserID    uint           `gorm:""`
	VideoID   uint           `gorm:""`
}

func (b *Barrage) Check() (exist bool, err error) {
	var count int64
	err = o.Model(&Barrage{}).Where("id", b.ID).Count(&count).Error
	return count > 0, err
}

func (b *Barrage) Query() error {
	return o.First(b, b.ID).Error
}

func (b *Barrage) Update() error {
	return o.Save(b).Error
}

func (b *Barrage) Delete() error {
	return o.Select(clause.Associations).Delete(b).Error
}
