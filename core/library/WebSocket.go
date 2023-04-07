package library

import (
	"bytes"
	"github.com/gorilla/websocket"
	"net/http"
)

// WebSocket 关于处理 WebSocket 连接服务的公共函数
type WebSocket struct {
	key      string //key
	addTime  int64  //添加时间
	lastTime int64  //最后使用时间

	clients           map[*websocket.Conn]bool     // 客户端连接映射表
	Broadcast         chan []byte                  // 广播消息通道
	Register          chan *websocket.Conn         // 连接注册通道
	Unregister        chan *websocket.Conn         // 连接注销通道
	SpecificBroadcast map[string][]*websocket.Conn // 指定广播对象连接列表
}

// HandleWebSocket 处理升级 WebSocket 连接
func (room *WebSocket) HandleWebSocket(CH HttpInfo) (bool, *WebSocket, *websocket.Conn) {

	// 检查错误是否为已关闭的连接。如果是，则打印一条日志并立即中止读取循环。
	if !websocket.IsWebSocketUpgrade(CH.Request) {
		http.Error(CH.ResponseWriter, "Expected WebSocket Upgrade", http.StatusBadRequest)
		return false, room, new(websocket.Conn)
	}

	// 将 HTTP 请求转为 WebSocket 连接
	//conn, err := websocket.Accept(CH.ResponseWriter, CH.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
	conn, err := websocket.Upgrade(CH.ResponseWriter, CH.Request, nil, 1024, 1024)
	if err != nil {
		SetLog(err, "转为 WebSocket 连接失败")
		return false, room, conn
	}

	//在中止循环之前需要调用defer room.unregister <- conn将该连接注销掉。
	defer func() {
		room.Unregister <- conn
	}()

	return true, room, conn
}

// Run 启动连接
func (room *WebSocket) Run() {
	for {
		select { // 监听事件通道,阻塞主程序
		//select{} 是一个空的 select 语句，它会一直阻塞主程序，直到有可用的 case。如果所有的 case 都没有准备好（或者没有 case），那么该 select 语句就会一直阻塞等待，不会导致 CPU 占用率增加。
		//在实现中，select 语句被用于监听三个 channel 的操作：room.register、room.unregister 和 room.broadcast。只有这些 channel 中有数据传入时，相应的 case 才会被执行，否则程序将一直阻塞等待。
		case conn := <-room.Register:
			// 将新连接注册到映射表中，并打印连接信息
			room.clients[conn] = true
			SetLog(conn.RemoteAddr(), "有新连接注册到映射表中")
		case conn := <-room.Unregister:
			// 从映射表中注销指定连接，并打印连接信息
			if _, ok := room.clients[conn]; ok {
				delete(room.clients, conn)
				SetLog(conn.RemoteAddr(), "新连接从映射表中注销")
			}
		case message := <-room.Broadcast:
			// 广播消息到所有连接的客户端
			for conn := range room.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					SetLog(conn.RemoteAddr().String()+":"+err.Error(), "无法向客户端发送消息")
					delete(room.clients, conn)
				}
			}
		}
	}
}

// parseTarget 解析广播消息中的指定对象
func (room *WebSocket) parseTarget(message []byte) (string, bool) {
	if len(message) > 0 && message[0] == '@' {
		if i := bytes.IndexByte(message, ' '); i > 1 {
			return string(message[1:i]), true
		}
	}
	return "", false
}

// Airing 广播内容
func (room *WebSocket) Airing(CH HttpInfo, Clawback func(message []byte) []byte) *WebSocket {
	// 循环读取客户端发送的消息并将其广播到所有连接的客户端
	for {
		_, message, err := CH.ThisConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				SetLog(err.Error(), "连接已关闭")
			} else {
				SetLog(err.Error(), "无法从客户端读取消息")
			}
			break
		}

		if len(message) == 0 { // 忽略空消息
			continue
		}

		// 如果指定了广播对象，则只向该对象中的连接发送消息
		if target, ok := room.parseTarget(message); ok {
			if cones, exists := room.SpecificBroadcast[target]; exists {
				for _, c := range cones {
					c.WriteMessage(websocket.TextMessage, Clawback(message))
				}
			}
		} else {
			// 否则把mod的内容返回给前端
			room.Broadcast <- Clawback(message)
		}
	}

	return room
}
