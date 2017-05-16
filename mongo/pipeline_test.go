package mongo

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"loyocloud-oa/models/crm/order"
	oceanTest "ocean/testing"
)

type order4Test struct {
	base order.OrderInfo
}

func (o *order4Test) GetMgoInfo() (string, string, string) { //database,collection,session
	return "OceanTest", "order", "master"
}

func (o *order4Test) GetId() bson.ObjectId {
	return o.base.GetId()
}

func TestOne(t *testing.T) {
	commonHandlers := oceanTest.New(openDBHandler)
	h := commonHandlers.ThenFunc(testOneHandler)
	h.RunTest(t)
}

func testOneHandler(t *testing.T) {
	if preError != nil {
		t.Errorf("Error occurred in the previous step")
		return
	}

	//get pipe
	o := NewOrm()
	pipe := o.Pipe(new(order4Test))

	//construct query condition(s)
	cond := M("company_id", bson.ObjectIdHex("56a09438e44c36717f7b0d86"))

	//take aciton here
	result := order.OrderInfo{}
	if err := pipe.Do(Match, cond).One(&result); err != nil {
		t.Errorf("err: query db meet error: %s", err.Error())
		return
	} else {
		if result.Id.Hex() != "57aecb87e44c36181c000100" {
			t.Errorf("record Id = %s, want %s", result.Id.Hex(), "57aecb87e44c36181c000100")
		}
	}
}

func TestAll(t *testing.T) {
	commonHandlers := oceanTest.New(openDBHandler)
	h := commonHandlers.ThenFunc(testAllHandler)
	h.RunTest(t)
}

func testAllHandler(t *testing.T) {
	if preError != nil {
		t.Errorf("Error occurred in the previous step")
		return
	}
	//get pipe
	o := NewOrm()
	pipe := o.Pipe(new(order4Test))

	//construct query condition(s)
	cond := M(string(Or), []interface{}{
		M("company_id", bson.ObjectIdHex("56a09438e44c36717f7b0d86")),
		M("company_id", bson.ObjectIdHex("57bd645ad33c6538a797a538")),
	})

	//take action(s)
	results := []*order.OrderInfo{}
	if err := pipe.Do(Match, cond).All(&results); err != nil {
		t.Errorf("err: query db meet error: %s", err.Error())
		return
	} else {
		if len(results) != 26 {
			t.Errorf("err result: %d, want %d", len(results), 4)
		}
	}
}

func TestCount(t *testing.T) {
	commonHandlers := oceanTest.New(openDBHandler)
	h := commonHandlers.ThenFunc(testCountHandler)
	h.RunTest(t)
}

func testCountHandler(t *testing.T) {
	if preError != nil {
		t.Errorf("Error occurred in the previous step")
		return
	}

	//get pipe
	o := NewOrm()
	pipe := o.Pipe(new(order4Test))

	//construct query conditions(s)
	cond := M("company_id", bson.ObjectIdHex("57bd645ad33c6538a797a538"))

	//take action here
	if num, err := pipe.Do(Match, cond).Count(); err != nil {
		t.Errorf("err: query record count meet error: %s", err.Error())
	} else {
		if num != 3 {
			t.Errorf("num = %d, want %d", num, 3)
		}
	}
}

func TestSkip(t *testing.T) {
	commonHandlers := oceanTest.New(openDBHandler)
	h := commonHandlers.ThenFunc(testSkipHandler)
	h.RunTest(t)
}

func testSkipHandler(t *testing.T) {
	if preError != nil {
		t.Errorf("Error occurred in the previous step")
		return
	}

	//get pipe
	o := NewOrm()
	pipe := o.Pipe(new(order4Test))

	//construct query conditions(s)
	cond := M("company_id", bson.ObjectIdHex("56a09438e44c36717f7b0d86"))

	//take action here
	results := []*order.OrderInfo{}
	var (
		num int
		err error
	)
	if num, err = pipe.Do(Match, cond).Count(); err != nil {
		t.Error("query db records number meet error: ", err.Error())
		return
	}
	if num <= 10 {
		t.Error("not enough records for testing, the number must > 10")
		return
	}
	//refresh pipe
	pipe = o.Pipe(new(order4Test))
	if err := pipe.Do(Match, cond).Skip(10).All(&results); err != nil {
		t.Error("query db meet error: ", err.Error())
		return
	} else {
		if len(results) != (num - 10) {
			t.Errorf("got %d results, expected %d", len(results), num-10)
			return
		}
	}
}

func testLimitHandler(t *testing.T) {
	if preError != nil {
		t.Errorf("Error occurred in the previous step")
		return
	}

	//get pipe
	o := NewOrm()
	pipe := o.Pipe(new(order4Test))

	//construct query conditions(s)
	cond := M("company_id", bson.ObjectIdHex("56a09438e44c36717f7b0d86"))

	//take action here
	results := []*order.OrderInfo{}
	var (
		num int
		err error
	)
	if num, err = pipe.Do(Match, cond).Count(); err != nil {
		t.Error("query db records number meet error: ", err.Error())
		return
	}
	if num <= 10 {
		t.Error("not enough records for testing, the number must > 10")
		return
	}
	//refresh pipe
	pipe = o.Pipe(new(order4Test))
	if err := pipe.Do(Match, cond).Limit(10).All(&results); err != nil {
		t.Error("query db meet error: ", err.Error())
		return
	} else {
		if len(results) != 10 {
			t.Errorf("got %d results, expected %d", len(results), 10)
			return
		}
	}
}
