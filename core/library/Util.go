package library

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Util 公共工具

// SetLog 输出日志内容
func SetLog(log any, info string) {
	//当前时间戳
	t1 := time.Now().Unix() //1564552562
	t2 := time.Unix(t1, 0).String()
	fmt.Print("[" + t2 + "][" + info + "][info]：[") //输出头描述
	fmt.Print(log)                                  //输出内容
	fmt.Println("]")                                //直接结尾
}

// GetRequests //获取所有请求内容（用于输出日志,直接返回json字符串）
func GetRequests(CH HttpInfo) string {
	// 对于cli模式的特殊照顾
	if CH.IsCli {
		r := make(map[string]interface{})
		for k, v := range CH.Mount {
			r[k] = v
		}
		return MapToJson(r)
	}

	// 正常的http模式
	if strings.Contains(CH.GetHeader("Content-Type"), "application/json") {
		return CH.Body
	} else {
		r := make(map[string]interface{})
		for k, v := range CH.Form {
			r[k] = v
		}
		return MapToJson(r)
	}
}

// OutJson json输出
func OutJson(CH HttpInfo, OutData map[string]interface{}) {
	jsonBytes, err := json.Marshal(OutData) //转换json
	if err != nil {
		fmt.Println(err)
		fprintf, err := fmt.Fprintf(CH.ResponseWriter, "{\"code\": \"1\", \"route\": \"输出错误\"}") // 发送响应到客户端
		if err != nil {
			SetLog(fprintf, "错误输出")
			return
		}
		return
	}

	//CH.ResponseWriter.WriteHeader(200) //设置响应码

	//对cli模式的特殊照顾
	if CH.IsCli {
		SetLog("["+CH.ClientRealIP()+"]/  :  "+MapToJson(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志
		SetLog(string(jsonBytes), "输出页面")
		return
	}

	SetLog("["+CH.ClientRealIP()+"]"+CH.Request.URL.RequestURI()+"  :  "+MapToJson(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志

	fprintf, err := fmt.Fprintf(CH.ResponseWriter, string(jsonBytes)) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}

// OutHtml 输出http
func OutHtml(CH HttpInfo, html string, OutData map[string]interface{}) {
	data, err := ioutil.ReadFile("./app/template/" + html)
	if err != nil {
		SetLog(err, "错误读取文件")
		return
	}
	html = string(data)

	for k, v := range OutData {
		html = strings.Replace(html, "<{$"+k+"}>", InterfaceToString(v), -1)
	}

	//对cli模式的特殊照顾
	if CH.IsCli {
		SetLog("["+CH.ClientRealIP()+"]/  :  "+MapToJson(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志
		SetLog(html, "输出页面")
		return
	}

	SetLog("["+CH.ClientRealIP()+"]"+CH.Request.URL.RequestURI()+"  :  "+MapToJson(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志

	fprintf, err := fmt.Fprintf(CH.ResponseWriter, html) // 发送响应到客户端
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
func RandStr(length int, isInt bool) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if isInt {
		letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// StrToHex 字符串转16进制
func StrToHex(str string) string {
	return hex.EncodeToString([]byte(str))
}

// HexToStr 16进制转字符串
func HexToStr(hexStringData string) string {
	hexData, _ := hex.DecodeString(hexStringData)
	return string(hexData)
}

// ReverseString 字符串倒叙
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// Enhex 16进制转加密串（二次传输加密用）
func Enhex(hex string) string {
	var ValNum = [4]map[string]string{
		map[string]string{"1": "g", "2": "h", "3": "i", "4": "j", "5": "k", "6": "l", "7": "m", "8": "n", "9": "o", "0": "p"},
		map[string]string{"1": "g", "2": "r", "3": "s", "4": "t", "5": "u", "6": "v", "7": "w", "8": "x", "9": "y", "0": "z"},
		map[string]string{"1": "G", "2": "H", "3": "I", "4": "J", "5": "K", "6": "L", "7": "M", "8": "N", "9": "O", "0": "P"},
		map[string]string{"1": "Q", "2": "R", "3": "S", "4": "T", "5": "U", "6": "V", "7": "W", "8": "X", "9": "Y", "0": "Z"},
	}

	// 字符串倒序
	hex = ReverseString(hex)

	//在字符串前面加一个随机数字
	str := strconv.Itoa(rand.Intn(9) + 1)

	//随机数位转大写
	HexS := []rune(hex)
	for i1 := 0; i1 < len(HexS); i1++ {
		Val := string(HexS[i1])

		if (HexS[i1] >= 97 && HexS[i1] <= 122) || (HexS[i1] >= 65 && HexS[i1] <= 90) {
			if rand.Intn(2) == 1 {
				str = str + Val
			} else {
				str = str + strings.ToUpper(Val)
			}
		} else {
			str = str + ValNum[rand.Intn(4)][Val]
		}

		if rand.Intn(20) < 5 {
			str = str + strconv.Itoa(rand.Intn(8)+1)
		}
	}

	//插入随机字母
	str = str + RandStr(1, false)

	//遵循base64规则，补充4倍位
	for i2 := 0; i2 < (len(str) % 4); i2++ {
		str = str + "="
	}

	return str
}

// Dehex 加密串转16进制（二次传输加密用）
func Dehex(hex string) string {
	var ValNum = map[string]string{
		"g": "1", "h": "2", "i": "3", "j": "4", "k": "5", "l": "6", "m": "7", "n": "8", "o": "9", "p": "0",
		"r": "2", "s": "3", "t": "4", "u": "5", "v": "6", "w": "7", "x": "8", "y": "9", "z": "0",
		"G": "1", "H": "2", "I": "3", "J": "4", "K": "5", "L": "6", "M": "7", "N": "8", "O": "9", "P": "0",
		"Q": "1", "R": "2", "S": "3", "T": "4", "U": "5", "V": "6", "W": "7", "X": "8", "Y": "9", "Z": "0",
	}

	//去除等号
	hex = strings.Replace(hex, "=", "", -1)

	//去除数字(由于首位是数字，就顺带去除了)
	hex = regexp.MustCompile(`\d+`).ReplaceAllString(hex, "")

	//去除最后一个字符
	hex = hex[0 : len(hex)-1]

	// 字符串倒序
	hex = ReverseString(hex)

	//字母转数字
	str := ""
	HexS := []rune(hex)
	for i1 := 0; i1 < len(HexS); i1++ {
		Val := string(HexS[i1])
		if _, ok := ValNum[Val]; ok {
			//如果存在key，就是数字演变的字母，即换回数字
			str = str + ValNum[Val]
		} else {
			str = str + Val
		}
	}

	// 字符转小写
	str = strings.ToLower(str)

	return str
}
