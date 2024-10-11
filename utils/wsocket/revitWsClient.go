// Package wsocket
/**
 * @author ErSan
 * @email  mlt131220@163.com
 * @date   2024/4/21 17:12
 * @description Revit websocket客户端
 */
package wsocket

import (
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/gorilla/websocket"
	"net/url"
)

type WebsocketClientManager struct {
	conn        *websocket.Conn
	addr        *string
	path        string
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
}

// RevitWsClient 创建websocket客户端
func RevitWsClient(addrIp, addrPort, path string, timeout int) *WebsocketClientManager {
	addrString := addrIp + ":" + addrPort
	var sendChan = make(chan string, 10)
	var recvChan = make(chan string, 10)
	var conn *websocket.Conn
	return &WebsocketClientManager{
		addr:        &addrString,
		path:        path,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
	}
}

// dail 链接服务端
func (wsc *WebsocketClientManager) dail() {
	var err error
	u := url.URL{Scheme: "ws", Host: *wsc.addr, Path: wsc.path}
	fmt.Println("[revit ws] connecting to revit:", u.String())
	logs.Info("[revit ws] connecting to revit:", u.String())
	wsc.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	wsc.isAlive = true
	fmt.Println("[revit ws] connecting successful！！！")
	logs.Info("[revit ws] connecting successful！！！")
}

// SendMsg 发送消息
func (wsc *WebsocketClientManager) SendMsg(msg string) {
	wsc.sendMsgChan <- msg
}

// ReadMsg 接收消息
func (wsc *WebsocketClientManager) ReadMsg() (data string) {
	data = <-wsc.recvMsgChan
	return
}

// sendMsgThread 发送消息
func (wsc *WebsocketClientManager) sendMsgThread() {
	go func() {
		for {
			msg := <-wsc.sendMsgChan
			err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("[revit ws] write error:", err)
				logs.Error("[revit ws] write error:", err)
				continue
			}
			fmt.Println("[revit ws] send:", msg)
			logs.Notice("[revit ws] send:", msg)
		}
	}()
}

// readMsgThread 读取消息
func (wsc *WebsocketClientManager) readMsgThread() {
	go func() {
		for {
			if wsc.conn != nil {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					fmt.Println("[revit ws] read error:", err)
					logs.Error("[revit ws] read error:", err)
					wsc.isAlive = false
					// 出现错误，退出读取，尝试重连
					break
				}
				fmt.Println("[revit ws] receive:", message)
				logs.Notice("[revit ws] receive:", message)
				// 需要读取数据，不然会阻塞
				wsc.recvMsgChan <- string(message)
			}
		}
	}()
}

// Start 开启服务并重连
func (wsc *WebsocketClientManager) Start() {
	if wsc.isAlive == false {
		wsc.dail()
		wsc.sendMsgThread()
		wsc.readMsgThread()
	}
}

// IsConnected 判断是否连接
func (wsc *WebsocketClientManager) IsConnected() bool {
	return wsc.isAlive
}
