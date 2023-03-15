这里放路由文件[README.md](README.md)

### 路由例子
```
AddRe.Get("路由", "函数名", "控制器名")

	AddRe.Get("/order/order_list", "OrderList", "CtlOrder")
	AddRe.Post("/order/order_list", "OrderList", "CtlOrder")
	AddRe.AddRoute([2]string{"Post", "Get"}, "/test", "TestA", "CtlOrder")
	AddRe.AddRoute([2]string{"Post", "Get"}, "/(?P<ct>[a-z|A-Z|\\d+]+)/(?P<ac>[a-z|A-Z|\\d+]+)", "{ac}", "{ct}")
``` 

### 路由说明
	1.路由正则遵循golang分组命名规则
	2.分组命名获取到的内容会对{}内的内容进行替换
	3.分组命名获取到的参数可以直接使用R函数获取到

