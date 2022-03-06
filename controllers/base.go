package controllers

import (
	"GoliGoli/models"
	"GoliGoli/utils"
	"errors"
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"google.golang.org/protobuf/proto"
	"net/http"
)

const EndPoint = "http://localhost:8080"

type (
	Reply struct {
		Code int         `json:"code"`
		Msg  interface{} `json:"msg"`
		Data interface{} `json:"data"`
	}
	BaseController struct {
		beego.Controller
		user models.User
	}
	ErrorMsg struct {
		Error string
		Value interface{}
	}
	ErrorList map[string]ErrorMsg
	Json      map[string]interface{}
)

func (c *BaseController) Reply(reply interface{}) {
	c.Data["json"] = reply
	if err := c.ServeJSON(); err != nil {
		utils.Display("err3", err)
	}
}

func (c *BaseController) ReplyError(err error) {
	var reply Reply
	reply.Code = http.StatusInternalServerError
	reply.Msg = err.Error()
	c.Reply(reply)
}

func (c *BaseController) Bind(obj interface{}) error {
	ct, exist := c.Ctx.Request.Header["Content-Type"]
	if !exist || len(ct) == 0 {
		//return c.Ctx.BindJSON(obj)
		return c.Ctx.BindForm(obj)
	}
	i, l := 0, len(ct[0])
	for i < l && ct[0][i] != ';' {
		i++
	}
	switch ct[0][0:i] {
	case context.ApplicationJSON:
		return c.Ctx.BindJSON(obj)
	case context.ApplicationXML, context.TextXML:
		return c.Ctx.BindXML(obj)
	case context.ApplicationForm, "multipart/form-data":
		return c.Ctx.BindForm(obj)
	case context.ApplicationProto:
		return c.Ctx.BindProtobuf(obj.(proto.Message))
	case context.ApplicationYAML:
		return c.Ctx.BindYAML(obj)
	default:
		return errors.New("Unsupported Content-Type:" + ct[0])
	}
}

func (c *BaseController) ShouldBind(form interface{}) error {
	var reply Reply
	var valid validation.Validation
	if err := c.Bind(form); err != nil {
		c.ReplyError(err)
		return err
	}
	//utils.Display(form)
	if ok, err := valid.Valid(form); err != nil {
		c.ReplyError(err)
		return err
	} else if !ok {
		reply.Code = http.StatusBadRequest
		reply.Msg = "Unexpected params."
		errs := make(ErrorList, len(valid.Errors))
		for _, v := range valid.Errors {
			errs[v.Field] = ErrorMsg{v.Message, v.Value}
		}
		reply.Data = errs
		c.Reply(reply)
		return errors.New("unexpected params")
	}
	return nil
}

func (c *BaseController) Auth(permit ...func(int) bool) error {
	// todo: divide into an identity reader and an identity authorizer
	var reply Reply
	if id := c.GetSession("user"); id == nil {
		reply.Code = http.StatusUnauthorized
		reply.Msg = "User not logged in."
		c.Reply(reply)
		return errors.New("user not logged in")
	} else {
		uid, ok := id.(uint)
		if !ok {
			err := errors.New("userid is not uint")
			c.ReplyError(err)
			utils.Display(err)
			return err
		}
		c.user.ID = uid
	}
	if exist, err := c.user.Check(); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return err
	} else if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user does not exist."
		if err := c.DelSession("user"); err != nil {
			utils.Display(err)
		}
		c.Reply(reply)
		return errors.New("user does not exist")
	}
	if err := c.user.Query(); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return err
	}
	if len(permit) > 0 {
		if !permit[0](c.user.State) {
			reply.Code = http.StatusForbidden
			reply.Msg = "Permission denied."
			c.Reply(reply)
			return errors.New("permission denied")
		}
	}
	return nil
}
