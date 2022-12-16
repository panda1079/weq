package core

type Request struct {
	RequestList []map[string]string
}

// Get 添加get路由
func (r *Request) Get(route string, ac string, ct string) {
	var method :=[1]

//["GET"]
	var elm = map[string]string{"method": method, "route": route, "ac": ac, "ct": ct} //定义插入数组
	r.RequestList = append(r.RequestList, elm)                                        //把路由插入
}

// Post 添加post路由
func (r *Request) Post(route string, ac string, ct string) {
	var elm = map[string]string{"route": route, "ac": ac, "ct": ct} //定义插入数组
	r.RequestList = append(r.RequestList, elm)                      //把路由插入
}
