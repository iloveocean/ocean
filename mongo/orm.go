package mongo

type orm struct{}

func NewOrm() *orm {
	return new(orm)
}

func (o *orm) Pipe(model IModel) PipeLine {
	d, c, s := model.GetMgoInfo()
	pl := &pipeline{
		d: d,
		c: c,
		s: s,
	}
	return pl
}
