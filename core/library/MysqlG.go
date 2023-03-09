package library

import (
	"config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //导入包但不使用，init()
	"strconv"
)

// MysqlG 关于连接的公共函数
type MysqlG struct {
	Connections map[string]*sql.DB
}

func (r *MysqlG) InitMysql() {
	//先把内存分配定下来，不如会出现 panic: assignment to entry in nil map
	r.Connections = make(map[string]*sql.DB)

	//获取启动配置
	deploy := config.Mysql{}
	con := deploy.Run()

	//循环加载mysql数据库配置
	for k1, v1 := range con {
		//SetLog(v1, "输出配置详情")
		b := r.Connect(k1, v1)
		if !b {
			SetLog("加载配置"+k1+"失败", "加载配置失败")
			return
		}
	}
}

func (r *MysqlG) Connect(mName string, con map[string]string) bool {
	//[map[MAX_NUM:20 TIME_OUT:3 charset:utf8mb4 database:test host:127.0.0.1 password:123456 port:3306 timeout:5 user:root]]
	//用户名:密码啊@tcp(ip:端口)/数据库的名字
	dsn := con["user"] + ":" + con["password"] + "@tcp(" + con["host"] + ":" + con["port"] + ")/" + con["database"] + "?charset=" + con["charset"] + "&parseTime=True&loc=Local&timeout=" + con["TIME_OUT"] + "s&readTimeout=" + con["TIME_OUT"] + "s"

	//连接数据集
	db, err := sql.Open("mysql", dsn) //open不会检验用户名和密码
	if err != nil {
		SetLog(err, "数据库错误")
		//rr := Connection{db}
		//r.Connections[mName] = rr
		return false
	}

	err = db.Ping() //尝试连接数据库
	if err != nil {
		SetLog(err, "连接数据库错误")
		//rr := Connection{db}
		//r.Connections[mName] = rr
		return false
	}

	SetLog(err, "连接数据库"+con["database"]+"成功~")

	//设置数据库连接池的最大连接数
	MaxNum, _ := strconv.Atoi(con["MAX_NUM"])
	db.SetMaxIdleConns(MaxNum)

	r.Connections[mName] = new(sql.DB)
	r.Connections[mName] = db

	return true
}
