package routes

import (
	"core/library"
)

func Web() map[string]map[string]string {
	var AddRe = library.Request{} //core.Request.Post()

	//AddRe.AddRoute([2]string{"Post", "Get"}, "/", "TestA", "CtlTest") //如需使用空路由的请使用这个，个人不建议使用空路由

	// ====================== 下面开始写路由 ====================== //

	AddRe.Get("/test/test", "Test", "CtlTest")
	AddRe.Post("/test/test", "Test", "CtlTest")
	AddRe.AddRoute([2]string{"Post", "Get"}, "/test", "TestA", "CtlTest")

	AddRe.WS("/aaa", "TestB", "CtlTest") //websocket类型

	AddRe.AddRoute([2]string{"Post", "Get"}, "/(?P<ct>[a-z|A-Z|\\d+]+)/(?P<ac>[a-z|A-Z|\\d+]+)", "{ac}", "{ct}")

	// ====================== 上面开始写路由 ====================== //

	return AddRe.RequestList
}
