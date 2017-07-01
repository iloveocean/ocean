package mongo

import (
	"encoding/json"
	"errors"

	"gopkg.in/mgo.v2"
)

type QuerySet struct {
	session    *mgo.Session
	collection *mgo.Collection
	iter       *mgo.Iter
	query      *mgo.Query
	cond       *Condition
	sort       []string
	selector   map[string]bool
	limit      int
	d, c, s    string //name of DB, collection, session
}

func (q *QuerySet) All(result interface{}) error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	return q.iter.All(result)
}

func (q *QuerySet) Count() (int, error) {
	if err := q.initQuery(); err != nil {
		return 0, err
	}
	defer q.session.Close()
	return q.query.Count()
}

func (q *QuerySet) Delete() error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	changer := mgo.Change{
		Remove: true,
	}
	_, err := q.query.Apply(changer, Map{})
	return err
}

func (q *QuerySet) DeleteAll() error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	_, err := q.collection.RemoveAll(q.cond.conditions)
	return err
}

func (q *QuerySet) Distinct(field string, result interface{}) error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	err := q.query.Distinct(field, result)
	return err
}

func (q *QuerySet) Exist() (bool, error) {
	var (
		num int
		err error
	)
	if num, err = q.Count(); err != nil {
		return false, err
	}
	return (num > 0), nil
}

func (q *QuerySet) Fields(fields []string) *QuerySet {
	if len(fields) > 0 {
		for _, item := range fields {
			q.selector[item] = true
		}
	}
	return q
}

func (q *QuerySet) Filter(field string, value interface{}, opt Operator) *QuerySet {
	q.cond.addCond(field, value, opt)
	return q
}

//for debug
func (q *QuerySet) GetQueryDoc() string {
	docRaw, err := json.MarshalIndent(q.cond.conditions, "", " ")
	if err != nil {
		return errors.New("json marshal error: " + err.Error()).Error()
	}
	return string(docRaw)
}

func (q *QuerySet) Limit(n int) *QuerySet {
	q.limit = n
	return q
}

func (q *QuerySet) One(result interface{}) error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	return q.query.One(result)
}

func (q *QuerySet) Pagination(page *Pagination) error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()

	//stuff total number
	if total, err := q.query.Count(); err != nil {
		return err
	} else {
		page.TotalRecords = int64(total)
	}

	//stuff records
	if page.Records == nil {
		page.Records = Map{}
	}
	jump := (page.PageIndex - 1) * page.PageSize
	return q.query.Skip(int(jump)).Limit(int(page.PageSize)).All(page.Records)
}

func (q *QuerySet) SetCond(new *Condition) *QuerySet {
	q.cond.mergeCond(new)
	return q
}

func (q *QuerySet) Sort(sortFields map[string]int, keys ...string) *QuerySet {
	sort := []string{}
	if len(keys) > 0 {
		for _, item := range keys {
			if v, ok := sortFields[item]; ok {
				if v > 0 {
					sort = append(sort, item)
				} else {
					sort = append(sort, "-"+item)
				}
			}
		}
	} else {
		for k, v := range sortFields {
			if v > 0 {
				sort = append(sort, k)
			} else {
				sort = append(sort, "-"+k)
			}
		}
	}
	q.sort = sort
	return q
}

func (q *QuerySet) Update(cond Map, upsert bool) error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	ch := mgo.Change{
		Update:    Map{"$set": cond},
		ReturnNew: false,
		Upsert:    upsert,
	}
	_, err := q.query.Apply(ch, Map{})
	return err
}

func (q *QuerySet) UpdateWithCommand(command string, cond Map, upsert bool) error {
	if err := q.initQuery(); err != nil {
		return err
	}
	defer q.session.Close()
	ch := mgo.Change{
		Update:    Map{command: cond},
		ReturnNew: false,
		Upsert:    upsert,
	}
	_, err := q.query.Apply(ch, Map{})
	return err
}

//func (q *QuerySet) UpdateAll

func (q *QuerySet) initQuery() error {
	var err error
	if q.session, err = CopySession(q.s); err != nil {
		return err
	}
	q.collection = q.session.DB(q.d).C(q.c)
	q.query = q.collection.Find(q.cond.conditions)
	q.iter = q.query.Iter()

	if len(q.sort) > 0 {
		q.query.Sort(q.sort...)
	}

	if len(q.selector) > 0 {
		q.query.Select(q.selector)
	}

	if q.limit > 0 {
		q.query.Limit(q.limit)
	}
	return nil
}
