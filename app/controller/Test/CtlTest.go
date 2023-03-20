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
		"a": CH.R("a", "a"),
		"b": CH.R("b", "b"),
		"c": CH.R("c", "c"),
		"d": CH.R("d", "d"),
		"e": CH.R("e", "e"),
		"f": CH.R("f", "f"),
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

	//library.OutJson(CH.ResponseWriter, aab) //输出到web页面
}

func (r *CtlTest) TestA(CH library.HttpInfo) {

	postData := map[string]interface{}{
		"a":  CH.R("a", "a"),
		"b":  CH.R("b", "b"),
		"c":  CH.R("c", "c"),
		"d":  CH.R("d", "d"),
		"e":  CH.R("e", "e"),
		"f":  CH.R("f", "f"),
		"ct": CH.R("ct", "ctc"),
		"ac": CH.R("ac", "aca"),
	}

	//library.OutJson(CH.ResponseWriter, postData) //输出到web页面

	params := make(map[string]interface{})
	extend := make(map[string]string)

	library.SetLog(library.MakeRequest("https://api.baidu.com", params, extend), "请求内容")
	library.OutHtml(CH.ResponseWriter, "test.html", postData)

	library.SetLog(CH.ClientRealIP(), "当前IP")

	library.SetLog(library.RandStr(10), "随机字符串")

}
