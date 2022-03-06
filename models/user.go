package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type User struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      `gorm:"" json:"-"`
	UpdatedAt   time.Time      `gorm:"" json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Username    string         `gorm:"unique"`
	Password    string         `gorm:"size:255" json:"-"`
	Portrait    string         `gorm:""`
	Gender      string         `gorm:""`
	Age         int            `gorm:""`
	Address     string         `gorm:""`
	Email       string         `gorm:""`
	State       int            `gorm:"default:0"`
	PostList    VideoList      `gorm:"many2many:user_videos" json:"-"`
	LikeList    VideoList      `gorm:"many2many:user_likes" json:"-"`
	CollectList VideoList      `gorm:"many2many:user_collects" json:"-"`
	ShareList   VideoList      `gorm:"many2many:user_shares" json:"-"`
	BlockList   UserList       `gorm:"many2many:user_blocks" json:"-"`
	CommentList []*Comment     `gorm:"" json:"-"`
	BarrageList []*Barrage     `gorm:"" json:"-"`
	VideoList   VideoList      `gorm:"" json:"-"`
}

func (u User) Render(endpoint string) User {
	if u.Portrait != "" {
		u.Portrait = endpoint + u.Portrait
	}
	return u
}

func (u *User) Check() (exist bool, err error) {
	var count int64
	if u.ID == 0 {
		err = o.Model(&User{}).Where("username = ?", u.Username).Count(&count).Error
	} else {
		err = o.Model(&User{}).Where("id = ?", u.ID).Count(&count).Error
	}
	return count > 0, err
}

func (u *User) Create() error {
	return o.Model(u).Create(u).Error
}

func (u *User) Query() error {
	if u.ID == 0 {
		return o.Model(u).Where("username = ?", u.Username).First(u).Error
	} else {
		return o.Model(u).First(u, u.ID).Error
	}
}

func (u *User) Update() error {
	return o.Save(u).Error
}

func (u *User) Delete() error {
	return o.Select(clause.Associations).Delete(u).Error
}

func (u *User) Block(b *User) error {
	return o.Model(u).Association("BlockList").Append(b)
}

func (u *User) Unblock(b *User) error {
	return o.Model(u).Association("BlockList").Delete(b)
}

func (u *User) IsBlock(b *User) (bool, error) {
	var blockList []*User
	err := o.Model(u).Association("BlockList").Find(&blockList, b.ID)
	return len(blockList) > 0, err
}

func (u *User) QueryBlock() error {
	return o.Model(u).Association("BlockList").Find(&u.BlockList)
}

func (u *User) Like(v *Video) error {
	return o.Model(u).Association("LikeList").Append(v)
}

func (u *User) Unlike(v *Video) error {
	return o.Model(u).Association("LikeList").Delete(v)
}

func (u *User) IsLike(v *Video) (bool, error) {
	var likeList []*Video
	err := o.Model(u).Association("LikeList").Find(&likeList, v.ID)
	return len(likeList) > 0, err
}

func (u *User) QueryLike() error {
	return o.Model(u).Association("LikeList").Find(&u.LikeList)
}

func (u *User) Collect(v *Video) error {
	return o.Model(u).Association("CollectList").Append(v)
}

func (u *User) Uncollect(v *Video) error {
	return o.Model(u).Association("CollectList").Delete(v)
}

func (u *User) IsCollect(v *Video) (bool, error) {
	var collectList []*Video
	err := o.Model(u).Association("CollectList").Find(&collectList, v.ID)
	return len(collectList) > 0, err
}

func (u *User) QueryCollect() error {
	return o.Model(u).Association("CollectList").Find(&u.CollectList)
}

func (u *User) Share(v *Video) error {
	return o.Model(u).Association("ShareList").Append(v)
}

func (u *User) Unshare(v *Video) error {
	return o.Model(u).Association("ShareList").Delete(v)
}

func (u *User) IsShare(v *Video) (bool, error) {
	var shareList []*Video
	err := o.Model(u).Association("ShareList").Find(&shareList, v.ID)
	return len(shareList) > 0, err
}

func (u *User) QueryShare() error {
	return o.Model(u).Association("ShareList").Find(&u.ShareList)
}

func (u *User) CountPost(filter ...VideoFilter) (total int64, err error) {
	// fixme: when just filtered by user_id, will return count of the videos whose id is satisfy the filter
	// 		  I don't known how to fix it now, it will be fixed when I get a solve for the problem
	if len(filter) > 0 {
		err := o.Where("ID", u.ID).Preload("PostList", func(db *gorm.DB) *gorm.DB {
			m := filter[0].filter(db)
			m.Model(&Video{}).Count(&total)
			return m
		}).Find(&User{}).Error
		if filter[0].filtered {
			return total, err
		}
	}
	return o.Model(u).Association("PostList").Count(), nil
}

func (u *User) QueryPost(filter ...VideoFilter) error {
	if len(filter) > 0 {
		return o.Where("ID", u.ID).Preload("PostList", func(db *gorm.DB) *gorm.DB {
			return filter[0].paginate(filter[0].filter(db))
		}).Find(u).Error
	} else {
		return o.Preload("PostList").Find(u).Error
	}
}

func (u *User) QueryVideo() error {
	return o.Model(u).Association("VideoList").Find(&u.VideoList)
}

func (u *User) Comment(v *Video, c *Comment) error {
	c.UserID = u.ID
	c.VideoID = v.ID
	return o.Create(c).Error
}

func (u *User) QueryComment() error {
	return o.Model(u).Association("CommentList").Find(&u.CommentList)
}

func (u *User) Barrage(v *Video, b *Barrage) error {
	b.UserID = u.ID
	b.VideoID = v.ID
	return o.Create(b).Error
}

func (u *User) QueryBarrage() error {
	return o.Model(u).Association("BarrageList").Find(&u.BarrageList)
}
