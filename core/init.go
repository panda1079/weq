package core

import (
	"fmt"
	"net/http"
)

func LoadRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.Method)
	fmt.Print("==========")
	fmt.Fprintf(w, "1231231321") // 发送响应到客户端
}

func Init() {

	http.HandleFunc("/", LoadRoute)
	http.ListenAndServe(":9091", nil)
}
