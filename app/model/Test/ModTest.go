package model

import (
	"core/library"
)

type ModTest struct {
}

func (r *ModTest) Test(Data map[string]string) map[string]interface{} {

	library.SetLog(Data, "数据日志")

	return map[string]interface{}{"code": "1", "msg": "子mod拉起成功"}
}
