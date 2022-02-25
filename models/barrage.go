package models

import (
	"time"
)

type Barrage struct {
	Id       int
	Content  string
	PostTime time.Time
	PlayTime time.Time
	State    int
	User     *User `orm:"rel(fk)"`
}

func (b *Barrage) Add() error {
	id, err := o.Insert(b)
	if err != nil {
		b.Id = int(id)
	}
	return err
}

func (b *Barrage) Get() error {
	return o.Read(b)
}

func (b *Barrage) Set() error {
	_, err := o.Update(b)
	return err
}

func (b *Barrage) Del() error {
	_, err := o.Delete(b)
	return err
}
