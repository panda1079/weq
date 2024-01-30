package Base

import (
	"core/library"
	"reflect"
	"strings"
)

// CtlBase 在这里写控制器的公共函数和组件
type CtlBase struct {
	SS library.ServerS
}

// SetTranslate 返回标注参数含义，好返回错误(中文提示，Util中的额外校验函数)
func (r *CtlBase) SetTranslate() map[string]map[string]string {
	return map[string]map[string]string{
		"username": {"name": "用户名", "extra": "IsUname", "tips": "用户名应为5-17个字母或数字字符"},
		"password": {"name": "密码", "extra": "IsPassword", "tips": "密码需以6-16位字母和数字的组合"},
		"sign":     {"name": "签名", "extra": ""},
	}
}

// CheckParameter 检查请求参数缺失
// @param  data  array  请求参数合集 [k1=>v1,k2=>v2]
// @param  need  array  需要检查的参数合集 [key1,key2]
// @return   bool
func (r *CtlBase) CheckParameter(CH library.HttpInfo, data map[string]interface{}, need []string) bool {
	var Translate = r.SetTranslate()

	for _, keySub := range need {
		var sense string
		value, ok := Translate[keySub]
		if ok {
			sense = value["name"]
		} else {
			sense = keySub
		}

		inspect := library.InterfaceToString(data[keySub])
		//存在性校验
		if library.Empty(inspect) || data[keySub] == "" {
			library.OutJson(CH, map[string]interface{}{"code": "0", "msg": "缺少" + sense + "参数"})
			return false
		}

		//额外函数校验
		if !library.Empty(Translate[keySub]) && Translate[keySub]["extra"] != "" {
			var fun = Translate[keySub]["extra"] //先撸出函数的字符串，免得识别不出来

			var (
				Verify     = &library.Verify{}
				methodArgs []reflect.Value
			)

			//压入数据
			methodArgs = append(methodArgs, reflect.ValueOf(data[keySub]))

			//反射调用(函数在 core/library/Verify.go)
			var ModBox = reflect.ValueOf(Verify).MethodByName(fun)
			var callOut = ModBox.Call(methodArgs)[0]

			//运行校验 (由于返回结果的类型比较坑，就需要做多几步了)
			if callOut.Kind() == reflect.Bool && !callOut.Bool() {

				if !library.Empty(Translate[keySub]) && Translate[keySub]["tips"] != "" {
					library.OutJson(CH, map[string]interface{}{"code": "0", "msg": "请输入正确的" + sense})
				} else {
					library.OutJson(CH, map[string]interface{}{"code": "0", "msg": Translate[keySub]["tips"]})
				}
				return false
			}
		}
	}
	return true
}

// GetSign 获取加密串
func (r *CtlBase) GetSign(params map[string]interface{}, key string) string {

	unSignParaString := library.JoinHttpCode(params, false, false)

	unSignParaString = strings.Replace(unSignParaString, "&#40;", "(", -1)
	unSignParaString = strings.Replace(unSignParaString, "&#41;", ")", -1)
	unSignParaString = strings.Replace(unSignParaString, "\\&quot;", "\"", -1)
	unSignParaString = strings.Replace(unSignParaString, "&#34;", "\"", -1)
	unSignParaString = strings.Replace(unSignParaString, "\\"+"\\", "\\", -1)

	return library.Md5(unSignParaString + "&key=" + key)
}
