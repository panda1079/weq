package controller

import (
	"core/library"
	"fmt"
	"model"
	Order "model/Order"
	"reflect"
)

type CtlOrder struct {
	out map[string]string
}

func (r *CtlOrder) OrderList(CH library.HttpInfo) {

	postData := map[string]string{
		"a": CH.R("a", "a"),
		"b": CH.R("b", "b"),
		"c": CH.R("c", "c"),
		"d": CH.R("d", "d"),
		"e": CH.R("e", "e"),
		"f": CH.R("f", "f"),
	}

	fmt.Print(postData)
	fmt.Print("\n")

	//------------- 以下是抽象化调用 --------------//

	//初始化MOD池
	var Mod = model.ModIndex{}
	var RegisterMessage = Mod.Init()

	var ModName = "ModOrder"
	var ModFunName = "OrderList"
	var methodArgs []reflect.Value

	//压入数据
	methodArgs = append(methodArgs, reflect.ValueOf(postData))

	//反射调用
	var ModBox = reflect.ValueOf(RegisterMessage[ModName]).MethodByName(ModFunName)
	var aaa = ModBox.Call(methodArgs)[0]
	fmt.Print(aaa)
	fmt.Print("\n")

	//------------- 以上是抽象化调用 --------------//

	//------------- 以下的正常拉起 --------------//

	Mod2 := Order.ModOrder{}
	var aab = Mod2.OrderList(postData) //直接调用
	fmt.Print(aab)
	fmt.Print("\n")

	//------------- 以上的正常拉起 --------------//

	library.OutJson(CH.ResponseWriter, postData) //输出到web页面
}

func (r *CtlOrder) TestA(CH library.HttpInfo) {

	postData := map[string]string{
		"a":  CH.R("a", "a"),
		"b":  CH.R("b", "b"),
		"c":  CH.R("c", "c"),
		"d":  CH.R("d", "d"),
		"e":  CH.R("e", "e"),
		"f":  CH.R("f", "f"),
		"ct": CH.R("ct", "ctc"),
		"ac": CH.R("ac", "aca"),
	}

	//library.SetLog(postData)
	library.OutJson(CH.ResponseWriter, postData) //输出到web页面

	library.SetLog(CH.GetHeader("Content-Type") == "application/json; charset=utf-8")
	library.SetLog(CH.GetHeader("content-type"))
	library.SetLog(CH.GetHeader("content-types"))
	//fmt.Print(CH.GetHeader("Content-Type"))
	//fmt.Print("\n")
	//fmt.Print(CH.GetHeader("content-type"))

}
