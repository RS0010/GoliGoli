package models

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
)

var o orm.Ormer

func init() {
	if err := orm.RegisterDataBase("default", "sqlite3", "./models/models.db"); err != nil {
		panic(err)
	}
	orm.RegisterModel(new(UserUserBlocks))
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Video))
	orm.RegisterModel(new(Comment))
	orm.RegisterModel(new(Barrage))
	if err := orm.RunSyncdb("default", true, false); err != nil {
		panic(err)
	}
	o = orm.NewOrm()
}
