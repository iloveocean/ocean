package mongo

import "fmt"

type orm struct{}

func NewOrm() *orm {
	return new(orm)
}

func (o *orm) Pipe(model IModel) PipeLine {
	d, c, s := model.GetMgoInfo()
	fmt.Println("d, c, s is :", d, c, s)
	pl := &pipeline{
		d: d,
		c: c,
		s: s,
	}
	return pl
}
