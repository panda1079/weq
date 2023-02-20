package public

type Request struct {
	RequestList map[string]map[string]string
}

// Get 添加get路由
func (r *Request) Get(route string, ac string, ct string) {
	var method = "GET"
	var reKey = method + "__" + route

	if r.RequestList == nil {
		r.RequestList = make(map[string]map[string]string)
	}
	if r.RequestList[reKey] == nil {
		r.RequestList[reKey] = make(map[string]string)
	}

	//["GET"]
	var elm = map[string]string{"method": method, "route": route, "ac": ac, "ct": ct} //定义插入数组
	r.RequestList[reKey] = elm                                                        //把路由插入
}

// Post 添加post路由
func (r *Request) Post(route string, ac string, ct string) {
	var method = "POST"
	var reKey = method + "__" + route

	if r.RequestList == nil {
		r.RequestList = make(map[string]map[string]string)
	}
	if r.RequestList[reKey] == nil {
		r.RequestList[reKey] = make(map[string]string)
	}

	//["POST"]
	var elm = map[string]string{"method": method, "route": route, "ac": ac, "ct": ct} //定义插入数组
	r.RequestList[reKey] = elm                                                        //把路由插入
}
