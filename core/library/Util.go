package library

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Util 公共工具

// SetLog 输出日志内容
func SetLog(logStr any, info string) {
	//当前时间戳
	t1 := time.Now().Unix() //1564552562
	t2 := time.Unix(t1, 0).String()
	str := fmt.Sprint("["+t2+"]["+info+"][info]：[", logStr, "]")

	//输出文件
	fmt.Println(str)

}

// GetRequests //获取所有请求内容（用于输出日志,直接返回json字符串）
func GetRequests(CH HttpInfo) string {
	// 对于cli模式的特殊照顾
	if CH.IsCli {
		r := make(map[string]interface{})
		for k, v := range CH.Mount {
			r[k] = v
		}
		return JsonEncode(r)
	}

	// 正常的http模式
	if strings.Contains(CH.GetHeader("Content-Type"), "application/json") {
		return CH.Body
	} else {
		r := make(map[string]interface{})
		for k, v := range CH.Form {
			vv := InterfaceToString(v[0])
			//对于加密内容的特殊处理
			if CH.TransmissionMod == 2 {
				vv = CH.Encryption.HexToStr(CH.Encryption.Dehex(InterfaceToString(vv)))
			}
			r[k] = vv
		}
		return JsonEncode(r)
	}
}

// OutStr 文本输出
func OutStr(CH HttpInfo, Str interface{}) {

	CH.ResponseWriter.WriteHeader(200) //设置响应码

	//对cli模式的特殊照顾
	if CH.IsCli {
		SetLog("["+CH.ClientRealIP()+"]/  :  "+InterfaceToString(Str)+"  :  "+GetRequests(CH), "正常推出") // 写个日志
		SetLog(InterfaceToString(Str), "输出页面")
		return
	}

	SetLog("["+CH.ClientRealIP()+"]"+CH.Request.URL.RequestURI()+"  :  "+InterfaceToString(Str)+"  :  "+GetRequests(CH), "正常推出") // 写个日志

	fprintf, err := fmt.Fprintf(CH.ResponseWriter, InterfaceToString(Str)) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}

// OutJson json输出
func OutJson(CH HttpInfo, OutData map[string]interface{}) {
	OutStr := ""

	if InterfaceToString(CH.TransmissionMod) == "2" {
		Encryption := Encryption{}
		OutMap := map[string]interface{}{}
		//循环获取需要的参数
		for k, v := range OutData {
			OutMap[k] = Encryption.Enhex(Encryption.StrToHex(JsonEncode(v)))
		}
		OutStr = JsonEncode(OutMap) //转换json
	} else {
		OutStr = JsonEncode(OutData) //转换json
	}

	CH.ResponseWriter.WriteHeader(200) //设置响应码

	//对cli模式的特殊照顾
	if CH.IsCli {
		SetLog("["+CH.ClientRealIP()+"]/  :  "+JsonEncode(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志
		SetLog(OutStr, "输出页面")
		return
	}

	SetLog("["+CH.ClientRealIP()+"]"+CH.Request.URL.RequestURI()+"  :  "+JsonEncode(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志

	fprintf, err := fmt.Fprintf(CH.ResponseWriter, OutStr) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}

// OutHtml 输出http
func OutHtml(CH HttpInfo, html string, OutData map[string]interface{}) {
	data, err := os.ReadFile("./app/template/" + html)
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
		SetLog("["+CH.ClientRealIP()+"]/  :  "+JsonEncode(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志
		SetLog(html, "输出页面")
		return
	}

	SetLog("["+CH.ClientRealIP()+"]"+CH.Request.URL.RequestURI()+"  :  "+JsonEncode(OutData)+"  :  "+GetRequests(CH), "正常推出") // 写个日志

	fprintf, err := fmt.Fprintf(CH.ResponseWriter, html) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf, "错误输出")
		return
	}
}

// Time 获取时间戳
func Time() string {
	return InterfaceToString(time.Now().Unix())
}

// GetMillisecond 获取毫秒级时间戳
func GetMillisecond() string {
	return InterfaceToString(time.Now().UnixNano() / int64(time.Millisecond))
}

// StrToTime 获取指定时间时间戳 格式为（年-月-日 时:分:秒  2023-08-10 12:00:00）
func StrToTime(specifiedTime string) string {
	// 将字符串转换为 time.Time 类型
	t, err := time.Parse("2006-01-02 15:04:05", specifiedTime)
	if err != nil {
		SetLog(err, "错误输出")
		return "0"
	}
	// 获取时间戳
	return InterfaceToString(t.Unix() - 28800) //由于系统时间会直接加到东8的八小时，所以要减掉
}

// Date 时间戳转日期(通PHP->date)
//
//	Y ：年（四位数）大写
//	m : 月（两位数，首位不足补0） 小写
//	d ：日（两位数，首位不足补0） 小写
//	H：小时 带有首位零的 24 小时小时格式
//	i ：带有首位零的分钟
//	s ：带有首位零的秒（00 -59）
func Date(format string) string {
	now := time.Now() //获取时间

	//处理星期格式
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	format = strings.Replace(format, "Y", now.Format("2006"), -1)         //年
	format = strings.Replace(format, "m", now.Format("01"), -1)           //月
	format = strings.Replace(format, "w", fmt.Sprintf("%d", weekday), -1) //星期
	format = strings.Replace(format, "d", now.Format("02"), -1)           //日
	format = strings.Replace(format, "H", now.Format("15"), -1)           //小时
	format = strings.Replace(format, "i", now.Format("04"), -1)           //分钟
	format = strings.Replace(format, "s", now.Format("05"), -1)           //秒

	return format
}

// JsonEncode map转json
func JsonEncode(param interface{}) string {
	switch param.(type) {
	case map[string]interface{}:
		dataType, _ := json.Marshal(param)
		dataString := string(dataType)
		return dataString
	case []interface{}:
		dataType, _ := json.Marshal(param)
		dataString := string(dataType)
		return dataString
	default:
		return InterfaceToString(param)
	}
}

// JsonDecode json转map/array
func JsonDecode(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	if Empty(str) {
		return map[string]interface{}{}
	}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}

	tempMap = JsonDecodeHandling(tempMap)

	return tempMap
}

// JsonDecodeHandling 外包关于JsonDecode的特殊处理
func JsonDecodeHandling(tempMap map[string]interface{}) map[string]interface{} {
	// 遍历map并处理长数字以及其他可能的错误
	for key, value := range tempMap {
		tempMap[key] = JsonDecodeHandling2(value)
	}
	return tempMap
}

// JsonDecodeHandling2 外包关于JsonDecode的特殊处理2
func JsonDecodeHandling2(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return value.(string)
	case int:
		return strconv.Itoa(value.(int))
	case float32:
		if Float32Value, ok := value.(float32); ok {
			if Float32Value == 0 {
				//由于浮点数0的精度没有定义，被表示为空字符串，这里是对于0的处理
				return "0"
			} else {
				//对于长数字的处理
				float64str := float64(Float32Value)
				str := strconv.FormatFloat(float64str, 'f', -1, 32)
				return strings.TrimRight(strings.TrimRight(str, "0"), ".")
			}
		}
	case float64:
		if Float64Value, ok := value.(float64); ok {
			if Float64Value == 0 {
				//由于浮点数0的精度没有定义，被表示为空字符串，这里是对于0的处理
				return "0"
			} else {
				//对于长数字的处理
				str := strconv.FormatFloat(Float64Value, 'f', -1, 64)
				return strings.TrimRight(strings.TrimRight(str, "0"), ".")
			}
		}
	case nil:
		// 处理缺失字段或空值的情况
		switch v {
		case "age":
			return "N/A"
		default:
			return ""
		}
	case map[string]interface{}:
		// 如果本身就是map了就需要再深入挖掘一下了
		if mapSI, ok := value.(map[string]interface{}); ok {
			return JsonDecodeHandling(mapSI)
		}
	case []interface{}:
		//如果本身是数组的就需要进行多级循环了
		if arrI, ok := value.([]interface{}); ok {
			res := []interface{}{}
			for _, val := range arrI {
				res = append(res, JsonDecodeHandling2(val))
			}
			return res
		}

	default:
		SetLog([]interface{}{value, reflect.TypeOf(value)}, "解json遇到其他类型了")
		return value
	}
	return value
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
	case int64:
		res = strconv.FormatInt(inter.(int64), 10)
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

// InterfaceToFloat64 interface转float64
func InterfaceToFloat64(inter interface{}) float64 {
	return StringToFloat64(InterfaceToString(inter))
}

// ReflectValueToMap reflect.Value 对象 转map
func ReflectValueToMap(rv reflect.Value) map[string]interface{} {
	if !rv.IsValid() {
		//空处理
		return map[string]interface{}{}
	}

	result := make(map[string]interface{})
	if rv.Kind() == reflect.Map {
		for _, key := range rv.MapKeys() {
			value := rv.MapIndex(key).Interface()
			result[key.String()] = value
		}
	}

	return result
}

// InterfaceToArray interface转数组
func InterfaceToArray(inter interface{}) []interface{} {
	if arrI, ok := inter.([]interface{}); ok {
		return arrI
	} else {
		return []interface{}{}
	}
}

// ReflectValueToStr reflect.Value 对象 转string
func ReflectValueToStr(rv reflect.Value) string {
	if !rv.IsValid() {
		//空处理
		return ""
	}
	str, ok := rv.Interface().(string)
	if ok {
		return str
	} else {
		return ""
	}
}

// StringToInt64 string转int64
func StringToInt64(str string) int64 {
	ins64, _ := strconv.ParseInt(str, 10, 64)
	return ins64
}

// StringToBytes 将string转为[]byte
func StringToBytes(s string) []byte {
	RunBytes := []byte(s)
	return RunBytes
}

// StringToFloat64 将string转为float64
func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.00
	}
	return f
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

// Md5 把字符串转成md5
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// MakeRequest 发起http请求
// httpUrl    访问路径
// params 参数，该数组多于1个，表示为POST
// extend 请求伪造包头参数
// 返回的为一个请求状态，一个内容
func MakeRequest(httpUrl string, params map[string]interface{}, extend map[string]interface{}) map[string]interface{} {
	if httpUrl == "" {
		return map[string]interface{}{"code": "100"}
	}

	method := "GET"
	//参数数组多于1个，表示为POST
	if len(params) > 1 {
		method = "POST"
	}

	setData := url.Values{}
	for k, v := range params {
		setData.Set(k, InterfaceToString(v))
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, httpUrl, strings.NewReader(setData.Encode()))
	if err != nil {
		SetLog(err, "发起http请求错误1-创建请求失败")
	}

	//写入包头
	if len(params) > 1 {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (HTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	for k, v := range extend {
		req.Header.Set(k, InterfaceToString(v))
	}

	resp, err := client.Do(req)
	if err != nil {
		SetLog(err, "发起http请求错误2")
		return map[string]interface{}{
			"Code":   500,
			"Header": nil,
			"result": err.Error(),
		}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			SetLog(err, "发起http请求错误3")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		SetLog(err, "发起http请求错误4")
		return map[string]interface{}{
			"Code":   resp.StatusCode,
			"Header": resp.Header,
			"result": err.Error(),
		}
	}

	return map[string]interface{}{
		"Code":   resp.StatusCode,
		"Header": resp.Header,
		"result": string(body),
	}
}

// Empty 判断变量是否为空
func Empty(params interface{}) bool {
	switch val := params.(type) {
	case nil:
		return true
	case bool:
		return val == false
	case int, float64, float32:
		return val == 0
	case string:
		return val == ""
	case []interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	case chan interface{}:
		return val == nil
	case func() bool:
		return val == nil
	default:
		return reflect.ValueOf(val).IsZero()
	}
}

// FilePutContents 把一个字符串写入文件中
func FilePutContents(filePath string, content string, flags fs.FileMode) bool {
	err := os.WriteFile(filePath, []byte(content), flags)
	if err != nil {
		return false
	} else {
		return true
	}
}

// RedisRestrict 使用redis进行防刷限制
// @param key string 在redis中的key
// @param ttt int 超时时间
// @param msg string 自定义提示语
// @return bool
func RedisRestrict(SS ServerS, CH HttpInfo, key string, ttt time.Duration, msg string) bool {
	if Empty(msg) {
		msg = "您的请求太频繁了，请慢点"
	}
	//redis限制
	if SS.RDb.Connection("write").Exists(key) {
		OutJson(CH, map[string]interface{}{"code": "0", "msg": msg})
		return false
	} else {
		SS.RDb.Connection("write").Set(key, '1', ttt)
		return true
	}
}

// ArrayKeys 返回包含数组中所有键名的一个新数组
func ArrayKeys(sortedmap map[string]interface{}) []string {
	pairs := make([]string, 0, len(sortedmap))
	for k := range sortedmap {
		pairs = append(pairs, k)
	}

	return pairs
}

// Ksort 对关联数组按照键名进行升序排序：
func Ksort(sortedmap map[string]interface{}) []string {

	keys := ArrayKeys(sortedmap) //获取所有key
	sort.Strings(keys)           // 排序

	//由于map的键值对的迭代顺序是不确定的，所以只能返回数组
	return keys
}

// JoinHttpCode 把请求参数按规律拼接起来
// @param array paraMap 参数集
// @param boolean urlEncode 是否url转码
// @param boolean onEmpty 是否去除空值的元素
func JoinHttpCode(paraMap map[string]interface{}, urlEncode bool, onEmpty bool) string {
	buff := ""
	pairs := Ksort(paraMap)

	for _, k1 := range pairs {
		v1 := paraMap[k1]
		if onEmpty {
			if !Empty(v1) {
				v2 := InterfaceToString(v1)
				if urlEncode {
					v2 = url.QueryEscape(v2)
				}
				buff = buff + k1 + "=" + v2 + "&"
			}
		} else {
			v2 := InterfaceToString(v1)
			if urlEncode {
				v2 = url.QueryEscape(v2)
			}
			buff = buff + k1 + "=" + v2 + "&"
		}
	}

	length := len(buff)
	if length > 0 {
		buff = buff[:length-1]
	}

	return buff
}

// MapInterfaceCp 复制map，并且重新弄一个内存位储存
func MapInterfaceCp(data map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for key, value := range data {
		res[key] = value
	}
	return res
}

// Ternary 伪装的三元运算
func Ternary(condition interface{}, trueRun interface{}, falseRun interface{}) interface{} {
	if Empty(condition) {
		return falseRun
	} else {
		return trueRun
	}
}
