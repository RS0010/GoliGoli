package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Query",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           `/`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "PartUpdate",
			Router:           `/`,
			AllowHTTPMethods: []string{"patch"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Ban",
			Router:           `/:uid:int/ban`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Unban",
			Router:           `/:uid:int/ban`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Block",
			Router:           `/:uid:int/block`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Unblock",
			Router:           `/:uid:int/block`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Logout",
			Router:           `/login`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:UserController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Register",
			Router:           `/register`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Search",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Upload",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Query",
			Router:           `/:vid:int`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Update",
			Router:           `/:vid:int`,
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "PartUpdate",
			Router:           `/:vid:int`,
			AllowHTTPMethods: []string{"patch"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Delete",
			Router:           `/:vid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Ban",
			Router:           `/:vid:int/ban`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Unban",
			Router:           `/:vid:int/ban`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Barrage",
			Router:           `/:vid:int/barrage`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "QueryBarrage",
			Router:           `/:vid:int/barrage/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "UnBarrage",
			Router:           `/:vid:int/barrage/:bid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "UnCollect",
			Router:           `/:vid:int/collect`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Collect",
			Router:           `/:vid:int/collect`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Comment",
			Router:           `/:vid:int/comment`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "QueryComment",
			Router:           `/:vid:int/comment`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "UnComment",
			Router:           `/:vid:int/comment/:cid:int`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "ReplyComment",
			Router:           `/:vid:int/comment/:cid:int`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "UnLike",
			Router:           `/:vid:int/like`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Like",
			Router:           `/:vid:int/like`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "Share",
			Router:           `/:vid:int/share`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"] = append(beego.GlobalControllerRouter["GoliGoli/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "View",
			Router:           `/:vid:int/view`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
