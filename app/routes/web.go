package routes

import (
	"core/library"
)

func Web() map[string]map[string]string {
	var AddRe = library.Request{} //core.Request.Post()

	// ====================== 下面开始写路由 ====================== //

	AddRe.Get("/test/test", "Test", "CtlTest")
	AddRe.Post("/test/test", "Test", "CtlTest")
	AddRe.AddRoute([2]string{"Post", "Get"}, "/test", "TestA", "CtlTest")
	AddRe.AddRoute([2]string{"Post", "Get"}, "/(?P<ct>[a-z|A-Z|\\d+]+)/(?P<ac>[a-z|A-Z|\\d+]+)", "{ac}", "{ct}")

	// ====================== 上面开始写路由 ====================== //

	return AddRe.RequestList
}
