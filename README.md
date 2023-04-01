# TEST_GO
## 基于golang 1.18

### 运行命令
```
go run main.go
```

### URL模式测试
```
http://127.0.0.1:9091/CtlTest/TestA
http://127.0.0.1:9091/test/test
```

### cli模式测试
```
// ac:控制器
// ct：控制器内的方法
// from：使用get形式带入的参数（复杂的字符串请url编码之后再输入）

go run main.go -m cli -ct CtlTest -ac TestA -from "a=1&b=2&c=3&d=abc&d=%E6%98%AF%E4%B8%AD%E6%96%87url&e=是中文"
```

### 一个简单的类PHP的golang结构
1.R函数可以直接获取url中的自定义参数以及GET/POST和raw的json参数

2.Go语言官方没有实现Mysql的数据库驱动,database/sql包提供了保证SQL或类SQL数据库的泛用接口,并不提供具体的数据库驱动。需要安装Go Mysql Driver 依赖
```
go get -u github.com/go-sql-driver/mysql
```
 
