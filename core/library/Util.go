package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Util 公共工具

// SetLog 输出日志内容
func SetLog(a any, info string) {
	now := time.Now()                  //获取当前时间
	timestamp := now.Unix()            //时间戳
	timeObj := time.Unix(timestamp, 0) //将时间戳转为时间格式

	fmt.Print("[")                       //起~
	fmt.Print(timeObj)                   //输出时间
	fmt.Print("][" + info + "][info]：[") //输出头描述
	fmt.Print(a)                         //输出内容
	fmt.Print("]\n")                     //完毕，换行
}

// OutJson json输出
func OutJson(w http.ResponseWriter, OutData any) {
	jsonBytes, err := json.Marshal(OutData) //转换json
	if err != nil {
		fmt.Println(err)
		fprintf, err := fmt.Fprintf(w, "{\"code\": \"1\", \"route\": \"输出错误\"}") // 发送响应到客户端
		if err != nil {
			SetLog(fprintf, "错误输出")
			return
		}
		return
	}

	fprintf, err := fmt.Fprintf(w, string(jsonBytes)) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}

// OutHtml 输出http
func OutHtml(w http.ResponseWriter, html string, OutData map[string]string) {
	data, err := ioutil.ReadFile("./app/template/" + html)
	if err != nil {
		SetLog(err, "错误读取文件")
		return
	}
	html = string(data)

	for k, v := range OutData {
		html = strings.Replace(html, "<{$"+k+"}>", v, -1)
	}

	fprintf, err := fmt.Fprintf(w, html) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}
