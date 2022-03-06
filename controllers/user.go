package controllers

import (
	"GoliGoli/models"
	"GoliGoli/utils"
	"GoliGoli/utils/bcrypt"
	"net/http"
	"strconv"
	"time"
)

// test
type UserController struct {
	BaseController
}

//func (c *UserController) Prepare() {
//	c.EnableXSRF = true
//}

// @Title User Register
// @Description User Register
// @Param Username body string true "username"
// @Param Password body string true "password"
// @Success 200 {object} controllers.Reply
// @Success 400 {object} controllers.Reply
// @router /register [post]
func (c *UserController) Register() {
	var reply Reply
	var form struct {
		Username string `valid:"Required"`
		Password string `valid:"Required"`
	}
	// Get the form data and validate
	if err := c.ShouldBind(&form); err != nil {
		return
	}
	// Check whether the username already exists
	c.user.Username = form.Username
	if exist, err := c.user.Check(); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	} else if exist {
		reply.Code = http.StatusConflict
		reply.Msg = "The username already exists."
		reply.Data = map[string]string{"Username": form.Username}
		c.Reply(reply)
		return
	}
	// Create user
	if hash, err := bcrypt.Encrypt(form.Password); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	} else {
		c.user.Password = hash
		if err := c.user.Create(); err != nil {
			c.ReplyError(err)
			utils.Display(err)
			return
		}
		if err := c.SetSession("user", c.user.ID); err != nil {
			utils.Display(err)
		}
		reply.Code = http.StatusOK
		reply.Msg = "Register successfully."
		reply.Data = map[string]interface{}{"UserID": c.user.ID, "Username": form.Username}
		c.Reply(reply)
		return
	}
}

// @router /login [post]
func (c *UserController) Login() {
	var reply Reply
	var form struct {
		Username string `form:"username" valid:"Required"`
		Password string `form:"password" valid:"Required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		return
	}
	c.user.Username = form.Username
	if exist, err := c.user.Check(); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	} else if !exist {
		reply.Code = http.StatusNotFound
		reply.Msg = "The username does not exist."
		reply.Data = map[string]string{"Username": form.Username}
		c.Reply(reply)
		return
	}
	if err := c.user.Query(); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !bcrypt.Verify(c.user.Password, form.Password) {
		reply.Code = http.StatusUnauthorized
		reply.Msg = "The password is incorrect."
		reply.Data = map[string]interface{}{"UserID": c.user.ID, "Username": form.Username}
		c.Reply(reply)
		return
	} else {
		reply.Code = http.StatusOK
		reply.Msg = "Login successfully."
		reply.Data = map[string]interface{}{"UserID": c.user.ID, "Username": form.Username}
		if err := c.SetSession("user", c.user.ID); err != nil {
			utils.Display(err)
		}
		c.Reply(reply)
		return
	}
}

// @router /login [delete]
func (c *UserController) Logout() {
	var reply Reply
	if err := c.Auth(); err != nil {
		return
	}
	if err := c.DelSession("user"); err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Logout successfully."
	reply.Data = Json{"UserID": c.user.ID, "Username": c.user.Username}
	c.Reply(reply)
}

// @router / [get]
func (c *UserController) Query() {
	var reply Reply
	search, _ := c.GetBool("search", false)
	if search {
		uid, _ := c.GetInt("uid", 0)
		name := c.GetString("name", "")
		email := c.GetString("email", "")
		orderby := c.GetString("orderby", "id")
		ordertype := c.GetString("ordertype", "asc")
		page, _ := c.GetInt("page", 1)
		pagesize, _ := c.GetInt("pagesize", 20)
		if uid < 0 {
			uid = 0
		}
		filter := models.UserFilter{
			ID:        uint(uid),
			Name:      name,
			Email:     email,
			Page:      page,
			PageSize:  pagesize,
			OrderBy:   orderby,
			OrderType: ordertype,
		}
		list, err := filter.Search()
		if err != nil {
			c.ReplyError(err)
			utils.Display(err)
			return
		}
		reply.Code = http.StatusOK
		reply.Msg = "Query successfully."
		reply.Data = []models.User{}
		for _, user := range list {
			reply.Data = append(reply.Data.([]models.User), user.Render(EndPoint))
		}
		c.Reply(reply)
	} else {
		if err := c.Auth(); err != nil {
			return
		}
		reply.Code = http.StatusOK
		reply.Msg = "Query successfully."
		reply.Data = c.user.Render(EndPoint)
		c.Reply(reply)
	}
}

// @router / [put]
func (c *UserController) Update() {
	var reply Reply
	var form struct {
		Username    string
		Gender      string
		Age         int
		Address     string
		Email       string
		Password    string
		NewPassword string
	}
	reply.Msg = []string{}
	if err := c.Auth(); err != nil {
		return
	}
	if err := c.ShouldBind(&form); err != nil {
		return
	}
	if form.NewPassword != "" {
		if bcrypt.Verify(c.user.Password, form.Password) {
			c.user.Password = form.NewPassword
			reply.Msg = append(reply.Msg.([]string), "Password changed.")
		} else {
			reply.Msg = append(reply.Msg.([]string), "The password is incorrect.")
		}
	}
	if form.Username != "" && form.Username != c.user.Username {
		user := models.User{Username: form.Username}
		exist, err := user.Check()
		if err != nil {
			c.ReplyError(err)
			utils.Display(err)
			return
		}
		if !exist {
			c.user.Username = form.Username
		} else {
			reply.Msg = append(reply.Msg.([]string), "The username already exists.")
		}
	}
	err := c.SaveToFile("Portrait", "./data/portrait/"+strconv.Itoa(int(c.user.ID))+".jpg")
	if err == nil {
		c.user.Portrait = "/data/portrait/" + strconv.Itoa(int(c.user.ID)) + ".jpg"
		reply.Msg = append(reply.Msg.([]string), "Portrait changed.")
	}
	c.user.Gender = form.Gender
	c.user.Age = form.Age
	c.user.Address = form.Address
	c.user.Email = form.Email
	err = c.user.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if len(reply.Msg.(string)) == 0 {
		reply.Msg = "Update successfully."
	} else {
		reply.Msg = append(reply.Msg.([]string), "Update successfully.")
	}
	reply.Data = c.user.Render(EndPoint)
	c.Reply(reply)
}

// @router / [patch]
func (c *UserController) PartUpdate() {
	var reply Reply
	var form struct {
		Username    string
		Portrait    string
		Gender      string
		Age         int
		Address     string
		Email       string
		Password    string
		NewPassword string
	}
	reply.Msg = []string{}
	if err := c.Auth(); err != nil {
		return
	}
	if err := c.ShouldBind(&form); err != nil {
		return
	}
	if form.NewPassword != "" {
		if bcrypt.Verify(c.user.Password, form.Password) {
			c.user.Password = form.NewPassword
			reply.Msg = append(reply.Msg.([]string), "Password changed.")
		} else {
			reply.Msg = append(reply.Msg.([]string), "The password is incorrect.")
		}
	}
	if form.Username != "" && form.Username != c.user.Username {
		user := models.User{Username: form.Username}
		exist, err := user.Check()
		if err != nil {
			c.ReplyError(err)
			utils.Display(err)
			return
		}
		if !exist {
			c.user.Username = form.Username
		} else {
			reply.Msg = append(reply.Msg.([]string), "The username already exists.")
		}
	}
	err := c.SaveToFile("Portrait", "./data/portrait/"+strconv.Itoa(int(c.user.ID))+".jpg")
	if err == nil {
		c.user.Portrait = "/data/portrait/" + strconv.Itoa(int(c.user.ID)) + ".jpg"
		reply.Msg = append(reply.Msg.([]string), "Portrait changed.")
	}
	if form.Gender != "" {
		c.user.Gender = form.Gender
	}
	if form.Age != 0 {
		c.user.Age = form.Age
	}
	if form.Address != "" {
		c.user.Address = form.Address
	}
	if form.Email != "" {
		c.user.Email = form.Email
	}
	err = c.user.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if len(reply.Msg.(string)) == 0 {
		reply.Msg = "Update successfully."
	} else {
		reply.Msg = append(reply.Msg.([]string), "Update successfully.")
	}
	reply.Data = c.user.Render(EndPoint)
	c.Reply(reply)
}

// @router / [delete]
func (c *UserController) Delete() {
	var reply Reply
	var form struct {
		Password string `valid:"Required"`
	}
	if err := c.Auth(); err != nil {
		return
	}
	if err := c.ShouldBind(&form); err != nil {
		return
	}
	if bcrypt.Verify(c.user.Password, form.Password) {
		c.user.Username += "-" + strconv.FormatInt(time.Now().UnixMicro(), 10)
		err := c.user.Update()
		if err != nil {
			c.ReplyError(err)
			utils.Display(err)
			return
		}
		err = c.user.Delete()
		if err != nil {
			c.ReplyError(err)
			utils.Display(err)
			return
		}
		err = c.DelSession("user")
		if err != nil {
			utils.Display(err)
		}
		reply.Code = http.StatusOK
		reply.Msg = "Delete successfully."
	} else {
		reply.Code = http.StatusForbidden
		reply.Msg = "The password is incorrect."
	}
	reply.Data = Json{"UserID": c.user.ID, "Username": c.user.Username}
	c.Reply(reply)
}

// @router /:uid:int/ban [post]
func (c *UserController) Ban() {
	var reply Reply
	uid, _ := strconv.Atoi(c.Ctx.Input.Param(":uid"))
	if err := c.Auth(func(state int) bool {
		return state > 10
	}); err != nil {
		return
	}
	if uint(uid) == c.user.ID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You can't ban yourself."
		c.Reply(reply)
		return
	}
	user := models.User{ID: uint(uid)}
	exist, err := user.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user does not exist."
		c.Reply(reply)
		return
	}
	err = user.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if user.State < 0 {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user has been banned."
		c.Reply(reply)
		return
	}
	if user.State >= c.user.State {
		reply.Code = http.StatusForbidden
		reply.Msg = "You can not ban the user."
		c.Reply(reply)
		return
	}
	user.State = -1 - user.State
	err = user.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Ban successfully."
	reply.Data = Json{"UserID": user.ID, "Username": user.Username}
	c.Reply(reply)
}

// @router /:uid:int/ban [delete]
func (c *UserController) Unban() {
	var reply Reply
	uid, _ := strconv.Atoi(c.Ctx.Input.Param(":uid"))
	if err := c.Auth(func(state int) bool {
		return state > 10
	}); err != nil {
		return
	}
	user := models.User{ID: uint(uid)}
	exist, err := user.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user does not exist."
		c.Reply(reply)
		return
	}
	err = user.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if user.State >= 0 {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user has not been banned."
		c.Reply(reply)
		return
	}
	if -1-user.State >= c.user.State {
		reply.Code = http.StatusForbidden
		reply.Msg = "You can not unban the user."
		c.Reply(reply)
		return
	}
	user.State = -1 - user.State
	err = user.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Unban successfully."
	reply.Data = Json{"UserID": user.ID, "Username": user.Username}
	c.Reply(reply)
}

// @router /:uid:int/block [post]
func (c *UserController) Block() {
	var reply Reply
	uid, _ := strconv.Atoi(c.Ctx.Input.Param(":uid"))
	if err := c.Auth(); err != nil {
		return
	}
	if uint(uid) == c.user.ID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You can't block yourself."
		c.Reply(reply)
		return
	}
	user := models.User{ID: uint(uid)}
	exist, err := user.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user does not exist."
		c.Reply(reply)
		return
	}
	block, err := c.user.IsBlock(&user)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if block {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user has been blocked."
		c.Reply(reply)
		return
	}
	err = c.user.Block(&user)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Block successfully."
	reply.Data = Json{"UserID": user.ID, "Username": user.Username}
	c.Reply(reply)
}

// @router /:uid:int/block [delete]
func (c *UserController) Unblock() {
	var reply Reply
	uid, _ := strconv.Atoi(c.Ctx.Input.Param(":uid"))
	if err := c.Auth(); err != nil {
		return
	}
	user := models.User{ID: uint(uid)}
	exist, err := user.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user does not exist."
		c.Reply(reply)
		return
	}
	block, err := c.user.IsBlock(&user)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if block {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user has not been blocked."
		c.Reply(reply)
		return
	}
	err = c.user.Unblock(&user)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Unblock successfully."
	reply.Data = Json{"UserID": user.ID, "Username": user.Username}
	c.Reply(reply)
}
