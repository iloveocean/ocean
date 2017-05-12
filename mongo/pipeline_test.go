package mongo

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"loyocloud-oa/models/crm/order"
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

func connectDB() error {
	err := StartUp("master", "127.0.0.1:27017", "OceanTest", "", "")
	if err != nil {
		fmt.Println("err: 连接 MongoDB 失败...")
		return err
	} else {
		fmt.Println("连接 MongoDB 成功...")
		return nil
	}
}

func TestOne(t *testing.T) {
	//connect to DB
	if err := connectDB(); err != nil {
		t.Error("failed to connect test Mongo DB")
		return
	}
	defer Shutdown()

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
