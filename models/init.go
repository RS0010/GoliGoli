package models

import (
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	columns []string
)

var (
	o              *gorm.DB
	userColumns    columns
	videoColumns   columns
	commentColumns columns
	barrageFilter  columns
)

func init() {
	if db, err := gorm.Open(sqlite.Open("data/database/models.db"), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		o = db
	}
	if err := o.AutoMigrate(&User{}, &Video{}, &Comment{}, &Barrage{}); err != nil {
		panic(err)
	}
	userColumns.GetModelColumns(&User{})
	videoColumns.GetModelColumns(&Video{})
	commentColumns.GetModelColumns(&Comment{})
	barrageFilter.GetModelColumns(&Barrage{})
}

func (c *columns) GetModelColumns(model interface{}) {
	types, err := o.Migrator().ColumnTypes(model)
	if err != nil {
		panic(err)
	}
	for _, value := range types {
		*c = append(*c, value.Name())
	}
}

func (c *columns) Contains(field string) bool {
	for _, value := range videoColumns {
		if value == field {
			return true
		}
	}
	return false
}
