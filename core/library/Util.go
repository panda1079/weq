package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Util 公共工具

// SetLog 输出日志内容
func SetLog(a any, info string) {
	//当前时间戳
	t1 := time.Now().Unix() //1564552562
	t2 := time.Unix(t1, 0).String()

	fmt.Print("[" + t2 + "][" + info + "][info]：[") //输出头描述
	fmt.Print(a)                                    //输出内容
	fmt.Print("]\n")                                //完毕，换行
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
func OutHtml(w http.ResponseWriter, html string, OutData map[string]interface{}) {
	data, err := ioutil.ReadFile("./app/template/" + html)
	if err != nil {
		SetLog(err, "错误读取文件")
		return
	}
	html = string(data)

	for k, v := range OutData {
		html = strings.Replace(html, "<{$"+k+"}>", InterfaceToString(v), -1)
	}

	fprintf, err := fmt.Fprintf(w, html) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}

// MapToJson map转json
func MapToJson(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// JsonToMap json转map
func JsonToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}

// InterfaceToString interface转String
func InterfaceToString(inter interface{}) string {
	var res = ""
	switch inter.(type) {
	case string:
		res = inter.(string)
		break
	case int:
		res = strconv.Itoa(inter.(int))
		break
	case float32:
		res = fmt.Sprintf("%g", inter)
		break
	case float64:
		res = fmt.Sprintf("%g", inter)
		break
	}

	return res
}
