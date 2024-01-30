package Test

import (
	"core/library"
)

type TestModel struct {
	SS library.ServerS
}

// Test 关于mysql的测试
func (r *TestModel) Test(Data map[string]interface{}) map[string]interface{} {

	library.SetLog(Data, "数据日志")

	kv := map[string]interface{}{
		"add_time": "123456789",
	}

	kv1 := map[string]interface{}{
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
func (r *TestModel) Test2(Data map[string]interface{}) map[string]interface{} {
	library.SetLog(Data, "数据日志")

	var key = "TestKey"

	library.SetLog(r.SS.RDb.Connection("write").Set(key, library.JsonEncode(Data), 0), "插入数据")
	library.SetLog(r.SS.RDb.Connection("write").Exists(key), "查询数据是否存在")
	library.SetLog(r.SS.RDb.Connection("write").Get(key), "查询数据")
	library.SetLog(r.SS.RDb.Connection("write").Del(key), "删除数据")
	library.SetLog(r.SS.RDb.Connection("write").Get(key), "查询数据")

	return map[string]interface{}{"code": "1", "msg": "子mod2拉起成功"}
}

// Test3 关于 Socket的测试
func (r *TestModel) Test3(msg []byte) ([]string, []byte) {
	//"SELECT `id`,`add_time` FROM `test` WHERE add_time = '123456789' "

	//res := library.MapToJson(r.SS.MDb.Connection("write").GetOne(string(msg))) //查出内容直接输出

	res := library.JsonEncode(map[string]interface{}{"code": 1, "data": "进入的数据：" + string(msg)})

	//自定义发送id对象(也可以加逻辑选择性推送)
	specific := []string{
		//"666",
	}

	return specific, library.StringToBytes(res)
}
