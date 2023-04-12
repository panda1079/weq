package library

import (
	"github.com/gorilla/websocket"
)

// ServerS 关于服务传递的公共函数（不能进去之后继承，真是让人头秃的设置）
type ServerS struct {
	MDb MysqlG
	RDb RedisG
	WSk map[string]*WebSocket //WebSocket模块组 （先预设好，这样就能一直放着）
}

// InitServerS 加载配置，装载链接
func (r *ServerS) InitServerS() {

	//加载Mysql模块
	r.MDb = MysqlG{}
	r.MDb.InitMysql()

	//加载Redis模块
	r.RDb = RedisG{}
	r.RDb.InitRedis()

	//预制WebSocket模块组空间
	if r.WSk == nil {
		r.WSk = make(map[string]*WebSocket)
	}

}

// RunWsk 添加WebSocket元素
func (r *ServerS) RunWsk(key string, CH HttpInfo) (bool, *WebSocket, *websocket.Conn) {
	//定义参数
	time := Time()

	if _, ok := r.WSk[key]; !ok {
		// 创建聊天室对象
		room := &WebSocket{
			key:      key,
			addTime:  time,
			lastTime: time,

			clients:           make(map[*websocket.Conn]bool),     // 客户端连接映射表
			Broadcast:         make(chan []byte),                  // 广播消息通道
			Register:          make(chan *websocket.Conn),         // 连接注册通道
			Unregister:        make(chan *websocket.Conn),         // 连接注销通道
			SpecificBroadcast: make(map[string][]*websocket.Conn), // 连接注销通道
		}

		//协程拉起
		go room.Run()

		return room.HandleWebSocket(CH)
	} else {
		r.WSk[key].lastTime = time
		return r.WSk[key].HandleWebSocket(CH)
	}
}
