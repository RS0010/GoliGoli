package main

import (
	"GoliGoli/models"
	_ "GoliGoli/routers"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	user1 := &models.User{Username: "admin", Password: "admin"}
	user2 := &models.User{Username: "admins", Password: "admin"}
	user1.Add()
	user2.Add()
	//video := &models.Video{Title: "hello"}
	//video.Add()
	//fmt.Println(user1)
	//fmt.Println(user2)
	err := user1.Block(user2)
	//err := user2.Block(user1)
	fmt.Println(user1)
	fmt.Println(user2)
	fmt.Println(*user2.Blocked[0])
	fmt.Println(*user2.Blocks[0])
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(*user1.Blocks[0])
	beego.Run()
}
