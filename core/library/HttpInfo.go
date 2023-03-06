package library

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
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
	var ct = CH.Request.Header["Content-Type"]

	if len(ct) > 0 {
		return ct[0]
	}

	return ""
}

// RH 获取请求参数不编码html  key : 变量名    defaultValue : 默认值
func (CH *HttpInfo) RH(key string, defaultValue string) string {
	var values = CH.Request.URL.Query()
	var res = values.Get(key)

	if len(res) == 0 {
		res = CH.Request.PostFormValue(key) //获取post数据
	}

	if len(res) == 0 {
		//获取 raw参数
		m := make(map[string]interface{})
		Str := json.Unmarshal([]byte(CH.Body), &m)
		if Str == nil {
			switch m[key].(type) {
			case string:
				res = m[key].(string)
				break
			case int:
				res = strconv.Itoa(m[key].(int))
				break
			case float32:
				res = fmt.Sprintf("%g", m[key])
				break
			case float64:
				res = fmt.Sprintf("%g", m[key])
				break
			}
		}
	}

	if len(res) == 0 {
		if value, ok := CH.Mount[key]; ok {
			res = value
		}
	}

	if len(res) == 0 {
		return defaultValue
	}

	return res
}

// R 获取请求参数  key : 变量名    defaultValue : 默认值
func (CH *HttpInfo) R(key string, defaultValue string) string {
	var res = CH.RH(key, defaultValue)

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
