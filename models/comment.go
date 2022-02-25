package models

import (
	"time"
)

type Comment struct {
	Id      int
	Content string
	Time    time.Time
	Parent  int
	State   int
	User    *User `orm:"rel(fk)"`
}

func (c *Comment) Add() error {
	id, err := o.Insert(c)
	if err != nil {
		c.Id = int(id)
	}
	return err
}

func (c *Comment) Get() error {
	return o.Read(c)
}

func (c *Comment) Set() error {
	_, err := o.Update(c)
	return err
}

func (c *Comment) Del() error {
	_, err := o.Delete(c)
	return err
}
