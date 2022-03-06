package main

import (
	_ "GoliGoli/init"
	_ "GoliGoli/models"
	_ "GoliGoli/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}
