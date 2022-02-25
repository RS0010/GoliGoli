package models

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id       int        `orm:"auto"`
	Username string     `orm:"unique"`
	Password string     `orm:"size(255)"`
	Portrait string     `orm:"null"`
	Gender   string     `orm:"null"`
	Age      int        `orm:"null"`
	Address  string     `orm:"null"`
	Email    string     `orm:"null"`
	State    int        `orm:"default(0)"`
	Likes    []*Video   `orm:"rel(m2m);rel_table(user_video_likes)"`
	Collects []*Video   `orm:"rel(m2m);rel_table(user_video_collects)"`
	Shares   []*Video   `orm:"rel(m2m);rel_table(user_video_shares)"`
	Blocks   []*User    `orm:"rel(m2m);rel_through(GoliGoli/models.UserUserBlocks)"`
	Blocked  []*User    `orm:"reverse(many)"`
	videos   []*Video   `orm:"reverse(many)"`
	Comments []*Comment `orm:"reverse(many)"`
	Barrages []*Barrage `orm:"reverse(many)"`
}

type UserUserBlocks struct {
	Id      int
	Blocker *User `orm:"rel(fk)"`
	Blocked *User `orm:"rel(fk)"`
}

func (u *User) Add() error {
	id, err := o.Insert(u)
	if err != nil {
		u.Id = int(id)
	}
	return err
}

func (u *User) Get() error {
	return o.Read(u)
}

func (u *User) Set() error {
	_, err := o.Update(u)
	return err
}

func (u *User) Del() error {
	_, err := o.Delete(u)
	return err
}

func (u *User) Block(b *User) error {
	m2m := o.QueryM2M(u, "Blocks")
	m2m.Add(b)
	_, err := o.Insert(&UserUserBlocks{Blocker: u, Blocked: b})
	if err != nil {
		return err
	}
	_, err = o.LoadRelated(u, "Blocks")
	_, err = o.LoadRelated(u, "Blocked")
	_, err = o.LoadRelated(b, "Blocks")
	_, err = o.LoadRelated(b, "Blocked")
	fmt.Println(err)
	return err
}

func (u *User) Like(v *Video) error {
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		v.LikeCount++
		_, err := txOrm.Update(v)
		if err != nil {
			return err
		}
		m2m := txOrm.QueryM2M(u, "Likes")
		_, err = m2m.Add(v)
		if err != nil {
			return err
		}
		_, err = txOrm.LoadRelated(u, "Likes")
		return err
	})
	if err != nil {
		v.LikeCount--
	}
	return err
}

func (u *User) Unlike(v *Video) error {
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		v.LikeCount--
		_, err := txOrm.Update(v)
		if err != nil {
			return err
		}
		m2m := txOrm.QueryM2M(u, "Likes")
		_, err = m2m.Remove(v)
		if err != nil {
			return err
		}
		_, err = txOrm.LoadRelated(u, "Likes")
		return err
	})
	if err != nil {
		v.LikeCount++
	}
	return err
}

//func (u *User) Collect(v *Video) error {
//
//}
