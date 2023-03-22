package controller

import (
	"core/library"
	"model"
	ModTest "model/Test"
	"reflect"
)

type CtlTest struct {
	SS  library.ServerS
	out map[string]string
}

func (r *CtlTest) Test(CH library.HttpInfo) {

	postData := map[string]interface{}{
		"a": CH.R("a", 222, "int"),
		"b": CH.R("b", "b", "string"),
		"c": CH.R("c", 444, "int"),
		"d": CH.R("d", "d", "string"),
		"e": CH.R("e", "e", "string"),
		"f": CH.R("f", "f", "string"),
	}

	library.SetLog(postData, "输出打印")

	//------------- 以下是抽象化调用 --------------//

	//初始化MOD池
	var (
		Mod             = model.ModIndex{}
		RegisterMessage = Mod.Init(r.SS)
	)

	var (
		ModName    = "ModTest"
		ModFunName = "Test"
		methodArgs []reflect.Value
	)

	//压入数据
	methodArgs = append(methodArgs, reflect.ValueOf(postData))

	//反射调用
	var ModBox = reflect.ValueOf(RegisterMessage[ModName]).MethodByName(ModFunName)
	var aaa = ModBox.Call(methodArgs)[0]
	library.SetLog(aaa, "反射调用输出打印")

	//------------- 以上是抽象化调用 --------------//

	//------------- 以下的正常拉起 --------------//

	Mod2 := ModTest.ModTest{SS: r.SS}
	var aab = Mod2.Test2(postData) //直接调用
	library.SetLog(aab, "直接调用输出打印")

	//------------- 以上的正常拉起 --------------//

	library.OutJson(CH, aab) //输出到web页面
}

func (r *CtlTest) TestA(CH library.HttpInfo) {

	postData := map[string]interface{}{
		"a":  CH.R("a", 1, "int"),
		"b":  CH.R("b", "b", "string"),
		"c":  CH.R("c", 222, "int"),
		"d":  CH.R("d", "d", "string"),
		"e":  CH.R("e", "e", "string"),
		"f":  CH.R("f", "f", "string"),
		"ct": CH.R("ct", "ctc", "string"),
		"ac": CH.R("ac", "aca", "string"),
	}

	library.SetLog(postData, "postData") //输出到日志

	//library.OutJson(CH, postData) //输出到web页面

	params := make(map[string]interface{})
	extend := make(map[string]string)

	library.SetLog(library.MakeRequest("https://api.baidu.com", params, extend), "请求内容")
	library.OutHtml(CH, "test.html", postData) //输出html（允许轻微替换内容）

	library.SetLog(CH.ClientRealIP(), "当前IP")

	library.SetLog(library.RandStr(10, true), "随机字符串")

}
