package wsocket

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/gorilla/websocket"
)

// Message 消息
type Message struct {
	Name    string          `json:"name"`    // 用户名称
	Message ResponseMessage `json:"message"` // 消息内容
}

// ResponseMessage 响应消息
type ResponseMessage struct {
	Type       string      `json:"type"`       // 消息类型: error | bim2gltf | cad | chatroom-join | chatroom-leave
	Subscriber string      `json:"subscriber"` // 订阅者
	Data       interface{} `json:"data"`       // 消息内容
}

var (
	// 用于新连接用户的通道.
	subscribe = make(chan *Connection, 10)
	// 用户退出通道.
	unsubscribe = make(chan string, 10)
	// 广播消息通道.
	publish = make(chan Message, 10)
	// 订阅者列表
	subscribers = list.New()
)

func Join(conn *Connection) {
	subscribe <- conn
}

func Leave(user string) {
	unsubscribe <- user
}

// Publish 发消息
func Publish(message Message) {
	publish <- message
}

// GetConnByUName 获取用户连接
func GetConnByUName(uname string) *Connection {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(*Connection).uname == uname {
			return sub.Value.(*Connection)
		}
	}
	return nil
}

// 判断用户是否已存在
func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(*Connection).uname == user {
			return true
		}
	}
	return false
}

// broadcastWebSocket 向WebSocket用户广播消息
func broadcastWebSocket(message Message) {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// 立即发送事件给WebSocket用户
		ws := sub.Value.(*Connection).wsConn
		if ws != nil {
			if message.Message.Subscriber == "" {
				message.Message.Subscriber = message.Name
			}
			msgJsonByte, _ := json.Marshal(message.Message)
			if err := ws.WriteMessage(websocket.TextMessage, msgJsonByte); err != nil {
				fmt.Println("ws:broadcastWebSocket error:", err)
				// 用户退出连接池.
				unsubscribe <- sub.Value.(*Connection).uname
			}
		}
	}
}

// 处理所有传入的chan消息
func chatroom() {
	for {
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.uname) {
				//将用户添加到列表末尾
				subscribers.PushBack(sub)

				// 广播用户加入信息.
				var msgJson = ResponseMessage{
					Type:       "chatroom-join",
					Subscriber: sub.uname,
					Data:       sub.uname + " has joined the room.",
				}

				publish <- Message{Name: sub.uname, Message: msgJson}

				logs.Info("User %s join the room.", sub.uname)
			}
		case message := <-publish:
			broadcastWebSocket(message)
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(*Connection).uname == unsub {
					subscribers.Remove(sub)
					if sub.Value.(*Connection).wsConn != nil {
						sub.Value.(*Connection).Close()
					}
					// 广播用户退出信息.
					var msgJson = ResponseMessage{
						Type:       "chatroom-leave",
						Subscriber: unsub,
						Data:       unsub + " exit the room.",
					}
					publish <- Message{Name: unsub, Message: msgJson}
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}
