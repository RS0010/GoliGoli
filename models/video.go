package models

type Video struct {
	Id           int
	Title        string
	LikeCount    int
	CommentCount int
	CollectCount int
	ShareCount   int
	ViewCount    int
	Likes        []*User `orm:"reverse(many)"`
	Collects     []*User `orm:"reverse(many)"`
	Shares       []*User `orm:"reverse(many)"`
	Authors      []*User `orm:"reverse(many)"`
}

func (v *Video) Add() error {
	id, err := o.Insert(v)
	if err != nil {
		v.Id = int(id)
	}
	return err
}

func (v *Video) Get() error {
	return o.Read(v)
}

func (v *Video) Set() error {
	_, err := o.Update(v)
	return err
}

func (v *Video) Del() error {
	_, err := o.Delete(v)
	return err
}
