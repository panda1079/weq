package main // 声明此文件为可执行程序

import (
	"log"      // 日志记录
	"net/http" // 处理 HTTP 请求

	"github.com/gorilla/websocket" // 使用 Gorilla WebSocket 库来处理 WebSocket 连接
)

// 定义聊天室类型
type chatRoom struct {
	clients    map[*websocket.Conn]bool // 客户端连接映射表，用于存储客户端连接实例和其在线状态
	broadcast  chan []byte              // 广播消息通道，用于向所有客户端广播消息
	register   chan *websocket.Conn     // 连接注册通道，用于将新连接注册到聊天室中
	unregister chan *websocket.Conn     // 连接注销通道，用于从聊天室中注销指定连接
}

// HandleWebSocket 处理 WebSocket 连接
func (room *chatRoom) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 设置跨域请求头
	// 将 HTTP 请求升级为 WebSocket 连接
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}

	// 注册连接到聊天室中
	room.register <- conn
	defer func() {
		room.unregister <- conn
	}()

	// 循环读取客户端发送的消息并将其广播到所有连接的客户端
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message from client: %v", err)
			break
		}
		room.broadcast <- message
	}
}

// Run 启动聊天室，处理连接注册、注销、消息广播等逻辑
func (room *chatRoom) Run() {
	for {
		select {
		case conn := <-room.register:
			// 将新连接注册到映射表中，并打印连接信息
			room.clients[conn] = true
			log.Printf("Client %v joined the chat", conn.RemoteAddr())
		case conn := <-room.unregister:
			// 从映射表中注销指定连接，并打印连接信息
			if _, ok := room.clients[conn]; ok {
				delete(room.clients, conn)
				log.Printf("Client %v left the chat", conn.RemoteAddr())
			}
		case message := <-room.broadcast:
			// 广播消息到所有连接的客户端
			for conn := range room.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Failed to send message to client %v: %v", conn.RemoteAddr(), err)
					delete(room.clients, conn)
				}
			}
		}
	}
}

func main() {
	// 创建聊天室对象
	room := &chatRoom{
		clients:    make(map[*websocket.Conn]bool), // 客户端连接映射表
		broadcast:  make(chan []byte),              // 广播消息通道
		register:   make(chan *websocket.Conn),     // 连接注册通道
		unregister: make(chan *websocket.Conn),     // 连接注销通道
	}

	// 启动聊天室
	go room.Run()

	// 处理 WebSocket 连接请求
	http.HandleFunc("/ws", room.HandleWebSocket)

	// 监听端口并启动 HTTP 服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
