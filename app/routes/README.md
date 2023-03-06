这里放路由文件[README.md](README.md)

### 路由说明
```
AddRe.Get("路由", "函数名", "控制器名")

AddRe.Get("/aaa/bbb", "aaa", "bbb")
AddRe.Post("/aaa/bbb", "aaa", "bbb")
AddRe.AddRoute([2]string{"Post", "Get"}, "/test", "TestA", "bbb")

AddRe.Get("/aaa/{id}/bbb", "index", "$id")
AddRe.Post("/aaa/{id}/{ccc}", "$ccc", "$id")
AddRe.AddRoute([2]string{"Post", "Get"}, "/test/{aaa}/{bbb}", "$aaa", "$bbb")
``` 
