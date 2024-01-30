package library

import (
	"config"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" //导入包但不使用，init()
)

// MysqlG 关于连接的公共函数
type MysqlG struct {
	Connections    map[string]*sql.DB
	connectionName string //当前连接池名称

	Wheres   string        //ORM 查询条件
	Table    string        //ORM 查询数据表
	VisibleS []interface{} //ORM 返回内容结果
	GroupS   string        //ORM 分组条件
	OrderS   string        //ORM 排序条件
	OrmSql   string        //ORM 最后拼接的sql
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

	SetLog(err, "连接Mysql - "+con["database"]+" - "+mName+" 成功~")

	//设置数据库连接池的最大连接数
	MaxNum, _ := strconv.Atoi(con["MAX_NUM"])
	db.SetMaxIdleConns(MaxNum)

	//r.Connections[mName] = new(sql.DB)
	r.Connections[mName] = db

	return true
}

// Connection 连接名设置
func (r *MysqlG) Connection(name string) *MysqlG {

	//初始化一下
	r.Wheres = ""                //ORM 查询条件
	r.Table = ""                 //ORM 查询数据表
	r.VisibleS = []interface{}{} //ORM 返回内容结果
	r.GroupS = ""                //ORM 分组条件
	r.OrderS = ""                //ORM 排序条件
	r.OrmSql = ""                //ORM 最后拼接的sql

	r.connectionName = name
	return r
}

// GetConnName 获取connectionName
func (r *MysqlG) GetConnName() string {
	if Empty(r.connectionName) {
		return "write"
	}
	return r.connectionName
}

// Execute 执行一条sql语句，注意，此步操作没有经过参数绑定，有风险，最好还是使用delete,update,insert等
func (r *MysqlG) Execute(SqlStr string) (isOk bool, result sql.Result) {

	//函数执行完成后被调用,通过 recover() 函数捕获到之前发生的错误。
	defer func(SqlStr string) {
		if r := recover(); r != nil {
			SetLog([]interface{}{SqlStr, r}, "执行sql语句错误")
			isOk = false
		}
	}(SqlStr)

	result, err := r.Connections[r.GetConnName()].Exec(SqlStr)
	isOk = true
	if err != nil {
		isOk = false
		SetLog([]interface{}{SqlStr, err}, "Mysql操作数据错误")
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
func (r *MysqlG) GetAll(SqlStr string) (rows []map[string]interface{}) {

	//函数执行完成后被调用,通过 recover() 函数捕获到之前发生的错误。
	defer func(sqlStr string) {
		if r := recover(); r != nil {
			SetLog([]interface{}{SqlStr, r}, "执行sql语句错误")
			rows = []map[string]interface{}{}
		}
	}(SqlStr)

	list, _ := r.Connections[r.GetConnName()].Query(SqlStr) //把数据查出来
	fields, _ := list.Columns()                             //返回列名

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

// WhereIntToStr 把查询条件的不同格式转换成字符串
func (r *MysqlG) WhereIntToStr(Where interface{}) string {
	WhereStr := ""

	switch Where.(type) {
	case string:
		WhereStr = " WHERE " + InterfaceToString(Where)
	case []interface{}:
		if vv2, ok1 := Where.([]interface{}); ok1 {
			for _, v2 := range vv2 {
				// 对数组的处理
				if vvv2, ok2 := v2.([]string); ok2 {
					if len(vvv2) == 2 {
						r.Where(vvv2[0], vvv2[1])
					} else if len(vvv2) == 3 {
						r.Where(vvv2[0], vvv2[1], vvv2[2])
					} else if len(vvv2) == 4 {
						r.Where(vvv2[0], vvv2[1], vvv2[2], vvv2[3])
					}
				}
				// 对string的处理
				if vvv2, ok3 := v2.(string); ok3 {
					r.Where(vvv2)
				}
				// 对map的处理
				if vvv2, ok4 := v2.(map[string]interface{}); ok4 {
					for kkkk2, vvvv2 := range vvv2 {
						r.Where(kkkk2, vvvv2)
					}
				}
			}
		}
		if !Empty(r.Wheres) {
			WhereStr = strings.Replace(r.Wheres, "AND", " WHERE ", 1)
		}
	case map[string]interface{}:
		if vv3, ok := Where.(map[string]interface{}); ok {
			for k3, v3 := range vv3 {
				r.Where(k3, v3)
			}
		}
		if !Empty(r.Wheres) {
			WhereStr = strings.Replace(r.Wheres, "AND", " WHERE ", 1)
		}
	}

	//不允许无条件更新
	if WhereStr == "" || len(WhereStr) < 10 {
		return ""
	}

	return WhereStr
}

// UpDate 以新的$key_values更新mysql数据
// keyValues array('aid'=>1,'cid'=>2)
// where  string 字符语句
// tableName 表名
// updateOne 是否只更新一行数据
func (r *MysqlG) UpDate(keyValues map[string]interface{}, Where interface{}, tableName string, updateOne bool) int64 {
	if len(keyValues) == 0 {
		return 0
	}

	SqlStr := "UPDATE " + r.GraveAccent(tableName) + " SET "

	for k, v := range keyValues {
		SqlStr += r.GraveAccent(k) + " = '" + InterfaceToString(v) + "',"
	}
	SqlStr = SqlStr[0 : len(SqlStr)-1]

	SqlStr += r.WhereIntToStr(Where)

	// 默认只更新一条
	if updateOne {
		SqlStr = strings.TrimSpace(SqlStr)
		SqlLen := len(SqlStr)
		SqlSUb := strings.ToUpper(SqlStr[SqlLen-7 : SqlLen])
		if SqlSUb != "LIMIT 1" {
			SqlStr += " LIMIT 1 "
		}
	}

	// SetLog(SqlStr, "打印sql") //正式上限需要去掉

	isOk, result := r.Execute(SqlStr)
	if !isOk {
		return 0
	}

	i, _ := result.RowsAffected() //受影响行数

	// SetLog(i, "更新数据受影响行数") //正式上限需要去掉

	return i
}

// Insert 插入一条新的数据
// keyValues array
// tableName 表名
func (r *MysqlG) Insert(keyValues map[string]interface{}, tableName string) int64 {
	if len(keyValues) == 0 {
		return 0
	}

	var (
		keySql   = ""
		valueSql = ""
	)

	for k, v := range keyValues {
		keySql += r.GraveAccent(k) + ","
		valueSql += "'" + InterfaceToString(v) + "'" + ","
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
func (r *MysqlG) MultiInsert(keyValues []map[string]interface{}, tableName string) int64 {
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
			v1list = append(v1list, InterfaceToString(v2))
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
func (r *MysqlG) Delete(Where interface{}, tableName string) int64 {

	SqlStr := "DELETE FROM" + " " + r.GraveAccent(tableName) + r.WhereIntToStr(Where)

	//SetLog(SqlStr, "打印sql") //正式上限需要去掉

	isOk, result := r.Execute(SqlStr)
	if !isOk {
		return 0
	}

	i, _ := result.RowsAffected()
	return i
}

// Limit 根据当前页面和每页显示数目产生SQL的limit（方便拼接使用,非执行）
// page 当前页码
// pageNum 显示数目
// @return string 返回LIMIT语句
func (r *MysqlG) Limit(page interface{}, pageNum interface{}) string {

	//转换一下page类型
	pageInt, err := strconv.Atoi(fmt.Sprintf("%v", page))
	if err != nil {
		return " LIMIT 0,0 "
	}

	//转换一下pageNumInt类型
	pageNumInt, err := strconv.Atoi(fmt.Sprintf("%v", pageNum))
	if err != nil {
		return " LIMIT 0,0 "
	}

	if pageInt < 1 {
		pageInt = 1
	}
	offset := (pageInt - 1) * pageNumInt

	return " LIMIT " + strconv.FormatInt(int64(offset), 10) + ", " + strconv.FormatInt(int64(pageNumInt), 10) + " "
}

// ========================= 下面是有关于ORM的尝试操作 =========================//
// ============== ORM只支持简单的拼接查询工作，复杂查询请使用原生sql ===============//

// TableName ORM设置数据表
func (r *MysqlG) TableName(table string) *MysqlG {
	r.Table = table
	return r
}

// Visible ORM 自定义返回内容结果
func (r *MysqlG) Visible(Visible []interface{}) *MysqlG {
	r.VisibleS = Visible
	return r
}

// Where ORM 添加查询条件
func (r *MysqlG) Where(where ...interface{}) *MysqlG {
	var Woolen = len(where)
	if Woolen == 1 {
		//只有一个的，就是直接字符串的条件
		r.Wheres = r.Wheres + " AND " + InterfaceToString(where[0]) + " "
	} else if Woolen == 2 {
		//有两个的就是and条件
		r.Wheres = r.Wheres + " AND " + r.GraveAccent(InterfaceToString(where[0])) + " = '" + InterfaceToString(where[1]) + "' "
	} else if Woolen == 3 {
		// 有三个那就是自定义条件了
		r.Wheres = r.Wheres + " AND " + r.GraveAccent(InterfaceToString(where[0])) + " " + InterfaceToString(where[2]) + " '" + InterfaceToString(where[1]) + "' "
	} else if Woolen == 4 {
		// 有四个那就是不但自定义条件，而且要自定义与或
		r.Wheres = r.Wheres + " " + InterfaceToString(where[3]) + " " + r.GraveAccent(InterfaceToString(where[0])) + " " + InterfaceToString(where[2]) + " '" + InterfaceToString(where[1]) + "' "
	}
	return r
}

// WhereBatch ORM 批量添加内容
func (r *MysqlG) WhereBatch(where []interface{}) *MysqlG {

	for _, v1 := range where {
		switch v1.(type) {
		case string:
			// 直接添加字符串条件
			// 例子："add_time = 1656904108"
			r.Where(v1)
		case []interface{}:
			// 多个复杂内容通过数组形式添加，循环还得循环 例子：
			//where := []interface{}{
			//	[]string{"id", "37", ">"},
			//	[]string{"id", "38"},
			//	"add_time = 1656904108",
			//	map[string]interface{}{"modifie_time": "1657019178", "product_id": "1"},
			//}

			if vv2, ok1 := v1.([]interface{}); ok1 {
				for _, v2 := range vv2 {
					// 对数组的处理
					if vvv2, ok2 := v2.([]string); ok2 {
						if len(vvv2) == 2 {
							r.Where(vvv2[0], vvv2[1])
						} else if len(vvv2) == 3 {
							r.Where(vvv2[0], vvv2[1], vvv2[2])
						} else if len(vvv2) == 4 {
							r.Where(vvv2[0], vvv2[1], vvv2[2], vvv2[3])
						}
					}
					// 对string的处理
					if vvv2, ok3 := v2.(string); ok3 {
						r.Where(vvv2)
					}
					// 对map的处理
					if vvv2, ok4 := v2.(map[string]interface{}); ok4 {
						for kkkk2, vvvv2 := range vvv2 {
							r.Where(kkkk2, vvvv2)
						}
					}
				}
			}
		case map[string]interface{}:
			// 简单粗暴，k = v ， 循环就完事了
			// 例子：map[string]interface{}{"modifie_time": "1657019178", "product_id": "1"}
			if vv3, ok := v1.(map[string]interface{}); ok {
				for k3, v3 := range vv3 {
					r.Where(k3, v3)
				}
			}

		}
	}

	return r
}

// Order ORM 添加排序条件
func (r *MysqlG) Order(order string) *MysqlG {
	r.OrderS = order
	return r
}

// Group ORM 设置分组条件
func (r *MysqlG) Group(group string) *MysqlG {
	r.GroupS = group
	return r
}

// SelectAssemble ORM 组装sql语句
func (r *MysqlG) SelectAssemble() *MysqlG {
	if Empty(r.Table) {
		//直接让OrmSql滞空，这样就get不出来数据了
		return r
	}

	// 设置默认查询对象
	ss := "*"
	if !Empty(r.VisibleS) {
		ss = strings.Replace(JsonEncode(r.VisibleS), "[", "", -1)
		ss = strings.Replace(ss, "]", "", -1)
		ss = strings.Replace(ss, "\"", "", -1)
		ss = strings.Replace(ss, ",", "`,`", -1)
		ss = "`" + ss + "`"
	}

	where := ""
	if !Empty(r.Wheres) {
		where = strings.Replace(r.Wheres, "AND", "WHERE ", 1)
	}

	r.OrmSql = "SELECT " + ss + " FROM " + r.GraveAccent(r.Table) + where + " " + r.GroupS + " " + r.OrderS

	return r
}

// GetSelectOrmSql 获取最终sql语句
func (r *MysqlG) GetSelectOrmSql() string {
	return r.SelectAssemble().OrmSql
}

// Get ORM获取单条数据
func (r *MysqlG) Get(where ...interface{}) map[string]interface{} {
	if len(where) > 0 {
		return r.GetOne(r.WhereBatch(where).GetSelectOrmSql())
	} else {
		return r.GetOne(r.GetSelectOrmSql())
	}
}

// All ORM获取所有数据
func (r *MysqlG) All(where ...interface{}) (rows []map[string]interface{}) {
	if len(where) > 0 {
		return r.WhereBatch(where).GetAll(r.GetSelectOrmSql())
	} else {
		return r.GetAll(r.GetSelectOrmSql())
	}
}

// ========================= 上面是有关于ORM的尝试操作 =========================//
