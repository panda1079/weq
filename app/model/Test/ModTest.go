package model

import (
	"core/library"
)

type ModTest struct {
	DB library.MysqlG
}

func (r *ModTest) Test(Data map[string]string) map[string]interface{} {

	library.SetLog(Data, "数据日志")

	kv := map[string]string{
		"add_time": "123456789",
	}

	kv1 := map[string]string{
		"add_time": "987654321",
	}

	// 查询数据
	library.SetLog(r.DB.Connection("write").Insert(kv, "h_game_client"), "插入数据")
	library.SetLog(r.DB.Connection("write").GetOne("SELECT `id`,`add_time` FROM `test` WHERE add_time = '123456789' "), "输出数据")
	library.SetLog(r.DB.Connection("write").UpDate(kv1, "`add_time`='123456789'", "h_game_client", true), "输出数据")
	library.SetLog(r.DB.Connection("write").GetOne("SELECT `id`,`add_time` FROM `test` WHERE add_time = '987654321' "), "输出数据")
	library.SetLog(r.DB.Connection("write").Delete("`add_time`='987654321'", "h_game_client"), "输出数据")
	library.SetLog(r.DB.Connection("write").GetOne("SELECT `id`,`add_time` FROM `test` WHERE add_time = '987654321' "), "输出数据")

	return map[string]interface{}{"code": "1", "msg": "子mod拉起成功"}
}
