package utils

import (
	"fmt"
	"github.com/beego/beego/v2/core/utils"
	beego "github.com/beego/beego/v2/server/web"
	"reflect"
	"runtime"
)

func Display(variables ...interface{}) {
	if beego.BConfig.RunMode == "dev" {
		var data []interface{}
		for index, item := range variables {
			_, file, line, _ := runtime.Caller(1)
			data = append(data, fmt.Sprintf(
				"--------------------------------------------------\n%s:%d#%d\n%s",
				file, line, index, reflect.TypeOf(item).String()), item)
		}
		utils.Display(data...)
		fmt.Println("--------------------------------------------------")
	}
}
