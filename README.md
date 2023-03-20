# TEST_GO
## 基于golang 1.18

### 运行命令
```
go run main.go
```

### 测试URL
```
http://127.0.0.1:9091/CtlTest/TestA
http://127.0.0.1:9091/test/test
```

### 一个简单的类PHP的golang结构
1.R函数可以直接获取url中的自定义参数以及GET/POST和raw的json参数

2.Go语言官方没有实现Mysql的数据库驱动,database/sql包提供了保证SQL或类SQL数据库的泛用接口,并不提供具体的数据库驱动。需要安装Go Mysql Driver 依赖
```
go get -u github.com/go-sql-driver/mysql
```
 
