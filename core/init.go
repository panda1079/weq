package core

import (
	"controller"
	"core/library"
	"fmt"
	"net/http"
	"reflect"
	"routes"
	"strings"
)

var routeList = routes.Web()
var RegisterMessage = make(map[string]interface{})

// LoadRoute 加载控制器函数
func LoadRoute(w http.ResponseWriter, r *http.Request) {
	//获得访问路径并去掉get参数
	var reUrl = strings.Split(r.URL.RequestURI(), "?")
	var route = reUrl[0]

	//获取当前url的路由设置 map[ac:order_list ct:CtlOrder method:GET route:/order/order_list]
	var RInfo = routeList[r.Method+"__"+route]
	//fmt.Print(r.Method + "__" + route)

	//判断是否存在路由
	if RInfo != nil {
		for k1, v1 := range RegisterMessage {
			if k1 == RInfo["ct"] { //找到控制器

				//预创建控制器对象
				var methodArgs []reflect.Value
				var CtlBox = reflect.ValueOf(v1).MethodByName(RInfo["ac"])

				//把http信息压入结构体内
				var HttpInfo = library.HttpInfo{}
				HttpInfo.ResponseWriter = w
				HttpInfo.Request = r
				methodArgs = append(methodArgs, reflect.ValueOf(HttpInfo))

				//把包含http内容的结构体推给控制器
				CtlBox.Call(methodArgs)

				return
			}
		}

	} else {
		fprintf, err := fmt.Fprintf(w, "{\"code\": \"1\", \"route\": \"路由不存在\"}") // 发送错误响应到客户端
		if err != nil {
			library.SetLog(fprintf)
			return
		}
	}
}

func Init() {
	//初始化控制器池
	var ctl = controller.CtlIndex{}
	RegisterMessage = ctl.Init()

	//拉起http-web服务
	http.HandleFunc("/", LoadRoute)
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		library.SetLog(err)
		return
	}
}
