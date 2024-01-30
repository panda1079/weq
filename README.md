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

### socket模式测试（socket是延长的tcp，就直接做到一起了）
```
go run main.go

//id 可以为选择性推送的标记
// 测试链接 ： http://127.0.0.1:9091/test?id=666
```

### 一个简单的类PHP的golang结构
1.R函数可以直接获取url中的自定义参数以及GET/POST和raw的json参数

2.Go语言官方没有实现Mysql的数据库驱动,database/sql包提供了保证SQL或类SQL数据库的泛用接口,并不提供具体的数据库驱动。需要安装Go Mysql Driver 依赖
```
go get -u github.com/go-sql-driver/mysql
go get -u github.com/go-redis/redis
```
 
3.Go语言没有完善的web socker包，需要导入websocket (待完善)
```
go get github.com/gorilla/websocket
```

4. 如果国内的get不到，那就需要使用国内源
```
1. 七牛 CDN
go env -w  GOPROXY=https://goproxy.cn,direct

2. 阿里云
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

3. 官方
go env -w  GOPROXY=https://goproxy.io,direct
```

1.写一个可以每分钟都在检测go程序是否运行的成勋
2.通过crontab进行进程守护
