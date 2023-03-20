package model

import (
	"core/library"
)

type ModTest struct {
	SS library.ServerS
}

// Test 关于mysql的测试
func (r *ModTest) Test(Data map[string]interface{}) map[string]interface{} {

	library.SetLog(Data, "数据日志")

	kv := map[string]string{
		"add_time": "123456789",
	}

	kv1 := map[string]string{
		"add_time": "987654321",
	}

	// 查询数据
	library.SetLog(r.SS.MDb.Connection("write").Insert(kv, "test"), "插入数据")
	library.SetLog(r.SS.MDb.Connection("write").GetOne("SELECT `id`,`add_time` FROM `test` WHERE add_time = '123456789' "), "输出数据")
	library.SetLog(r.SS.MDb.Connection("write").UpDate(kv1, "`add_time`='123456789'", "test", true), "输出数据")
	library.SetLog(r.SS.MDb.Connection("write").GetOne("SELECT `id`,`add_time` FROM `test` WHERE add_time = '987654321' "), "输出数据")
	library.SetLog(r.SS.MDb.Connection("write").Delete("`add_time`='987654321'", "test"), "输出数据")
	library.SetLog(r.SS.MDb.Connection("write").GetOne("SELECT `id`,`add_time` FROM `test` WHERE add_time = '987654321' "), "输出数据")

	return map[string]interface{}{"code": "1", "msg": "子mod1拉起成功"}
}

// Test2 关于Redis的测试
func (r *ModTest) Test2(Data map[string]interface{}) map[string]interface{} {
	library.SetLog(Data, "数据日志")

	var key = "TestKey"

	library.SetLog(r.SS.RDb.Connection("write").Set(key, library.MapToJson(Data), 0), "插入数据")
	library.SetLog(r.SS.RDb.Connection("write").Exists(key), "查询数据是否存在")
	library.SetLog(r.SS.RDb.Connection("write").Get(key), "查询数据")
	library.SetLog(r.SS.RDb.Connection("write").Del(key), "删除数据")
	library.SetLog(r.SS.RDb.Connection("write").Get(key), "查询数据")
	return map[string]interface{}{"code": "1", "msg": "子mod2拉起成功"}

}
