package library

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Util 公共工具

func SetLog(a ...any) {
	fmt.Print(a)
	fmt.Print("\n")
}

func OutJson(w http.ResponseWriter, OutData any) {
	jsonBytes, err := json.Marshal(OutData) //转换json
	if err != nil {
		fmt.Println(err)
		fprintf, err := fmt.Fprintf(w, "{\"code\": \"1\", \"route\": \"输出错误\"}") // 发送响应到客户端
		if err != nil {
			SetLog(fprintf)
			return
		}
		return
	}

	fprintf, err := fmt.Fprintf(w, string(jsonBytes)) // 发送响应到客户端
	if err != nil {
		SetLog(fprintf)
		return
	}
}
