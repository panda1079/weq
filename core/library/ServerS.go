package library

// ServerS 关于服务传递的公共函数（不能进去之后继承，真是让人头秃的设置）
type ServerS struct {
	MDb MysqlG
	RDb RedisG
}

// InitServerS 加载配置，装载链接
func (r *ServerS) InitServerS() {

	//加载Mysql模块
	r.MDb = MysqlG{}
	r.MDb.InitMysql()

	//加载Redis模块
	r.RDb = RedisG{}
	r.RDb.InitRedis()
}
