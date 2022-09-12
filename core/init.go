package core

import (
	"fmt"
	"log"
	"net/http"
)

func Load_route() {
	fmt.Println(routes.routes)

	http.HandleFunc("/", sayHelloWorld)
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "1231231321") // 发送响应到客户端
}
