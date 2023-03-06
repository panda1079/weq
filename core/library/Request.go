package library

import (
	"net/http"
	"regexp"
	"strings"
)

// Request 关于路由的公共函数
type Request struct {
	RequestList map[string]map[string]string
}

// AddRe 添加路由内容到数组 "GET__/order/order_list":map[string]string{"ac":"order_list","ct":"CtlOrder","method":"GET","route":"/order/order_list"}
func (r *Request) AddRe(method string, route string, ac string, ct string) {
	var reKey = method + "__" + route

	if r.RequestList == nil {
		r.RequestList = make(map[string]map[string]string)
	}
	if r.RequestList[reKey] == nil {
		r.RequestList[reKey] = make(map[string]string)
	}

	//Add ["GET"] // ["POST"]
	var elm = map[string]string{"method": method, "route": route, "ac": ac, "ct": ct} //定义插入数组
	r.RequestList[reKey] = elm                                                        //把路由插入
}

// Get 添加get路由
func (r *Request) Get(route string, ac string, ct string) {
	r.AddRe("GET", route, ac, ct) //集中插入
}

// Post 添加post路由
func (r *Request) Post(route string, ac string, ct string) {
	r.AddRe("POST", route, ac, ct) //集中插入                                                     //把路由插入
}

// AddRoute 支持多种请求方式
func (r *Request) AddRoute(methods [2]string, route string, ac string, ct string) {
	for _, value1 := range methods {
		if value1 == "Get" {
			r.Get(route, ac, ct)
		}

		if value1 == "Post" {
			r.Post(route, ac, ct)
		}
	}
}

// GetRInfo 获取适合的路由
func (r *Request) GetRInfo(Rr *http.Request, routeList map[string]map[string]string, route string) (map[string]string, map[string]string) {
	var run = make(map[string]string)
	var Mount = make(map[string]string)

	// 1.先找出非正则的，再找出正则的
	if value, ok := routeList[Rr.Method+"__"+route]; ok {
		run = value
	} else {
		// 2.正则的需要有内容才算
		for _, Value1 := range routeList {

			var re = regexp.MustCompile(Value1["route"])
			var match = re.FindStringSubmatch(route)

			if len(match) > 0 {
				var groupNames = re.SubexpNames()

				run = Value1
				for i, name := range groupNames {
					if i != 0 && name != "" {
						//SetLog(name + "-----" + match[i])

						run["ac"] = strings.Replace(run["ac"], "{"+name+"}", match[i], -1)
						run["ct"] = strings.Replace(run["ct"], "{"+name+"}", match[i], -1)

						Mount[name] = match[i] //向请求参数内添加额外参数（请勿与请求参数起冲突，否则将会替换掉请求参数）
					}
				}
				break

			}
		}
	}

	return Mount, run
}
