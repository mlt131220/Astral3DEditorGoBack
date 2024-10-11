package wsocket

import (
	json2 "encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	// 用户标识符
	uname  string
	wsConn *websocket.Conn
	//读取websocket的channel
	inChan    chan []byte
	closeChan chan byte
	mutex     sync.Mutex
	//closeChan 状态
	isClosed bool
}

// InitConnection 初始化长连接
func InitConnection(uname string, wsConn *websocket.Conn) {
	conn := &Connection{
		uname:     uname,
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	// 加入到连接池
	Join(conn)
	// TODO 当前有一个问题，会立马执行Leave，导致连接池中没有连接，所以无法广播消息；但是如果不执行Leave，会导致连接池中有很多无效连接（当前在广播分发时判断，未连接则使其退出）
	//defer Leave(uname)

	//启动读协程
	go conn.readLoop()
}

// ReadMessage 读取websocket消息
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("ws:connection is closed")
	}
	return
}

// Close 关闭连接
func (conn *Connection) Close() {
	//线程安全的Close,可重入
	conn.wsConn.Close()

	//只执行一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// readLoop 读协程
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}

		var json ResponseMessage
		json2.Unmarshal(data, &json)

		// 前端刷新页面可能离开订阅池，这里需要重新加入订阅池
		if !isUserExist(subscribers, conn.uname) {
			Join(conn)
		}

		// 加入消息广播通道
		publish <- Message{Name: conn.uname, Message: json}

		//如果数据量过大阻塞在这里,等待inChan有空闲的位置！
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			//closeChan关闭的时候
			goto ERR
		}
	}
ERR:
	conn.Close()
}
