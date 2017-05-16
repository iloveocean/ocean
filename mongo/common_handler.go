package mongo

import (
	oceanTest "ocean/testing"
	"testing"
)

var preError error = nil

func openDBHandler(next oceanTest.TestHandler) oceanTest.TestHandler {
	//open mongodb must be the first step of test process,
	//so initialize preError == nil here
	preError = nil

	fn := func(t *testing.T) {
		preError = StartUp("master", "127.0.0.1:27017", "OceanTest", "", "")
		if preError != nil {
			t.Error("err: 连接 MongoDB 失败...")
			return
		}
		defer Shutdown()
		next.RunTest(t)
	}
	return oceanTest.TestHandlerFunc(fn)
}
