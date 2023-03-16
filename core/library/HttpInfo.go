package library

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"
)

// HttpInfo 请求内容
type HttpInfo struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Body           string
	Mount          map[string]string
}

// GetHeader 获取请求头信息
func (CH *HttpInfo) GetHeader(key string) string {
	return CH.Request.Header.Get(key)
}

// RH 获取请求参数不编码html  key : 变量名    defaultValue : 默认值
func (CH *HttpInfo) RH(key string, defaultValue string) string {
	//获取GET参数
	var values = CH.Request.URL.Query()
	var res = values.Get(key)

	if len(res) == 0 {
		res = CH.Request.PostFormValue(key) //获取post数据
	}

	if len(res) == 0 {
		//application/json的情形下获取raw参数
		if strings.Contains(CH.GetHeader("Content-Type"), "application/json") {
			m := make(map[string]interface{})
			Str := json.Unmarshal([]byte(CH.Body), &m)
			if Str == nil {
				res = InterfaceToString(m[key]) //转类型
			}
		}
	}

	//获取额外挂载参数
	if len(res) == 0 {
		if value, ok := CH.Mount[key]; ok {
			res = value
		}
	}

	//实在没有，就把默认值返回去吧
	if len(res) == 0 {
		return defaultValue
	}

	return res
}

// R 获取请求参数  key : 变量名    defaultValue : 默认值
func (CH *HttpInfo) R(key string, defaultValue string) interface{} {
	var res = CH.RH(key, defaultValue) //直接获取，省事

	var arg = html.EscapeString(strings.Trim(res, " ")) //转义html字符串
	arg = strings.Replace(arg, "(", "&#40;", -1)        //对小写括号进行处理
	arg = strings.Replace(arg, ")", "&#41;", -1)        //对小写括号进行处理
	arg = strings.Replace(arg, "=", "&#61;", -1)        //对等于括号进行处理
	return arg
}

// GetReUrl 获得访问路径并去掉get参数
func (CH *HttpInfo) GetReUrl() string {
	var reUrl = strings.Split(CH.Request.URL.RequestURI(), "?")
	return reUrl[0]
}
