package library

import (
	"html"
	"net/http"
	"strings"
)

// HttpInfo 请求内容
type HttpInfo struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// R 获取请求参数  key : 变量名    defaultValue : 默认值
func (r *HttpInfo) R(key string, defaultValue string) string {
	var values = r.Request.URL.Query()
	var res = values.Get(key)

	if len(res) == 0 {
		res = r.Request.PostFormValue(key) //获取post数据
	}

	if len(res) == 0 {
		return defaultValue
	}

	var arg = html.EscapeString(strings.Trim(res, " ")) //转义html字符串
	arg = strings.Replace(arg, "(", "&#40;", -1)        //对小写括号进行处理
	arg = strings.Replace(arg, ")", "&#41;", -1)        //对小写括号进行处理

	return arg
}
