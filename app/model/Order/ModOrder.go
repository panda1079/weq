package model

import "fmt"

type ModOrder struct {
}

func (r *ModOrder) OrderList(Data map[string]string) map[string]string {
	fmt.Print(Data)
	return map[string]string{"code": "1", "msg": "子mod拉起成功"}
}
