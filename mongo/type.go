package mongo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Operator string

const (
	Equal         Operator = "$eq"
	GT            Operator = "$gt"
	GTE           Operator = "$gte"
	LT            Operator = "$lt"
	LTE           Operator = "$lte"
	NotEqual      Operator = "$ne"
	IN            Operator = "$in"
	NotIn         Operator = "$nin"
	Or            Operator = "$or"
	And           Operator = "$and"
	Not           Operator = "not"
	Nor           Operator = "nor"
	Exists        Operator = "$exists"
	Type          Operator = "type"
	Coordinates   Operator = "coordinates"
	Mod           Operator = "mod"
	Like          Operator = "$regex"
	Text          Operator = "text"
	Where         Operator = "where"
	GeoWithIn     Operator = "geoWithin"
	GeoIntersects Operator = "geoIntersects"
	Near          Operator = "near"
	NearSphere    Operator = "nearSphere"
	All           Operator = "all"
	ElemMatch     Operator = "elemMatch"
	Size          Operator = "$size"
	Sum           Operator = "$sum"
	AddToSet      Operator = "$addToSet"
	Set           Operator = "$set"
	Inc           Operator = "$inc"
	SetIsSubset   Operator = "$setIsSubset"
	Filter        Operator = "$filter"
	Last          Operator = "$last"
	Cond          Operator = "$cond"

	Match   Operator = "$match"
	Project Operator = "$project"
	Unwind  Operator = "$unwind"
	Group   Operator = "$group"
	Limit   Operator = "$limit"
	Skip    Operator = "$skip"
	Sort    Operator = "$sort"
	Lookup  Operator = "$lookup"
	GeoNear Operator = "$geoNear"

	DateToStringOper Operator = "$dateToString"
	DayOfMonthOper   Operator = "$dayOfMonth"
	Concat           Operator = "$concat"
	Substr           Operator = "$substr"
)

type IModel interface {
	GetId() bson.ObjectId
	GetMgoInfo() (string, string, string) //database,collection,session
}

type DBOperator func(*mgo.Collection) error

type Pagination struct {
	PageIndex    int64       `json:"pageIndex,omitempty"`
	PageSize     int64       `json:"pageSize,omitempty"`
	TotalRecords int64       `json:"totalRecords"`
	Header       interface{} `json:"header,omitempty"`
	Records      interface{} `json:"records"`
}

type PipeLine interface {
	Do(Operator, interface{}, ...bool) PipeLine
	All(interface{}) error
	One(interface{}) error
	Count() (int, error)
	Limit(int) PipeLine
	Skip(int) PipeLine
	Pagination(*Pagination) error
	DocOfPipelines() string
	Project([]string, ...Map)
	Lookup(IModel, string, string, string) PipeLine
}
