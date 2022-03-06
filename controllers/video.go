package controllers

// todo: vid can't be zero

import (
	"GoliGoli/models"
	"GoliGoli/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	VideoReplyBrief struct {
		ID       uint
		Path     string
		Title    string
		Info     string
		Category string
	}

	VideoReplyDetail struct {
		VideoReplyBrief
		ViewCount      int
		LikerCount     int
		CollectorCount int
		SharerCount    int
		AuthorList     models.UserList
		UserID         uint
	}
)

type VideoController struct {
	BaseController
}

// @router / [get]
func (c *VideoController) Search() {
	var reply Reply
	title := c.GetString("title", "")
	category := c.GetString("category", "")
	userid, _ := c.GetInt("userid", 0)
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 20)
	orderby := c.GetString("orderby", "id")
	ordertype := c.GetString("ordertype", "asc")
	if userid < 0 {
		userid = 0
	}
	filter := models.VideoFilter{
		Title:     title,
		Category:  category,
		PostID:    uint(userid),
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
	reply.Data = []VideoReplyBrief{}
	for _, video := range list {
		v := video.Render(EndPoint)
		reply.Data = append(reply.Data.([]VideoReplyBrief), VideoReplyBrief{
			ID:       v.ID,
			Path:     v.Path,
			Title:    v.Title,
			Info:     v.Info,
			Category: v.Category,
		})
	}
	c.Reply(reply)
}

// @router /:vid:int [get]
func (c *VideoController) Query() {
	var reply Reply
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	err = video.QueryAuthor()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	v := video.Render(EndPoint)
	reply.Code = http.StatusOK
	reply.Msg = "Query successfully."
	reply.Data = VideoReplyDetail{
		VideoReplyBrief: VideoReplyBrief{
			ID:       v.ID,
			Path:     v.Path,
			Title:    v.Title,
			Info:     v.Info,
			Category: v.Category,
		},
		ViewCount:      v.ViewCount,
		LikerCount:     video.CountLiker(),
		CollectorCount: video.CountCollector(),
		SharerCount:    video.CountSharer(),
		UserID:         v.UserID,
		AuthorList:     video.AuthorList,
	}
	c.Reply(reply)
}

// @router / [post]
func (c *VideoController) Upload() {
	// todo: multiuser upload
	// todo: anonymity upload
	var reply Reply
	var form struct {
		Title    string
		Info     string
		Category string
	}
	err := c.Auth()
	if err != nil {
		return
	}
	err = c.ShouldBind(&form)
	if err != nil {
		return
	}
	_, f, err := c.GetFile("Video")
	if err != nil {
		reply.Code = http.StatusForbidden
		reply.Msg = "Upload failed."
		c.Reply(reply)
		return
	}
	var filename, title string
	if form.Title == "" {
		title = f.Filename
		filename = f.Filename
	} else {
		title = form.Title
		filename = form.Title + "." + f.Filename[strings.LastIndex(f.Filename, ".")+1:]
	}
	err = c.SaveToFile("Video", "./data/videos/"+filename)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	video := models.Video{
		Title:    title,
		Info:     form.Info,
		Path:     "/data/videos/" + filename,
		Category: form.Category,
	}
	err = video.Post(&c.user)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Upload successfully."
	v := video.Render(EndPoint)
	reply.Data = VideoReplyBrief{
		ID:       v.ID,
		Path:     v.Path,
		Title:    v.Title,
		Info:     v.Info,
		Category: v.Category,
	}
	c.Reply(reply)
}

// @router /:vid:int [put]
func (c *VideoController) Update() {
	var reply Reply
	var form struct {
		Title    string
		Info     string
		Category string
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = c.Auth()
	if err != nil {
		return
	}
	err = video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if c.user.ID != video.UserID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You are not the author of this video."
		c.Reply(reply)
		return
	}
	err = c.ShouldBind(&form)
	if err != nil {
		return
	}
	if form.Title != "" {
		video.Title = form.Title
	}
	video.Info = form.Info
	video.Category = form.Category
	err = video.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	v := video.Render(EndPoint)
	reply.Code = http.StatusOK
	reply.Msg = "Update successfully."
	reply.Data = VideoReplyBrief{
		ID:       v.ID,
		Path:     v.Path,
		Title:    v.Title,
		Info:     v.Info,
		Category: v.Category,
	}
	c.Reply(reply)
}

// @router /:vid:int [patch]
func (c *VideoController) PartUpdate() {
	var reply Reply
	var form struct {
		Title    string
		Info     string
		Category string
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = c.Auth()
	if err != nil {
		return
	}
	err = video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if c.user.ID != video.UserID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You are not the author of this video."
		c.Reply(reply)
		return
	}
	err = c.ShouldBind(&form)
	if err != nil {
		return
	}
	if form.Title != "" {
		video.Title = form.Title
	}
	if form.Info != "" {
		video.Info = form.Info
	}
	if form.Category != "" {
		video.Category = form.Category
	}
	err = video.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	v := video.Render(EndPoint)
	reply.Code = http.StatusOK
	reply.Msg = "Update successfully."
	reply.Data = VideoReplyBrief{
		ID:       v.ID,
		Path:     v.Path,
		Title:    v.Title,
		Info:     v.Info,
		Category: v.Category,
	}
	c.Reply(reply)
}

// @router /:vid:int [delete]
func (c *VideoController) Delete() {
	var reply Reply
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = c.Auth()
	if err != nil {
		return
	}
	err = video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if c.user.ID != video.UserID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You are not the author of this video."
		c.Reply(reply)
		return
	}
	err = video.Delete()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Delete successfully."
	reply.Data = Json{"id": vid}
	c.Reply(reply)
}

// @router /:vid:int/like [post]
func (c *VideoController) Like() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	like, err := c.user.IsLike(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if like {
		reply.Code = http.StatusForbidden
		reply.Msg = "You have liked this video."
		c.Reply(reply)
		return
	}
	err = c.user.Like(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Like successfully."
	reply.Data = Json{"id": video.ID}
	c.Reply(reply)
}

// @router /:vid:int/like [delete]
func (c *VideoController) UnLike() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	like, err := c.user.IsLike(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !like {
		reply.Code = http.StatusForbidden
		reply.Msg = "You have not liked this video."
		c.Reply(reply)
		return
	}
	err = c.user.Unlike(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Unlike successfully."
	reply.Data = Json{"id": video.ID}
	c.Reply(reply)
}

// @router /:vid:int/collect [post]
func (c *VideoController) Collect() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	collect, err := c.user.IsCollect(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if collect {
		reply.Code = http.StatusForbidden
		reply.Msg = "You have collected this video."
		c.Reply(reply)
		return
	}
	err = c.user.Collect(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Collect successfully."
	reply.Data = Json{"id": video.ID}
	c.Reply(reply)
}

// @router /:vid:int/collect [delete]
func (c *VideoController) UnCollect() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	collect, err := c.user.IsCollect(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !collect {
		reply.Code = http.StatusForbidden
		reply.Msg = "You have not collected this video."
		c.Reply(reply)
		return
	}
	err = c.user.Collect(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Uncollect successfully."
	reply.Data = Json{"id": video.ID}
	c.Reply(reply)
}

// @router /:vid:int/share [post]
func (c *VideoController) Share() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	share, err := c.user.IsShare(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if share {
		reply.Code = http.StatusForbidden
		reply.Msg = "You have shared this video."
		c.Reply(reply)
		return
	}
	err = c.user.Share(&video)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Share successfully."
	reply.Data = Json{"id": video.ID}
	c.Reply(reply)
}

// @router /:vid:int/view [post]
func (c *VideoController) View() {
	var reply Reply
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	err := video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	video.ViewCount++
	err = video.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "View successfully."
	c.Reply(reply)
}

// @router /:vid:int/comment [post]
func (c *VideoController) Comment() {
	var reply Reply
	var form struct {
		Content string `valid:"Required"`
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = c.Auth()
	if err != nil {
		return
	}
	err = c.ShouldBind(&form)
	if err != nil {
		return
	}
	comment := models.Comment{
		Content: form.Content,
		UserID:  c.user.ID,
		VideoID: video.ID,
	}
	err = c.user.Comment(&video, &comment)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Comment successfully."
	reply.Data = comment
	c.Reply(reply)
}

// @router /:vid:int/comment [get]
func (c *VideoController) QueryComment() {
	var reply Reply
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	uid, _ := c.GetInt("uid", 0)
	cid, _ := c.GetInt("cid", 0)
	orderby := c.GetString("orderby", "id")
	ordertype := c.GetString("ordertype", "asc")
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 20)
	filter := models.CommentFilter{
		ID:        uint(cid),
		VideoID:   uint(vid),
		UserID:    uint(uid),
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
	reply.Data = list
	c.Reply(reply)
}

// @router /:vid:int/comment/:cid:int [delete]
func (c *VideoController) UnComment() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	cid, _ := strconv.Atoi(c.Ctx.Input.Param(":cid"))
	comment := models.Comment{ID: uint(cid)}
	exist, err = comment.Check()
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Comment does not exist."
		c.Reply(reply)
		return
	}
	err = comment.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if comment.UserID != c.user.ID || video.UserID != c.user.ID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You can not delete this comment."
		c.Reply(reply)
		return
	}
	err = comment.Delete()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Delete successfully."
	reply.Data = comment
	c.Reply(reply)
}

// @router /:vid:int/comment/:cid:int [post]
func (c *VideoController) ReplyComment() {
	var reply Reply
	var form struct {
		Content string `valid:"Required"`
	}
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	cid, _ := strconv.Atoi(c.Ctx.Input.Param(":cid"))
	comment := models.Comment{ID: uint(cid)}
	exist, err = comment.Check()
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Comment does not exist."
		c.Reply(reply)
		return
	}
	err = c.ShouldBind(&form)
	if err != nil {
		return
	}
	replyComment := models.Comment{
		Content: form.Content,
		UserID:  c.user.ID,
		VideoID: video.ID,
		Parent:  cid,
	}
	err = c.user.Comment(&video, &replyComment)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Reply successfully."
	reply.Data = replyComment
	c.Reply(reply)
}

// @router /:vid:int/barrage [post]
func (c *VideoController) Barrage() {
	var reply Reply
	var form struct {
		Content string    `valid:"Required"`
		Time    time.Time `valid:"Required"`
	}
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = c.ShouldBind(&form)
	if err != nil {
		return
	}
	barrage := models.Barrage{
		Content:  form.Content,
		PlayTime: form.Time,
		UserID:   c.user.ID,
		VideoID:  video.ID,
	}
	err = c.user.Barrage(&video, &barrage)
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Barrage successfully."
	reply.Data = barrage
	c.Reply(reply)
}

// @router /:vid:int/barrage/:bid:int [delete]
func (c *VideoController) UnBarrage() {
	var reply Reply
	err := c.Auth()
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	bid, _ := strconv.Atoi(c.Ctx.Input.Param(":bid"))
	barrage := models.Barrage{ID: uint(bid)}
	exist, err = barrage.Check()
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Barrage does not exist."
		c.Reply(reply)
		return
	}
	err = barrage.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if barrage.UserID != c.user.ID || barrage.VideoID != video.ID {
		reply.Code = http.StatusForbidden
		reply.Msg = "You can not delete this barrage."
		c.Reply(reply)
		return
	}
	err = barrage.Delete()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Delete successfully."
	reply.Data = barrage
	c.Reply(reply)
}

// @router /:vid:int/barrage [get]
func (c *VideoController) QueryBarrage() {
	var reply Reply
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	bid, _ := c.GetInt("bid", 0)
	uid, _ := c.GetInt("uid", 0)
	orderby := c.GetString("orderby", "id")
	ordertype := c.GetString("ordertype", "asc")
	page, _ := c.GetInt("page", 1)
	pagesize, _ := c.GetInt("pagesize", 20)
	filter := models.BarrageFilter{
		ID:        uint(bid),
		VideoID:   uint(vid),
		UserID:    uint(uid),
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
	reply.Data = list
	c.Reply(reply)
}

// @router /:vid:int/ban [post]
func (c *VideoController) Ban() {
	var reply Reply
	err := c.Auth(func(state int) bool {
		return state > 10
	})
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if video.State < 0 {
		reply.Code = http.StatusForbidden
		reply.Msg = "The video has been banned."
		c.Reply(reply)
		return
	}
	video.State = -1
	err = video.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Ban successfully."
	v := video.Render(EndPoint)
	reply.Data = VideoReplyBrief{
		ID:       uint(vid),
		Path:     v.Path,
		Title:    v.Title,
		Info:     v.Info,
		Category: v.Category,
	}
	c.Reply(reply)
}

// @router /:vid:int/ban [delete]
func (c *VideoController) Unban() {
	var reply Reply
	err := c.Auth(func(state int) bool {
		return state > 10
	})
	if err != nil {
		return
	}
	vid, _ := strconv.Atoi(c.Ctx.Input.Param(":vid"))
	video := models.Video{ID: uint(vid)}
	exist, err := video.Check()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if !exist {
		reply.Code = http.StatusForbidden
		reply.Msg = "Video does not exist."
		c.Reply(reply)
		return
	}
	err = video.Query()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	if video.State >= 0 {
		reply.Code = http.StatusForbidden
		reply.Msg = "The user has not been banned."
		c.Reply(reply)
		return
	}
	video.State = 0
	err = video.Update()
	if err != nil {
		c.ReplyError(err)
		utils.Display(err)
		return
	}
	reply.Code = http.StatusOK
	reply.Msg = "Ban successfully."
	v := video.Render(EndPoint)
	reply.Data = VideoReplyBrief{
		ID:       uint(vid),
		Path:     v.Path,
		Title:    v.Title,
		Info:     v.Info,
		Category: v.Category,
	}
	c.Reply(reply)
}
