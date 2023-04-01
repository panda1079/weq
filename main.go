package main

import (
	"core"
	"flag"
)

//var m = string(*flag.String("m", "http", "string类型参数")) //设定启动模式
var m1 = flag.String("m", "http", "string类型参数") //设定启动模式

var ct1 = flag.String("ct", "", "string类型参数")     //设定控制器
var ac1 = flag.String("ac", "", "string类型参数")     //设定方法
var from1 = flag.String("from", "", "string类型参数") //设定其他参数（x-www-from-urlencoded格式）

func main() {
	// 获取内容
	flag.Parse()
	var m = string(*m1)

	if m == "http" {
		core.InitHttp()
	} else if m == "cli" {
		core.InitCli(string(*ct1), string(*ac1), string(*from1))
	}

}
