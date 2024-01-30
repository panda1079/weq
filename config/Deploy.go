package config

import (
	"strconv"
	"time"
)

type Deploy struct {
}

func (r *Deploy) Run() map[string]string {

	Year := strconv.Itoa(time.Now().Year())
	Month := time.Now().Format("01")
	return map[string]string{
		"LISTEN_ADDRESS": "0.0.0.0",
		"PORT":           "9000",
		"LOG_DIR":        "log/server-" + Year + "-" + Month + ".log", //日志写入位置
	}
}
