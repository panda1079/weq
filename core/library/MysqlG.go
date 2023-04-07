package library

import (
	"config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //导入包但不使用，init()
	"strconv"
	"strings"
)

// MysqlG 关于连接的公共函数
type MysqlG struct {
	Connections    map[string]*sql.DB
	connectionName string //当前连接池名称
}

// InitMysql 加载配置，装载链接
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
			SetLog("加载Mysql配置"+k1+"失败", "加载配置失败")
			return
		}
	}
}

// Connect 绑定解析，并返回对象
func (r *MysqlG) Connect(mName string, con map[string]string) bool {
	//[map[MAX_NUM:20 TIME_OUT:3 charset:utf8mb4 database:test host:127.0.0.1 password:123456 port:3306 timeout:5 user:root]]
	//用户名:密码啊@tcp(ip:端口)/数据库的名字
	dsn := con["user"] + ":" + con["password"] + "@tcp(" + con["host"] + ":" + con["port"] + ")/" + con["database"] + "?charset=" + con["charset"] + "&parseTime=True&loc=Local&timeout=" + con["TIME_OUT"] + "s&readTimeout=" + con["TIME_OUT"] + "s"

	//连接数据集
	db, err := sql.Open("mysql", dsn) //open不会检验用户名和密码
	if err != nil {
		SetLog(err, "Mysql错误")
		return false
	}

	err = db.Ping() //尝试连接数据库
	if err != nil {
		SetLog(err, "连接Mysql错误")
		return false
	}

	SetLog(err, "连接Mysql-"+con["database"]+"成功~")

	//设置数据库连接池的最大连接数
	MaxNum, _ := strconv.Atoi(con["MAX_NUM"])
	db.SetMaxIdleConns(MaxNum)

	//r.Connections[mName] = new(sql.DB)
	r.Connections[mName] = db

	return true
}

// Connection 连接名设置
func (r *MysqlG) Connection(name string) *MysqlG {
	r.connectionName = name
	return r
}

// Execute 执行一条sql语句，注意，此步操作没有经过参数绑定，有风险，最好还是使用delete,update,insert等
func (r *MysqlG) Execute(SqlStr string) (bool, sql.Result) {
	result, err := r.Connections[r.connectionName].Exec(SqlStr)
	isOk := true
	if err != nil {
		isOk = false
		SetLog(err, "Mysql操作数据错误")
	}
	return isOk, result
}

// GraveAccent 自动加上`符号
func (r *MysqlG) GraveAccent(item string) string {
	if strings.ToUpper(item[0:1]) != "`" {
		return "`" + item + "`"
	}
	return item
}

// GetAll 获取多条的数据
func (r *MysqlG) GetAll(sqlStr string) (rows []map[string]interface{}) {
	list, _ := r.Connections[r.connectionName].Query(sqlStr) //把数据查出来
	fields, _ := list.Columns()                              //返回列名

	//逐行读取并插入
	for list.Next() {
		scans := make([]interface{}, len(fields))
		row := make(map[string]interface{})

		for i := range scans {
			scans[i] = &scans[i]
		}

		list.Scan(scans...)

		for i, v := range scans {
			var value string
			if v != nil {
				value = fmt.Sprintf("%s", v)
			}
			row[fields[i]] = value
		}

		rows = append(rows, row)
	}

	return rows
}

// GetOne 获取单条的数据
func (r *MysqlG) GetOne(SqlStr string) map[string]interface{} {
	// 去除连续空格
	SqlStr = strings.TrimSpace(SqlStr)
	Len := len(SqlStr)
	SqlSUb := strings.ToUpper(SqlStr[Len-7 : Len])
	if SqlSUb != "LIMIT 1" {
		SqlStr += " LIMIT 1 "
	}

	//查询数据
	res := r.GetAll(SqlStr)
	if len(res) > 0 {
		return res[0]
	}

	return map[string]interface{}{}
}

// UpDate 以新的$key_values更新mysql数据
// keyValues array('aid'=>1,'cid'=>2)
// where  string 字符语句
// tableName 表名
// updateOne 是否只更新一行数据
func (r *MysqlG) UpDate(keyValues map[string]string, Where string, tableName string, updateOne bool) int64 {
	if len(keyValues) == 0 {
		return 0
	}

	SqlStr := "UPDATE " + r.GraveAccent(tableName) + " SET "

	for k, v := range keyValues {
		SqlStr += r.GraveAccent(k) + " = '" + v + "',"
	}
	SqlStr = SqlStr[0 : len(SqlStr)-1]

	SqlStr += " WHERE " + Where

	// 默认只更新一条
	if updateOne {
		SqlStr = strings.TrimSpace(SqlStr)
		SqlLen := len(SqlStr)
		SqlSUb := strings.ToUpper(SqlStr[SqlLen-7 : SqlLen])
		if SqlSUb != "LIMIT 1" {
			SqlStr += " LIMIT 1 "
		}
	}

	//SetLog(SqlStr, "打印sql") //正式上限需要去掉

	isOk, result := r.Execute(SqlStr)
	if !isOk {
		return 0
	}

	i, _ := result.RowsAffected() //受影响行数
	return i
}

// Insert 插入一条新的数据
// keyValues array
// tableName 表名
func (r *MysqlG) Insert(keyValues map[string]string, tableName string) int64 {
	if len(keyValues) == 0 {
		return 0
	}

	var (
		keySql   = ""
		valueSql = ""
	)

	for k, v := range keyValues {
		keySql += r.GraveAccent(k) + ","
		valueSql += "'" + v + "'" + ","
	}
	keySql = keySql[0 : len(keySql)-1]
	valueSql = valueSql[0 : len(valueSql)-1]

	SqlStr := "INSERT INTO" + " " + r.GraveAccent(tableName) + " (" + keySql + ") VALUES (" + valueSql + ")"

	//SetLog(SqlStr, "打印sql") //正式上限需要去掉

	isOk, result := r.Execute(SqlStr)
	if !isOk {
		return 0
	}

	id, _ := result.LastInsertId() //新增数据的ID
	return id
}

// MultiInsert 多条记录同时insert
// keyValues array
// tableName 表名
// ignore 是否添加IGNORE关键字
func (r *MysqlG) MultiInsert(keyValues []map[string]string, tableName string) int64 {
	if len(keyValues) == 0 {
		return 0
	}

	var (
		keyS     []string
		valueStr []string
	)

	for k, _ := range keyValues[0] {
		keyS = append(keyS, r.GraveAccent(k))
	}

	for _, v1 := range keyValues {
		var v1list []string
		for _, v2 := range v1 {
			v1list = append(v1list, v2)
		}
		valueStr = append(valueStr, "('"+strings.Join(v1list, "','")+"')") //这里需要数组转字符串
	}

	SqlStr := "INSERT INTO" + " " + r.GraveAccent(tableName) + " (" + strings.Join(keyS, ",") + ") VALUES " + strings.Join(valueStr, ",")

	//SetLog(SqlStr, "打印sql") //正式上限需要去掉

	isOk, result := r.Execute(SqlStr)
	if !isOk {
		return 0
	}

	i, _ := result.RowsAffected() //受影响行数
	return i
}

// Delete 删除数据
// where  string 字符语句
// tableName string 表名
func (r *MysqlG) Delete(Where string, tableName string) int64 {
	if len(Where) == 0 {
		return 0
	}

	SqlStr := "DELETE FROM" + " " + r.GraveAccent(tableName) + " WHERE " + Where

	//SetLog(SqlStr, "打印sql") //正式上限需要去掉

	isOk, result := r.Execute(SqlStr)
	if !isOk {
		return 0
	}

	i, _ := result.RowsAffected()
	return i
}
