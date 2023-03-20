package library

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
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

// In 判断字符串是否在数组种
func In(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

// MakeRequest 发起http请求
//url    访问路径
//params 参数，该数组多于1个，表示为POST
//extend 请求伪造包头参数
//返回的为一个请求状态，一个内容
func MakeRequest(url string, params map[string]interface{}, extend map[string]string) map[string]interface{} {
	if url == "" {
		return map[string]interface{}{"code": "100"}
	}

	//参数数组多于1个，表示为POST
	met := "GET"
	if len(params) > 1 {
		met = "POST"
	}

	paramStr, _ := json.Marshal(params) //转换格式

	client := &http.Client{}
	req, err := http.NewRequest(met, url, bytes.NewReader(paramStr))
	if err != nil {
		SetLog(err, "发起http请求错误1")
	}

	//写入包头
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	for k1, v1 := range extend {
		req.Header.Set(k1, v1)
	}

	resp, err := client.Do(req)
	var res = map[string]interface{}{
		"Code":   resp.StatusCode,
		"Header": resp.Header,
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SetLog(err, "发起http请求错误2")
	}

	res["result"] = string(body)
	return res
}

// RandStr 产生随机字符串
func RandStr(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
