package mongo

import (
	"encoding/json"
	"errors"

	"gopkg.in/mgo.v2"
)

type pipeline struct {
	pipelines  []Map
	d, c, s    string //db, collection and session name
	session    *mgo.Session
	collection *mgo.Collection
	pipe       *mgo.Pipe
}

func (p *pipeline) Do(op Operator, cond interface{}, selectors ...bool) PipeLine {
	if len(selectors) > 0 {
		for _, selector := range selectors {
			if !selector {
				return p
			}
		}
	}
	m := M(string(op), cond)
	p.pipelines = append(p.pipelines, m)
	return p
}

func (p *pipeline) All(result interface{}) error {
	if err := p.initPipe(); err != nil {
		return err
	}
	defer p.session.Close()
	return p.pipe.All(result)
}

func (p *pipeline) One(result interface{}) error {
	if err := p.initPipe(); err != nil {
		return err
	}
	defer p.session.Close()
	return p.pipe.One(result)
}

func (p *pipeline) Count() (int, error) {
	countCond := M("_id", nil).M("count", 1, Sum)
	p.Do(Group, countCond)
	if err := p.initPipe(); err != nil {
		return 0, err
	}
	defer p.session.Close()
	var counter struct {
		Count int `bson:count`
	}
	if err := p.pipe.One(&counter); err != nil {
		return 0, err
	}
	return counter.Count, nil
}

func (p *pipeline) Skip(num int) PipeLine {
	p.pipelines = append(p.pipelines, M(string(Skip), num))
	return p
}

func (p *pipeline) Limit(num int) PipeLine {
	p.pipelines = append(p.pipelines, M(string(Limit), num))
	return p
}

func (p *pipeline) Pagination(page *Pagination) error {
	if page == nil {
		return errors.New("input page parameter is null!")
	}
	if page.Records == nil {
		//var models []*interface{}
		page.Records = &[]*interface{}{}
	}

	//get total number
	if total, err := p.getCount(); err != nil {
		return err
	} else {
		page.TotalRecords = int64(total)
	}

	//get page
	if err := p.initPipe(); err != nil {
		return err
	}
	start := (page.PageIndex - 1) * page.PageSize
	return p.Skip(int(start)).Limit(int(page.PageSize)).All(page.Records)
}

func (p *pipeline) DocOfPipelines() string {
	docRaw, err := json.MarshalIndent(p.pipelines, "", " ")
	if err != nil {
		return errors.New("json marshal error: " + err.Error()).Error()
	}
	return string(docRaw)
}

func (p *pipeline) Project(fields []string, addConds ...Map) {
	projectMap := Map{}
	for _, field := range fields {
		projectMap.M(field, 1)
	}
	for _, value := range addConds {
		for k, v := range value {
			projectMap[k] = v
		}
	}
	p.Do(Project, projectMap)
}

func (p *pipeline) Lookup(join IModel, localField, foreignField, asField string) PipeLine {
	if join == nil || localField == "" || foreignField == "" || asField == "" {
		return p
	}
	//get join collection name
	_, collection, _ := join.GetMgoInfo()

	p.Do(Lookup, Map{
		"from":         collection,
		"localField":   localField,
		"foreignField": foreignField,
		"as":           asField,
	})
	return p
}

func (p *pipeline) initPipe() error {
	session, err := CopySession(p.s)
	if err != nil {
		return err
	}
	p.session = session
	p.collection = session.DB(p.d).C(p.c)
	p.pipe = p.collection.Pipe(p.pipelines)
	return nil
}

func (p *pipeline) getCount() (int, error) {
	tmp := append(p.pipelines, M("_id", nil).M("count", 1, Sum))
	session, err := CopySession(p.s)
	if err != nil {
		return 0, err
	}
	defer session.Close()
	var result struct {
		Count int `bson:count`
	}
	pipe := session.DB(p.d).C(p.c).Pipe(tmp)
	if err := pipe.One(&result); err != nil {
		return 0, err
	} else {
		return result.Count, nil
	}
}
