package model

import (
	"core/library"
	"fmt"
)

type ModTest struct {
	DB library.MysqlG
}

func (r *ModTest) Test(Data map[string]string) map[string]interface{} {

	library.SetLog(Data, "数据日志")
	// 查询数据
	var aaa = r.DB.Connection("write").GetAll("SELECT `id`,`name`,`system` FROM h_game_client")

	fmt.Print(aaa)

	return map[string]interface{}{"code": "1", "msg": "子mod拉起成功"}

}
