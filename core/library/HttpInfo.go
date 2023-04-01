package library

import (
	"encoding/json"
	"html"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// HttpInfo 请求内容
type HttpInfo struct {
	IsCli          bool
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Body           string
	Mount          map[string]string
	Form           url.Values
	MultipartForm  *multipart.Form
}

// GetHeader 获取请求头信息
func (CH *HttpInfo) GetHeader(key string) string {
	//对于cli的特殊关照
	if CH.IsCli {
		return ""
	}
	return CH.Request.Header.Get(key)
}

// RH 获取请求参数不编码html  key : 变量名    defaultValue : 默认值
func (CH *HttpInfo) RH(key string, defaultValue any) interface{} {
	//对于cli的特殊关照
	if CH.IsCli {
		if value, ok := CH.Mount[key]; ok {
			return value
		}

		//实在没有，就把默认值返回去吧
		return defaultValue
	}

	//获取GET参数
	var values = CH.Request.URL.Query()
	var res = values.Get(key)

	if len(res) == 0 {
		res = CH.Request.PostFormValue(key) //获取post数据
	}

	if len(res) == 0 {
		res = CH.Request.FormValue(key) //获取From数据
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

// R 获取请求参数  key : 变量名    defaultValue : 默认值    Type:类型（string/int）
func (CH *HttpInfo) R(key string, defaultValue any, Type string) interface{} {
	var res = InterfaceToString(CH.RH(key, defaultValue)) //直接获取，省事

	if Type == "string" {
		var arg = html.EscapeString(strings.Trim(res, " ")) //转义html字符串
		arg = strings.Replace(arg, "(", "&#40;", -1)        //对小写括号进行处理
		arg = strings.Replace(arg, ")", "&#41;", -1)        //对小写括号进行处理
		arg = strings.Replace(arg, "=", "&#61;", -1)        //对等于括号进行处理
		return arg
	} else if Type == "int" {
		ins, _ := strconv.Atoi(res)
		return ins
	} else if Type == "int64" {
		ins64, _ := strconv.ParseInt(res, 10, 64)
		return ins64
	}

	return defaultValue
}

// GetReUrl 获得访问路径并去掉get参数
func (CH *HttpInfo) GetReUrl() string {
	var reUrl = strings.Split(CH.Request.URL.RequestURI(), "?")
	return reUrl[0]
}

// ClientRealIP 获取用户的真实IP
func (CH *HttpInfo) ClientRealIP() string {
	ip := "127.0.0.1" //弄个初始值

	//对于cli的特殊关照
	if CH.IsCli {
		addRs, err := net.InterfaceAddrs()
		if err != nil {
			SetLog(err, "获取本机ip失败")
			return ip
		}

		for _, address := range addRs {
			// 检查ip地址判断是否回环地址
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
		return ip
	}

	ipSub, _, _ := net.SplitHostPort(CH.Request.RemoteAddr)
	if net.ParseIP(ipSub) != nil {
		ip = ipSub
	}

	if In(ip[0:3], []string{"127", "172", "192"}) {
		xri := CH.GetHeader("X-Real-IP")
		xff := CH.GetHeader("X-Forward-For")

		if net.ParseIP(xri) != nil {
			ip = xri
		}

		if net.ParseIP(xff) != nil {
			ip = xff
		}
	}
	return ip
}
