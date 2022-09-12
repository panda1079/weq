package main

import (
	// "fmt"
	// "log"
	// "net/http"
	// "strings"

	core "core"
)

var DEMO = 1 // 1:测试环境, 2:加载正式环境配置 0:正式运行

func main() {
	core.Load_route()
}

// func sayHelloWorld(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm() // 解析参数

// 	fmt.Println(r.Form)             // 在服务端打印请求参数
// 	fmt.Println("URL:", r.URL.Path) // 请求 URL
// 	fmt.Println("Scheme:", r.URL.Scheme)

// 	for k, v := range r.Form {
// 		fmt.Println(k, ":", strings.Join(v, ""))
// 	}
// 	fmt.Fprintf(w, "你好，学院君！") // 发送响应到客户端
// }

// func maina() {
// 	core.Maina()
// 	http.HandleFunc("/", sayHelloWorld)
// 	err := http.ListenAndServe(":9091", nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
