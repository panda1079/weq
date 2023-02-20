package core

import (
	"fmt"
	"net/http"
	"routes"
	"strings"
	//"controller"
)

var routeList = routes.Web()

func LoadRoute(w http.ResponseWriter, r *http.Request) {

	//获得访问路径并去掉get参数
	var reUrl = strings.Split(r.URL.RequestURI(), "?")
	var route = reUrl[0]

	//判断是否存在路由
	if routeList[r.Method+"__"+route] != nil {
		fmt.Print(routeList[r.Method+"__"+route])
		fmt.Print("\n")

		//var ctl = controller.

		//fmt.Fprintf(w, "111111") // 发送响应到客户端

	} else {
		fmt.Fprintf(w, "22222222") // 发送响应到客户端
	}

}

func Init() {
	http.HandleFunc("/", LoadRoute)
	http.ListenAndServe(":9091", nil)
}
