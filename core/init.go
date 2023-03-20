package core

import (
	"bytes"
	"config"
	"controller"
	"core/library"
	"net/http"
	"reflect"
	"routes"
)

var RegisterMessage = make(map[string]interface{})
var routeList = routes.Web() //加载路由

// LoadRoute 加载控制器函数
func LoadRoute(w http.ResponseWriter, r *http.Request) {

	//初始化http处理结构体，把http信息压入结构体内
	var HttpInfo = library.HttpInfo{}
	HttpInfo.ResponseWriter = w
	HttpInfo.Request = r

	//预制body内容raw访问使用
	var buf = new(bytes.Buffer)
	from, err := buf.ReadFrom(r.Body)
	if err != nil {
		library.SetLog(from, "错误输出")
		library.SetLog(err, "错误输出")
		library.OutJson(w, map[string]string{"code": "0", "msg": "预制body失败"})
		return
	}
	HttpInfo.Body = buf.String()

	//获得访问路径并去掉get参数
	var route = HttpInfo.GetReUrl()

	//获取当前url的路由设置 map[ac:order_list ct:CtlOrder method:GET route:/order/order_list]
	lr := library.Request{}
	Mount, RInfo := lr.GetRInfo(r, routeList, route)
	HttpInfo.Mount = Mount

	//判断是否存在路由
	if RInfo != nil {
		//循环控制器列表
		for k1, v1 := range RegisterMessage {
			//fmt.Print(k1, RInfo["ct"])
			//fmt.Print("\n")
			if k1 == RInfo["ct"] { //找到控制器
				//预创建控制器对象
				var methodArgs []reflect.Value
				methodArgs = append(methodArgs, reflect.ValueOf(HttpInfo))

				//把包含http内容的结构体推给控制器
				var CtlBox = reflect.ValueOf(v1).MethodByName(RInfo["ac"])
				CtlBox.Call(methodArgs)

				//完事了就直接退出
				return
			}
		}

	} else {
		library.OutJson(w, map[string]string{"code": "0", "msg": "路由不存在"})
	}
}

func Init() {

	//加载预设服务模块
	var SS = library.ServerS{}
	SS.InitServerS()

	//初始化控制器池
	var ctl = controller.CtlIndex{}
	RegisterMessage = ctl.Init(SS)

	//获取启动配置
	deploy := config.Deploy{}
	con := deploy.Run()

	http.HandleFunc("/", LoadRoute)
	err := http.ListenAndServe(con["LISTEN_ADDRESS"]+":"+con["PORT"], nil)
	if err != nil {
		library.SetLog(err, "错误输出")
		return
	}
}
