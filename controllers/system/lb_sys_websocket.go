package system

import (
	"es-3d-editor-go-back/server"
	"es-3d-editor-go-back/utils/wsocket"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
	"net/http"
)

// LbSysWebSocketController 通用-WebSocket
type LbSysWebSocketController struct {
	beego.Controller
}

// URLMapping ...
func (c *LbSysWebSocketController) URLMapping() {
	c.Mapping("WsHandler", c.WsHandler)
}

var (
	upgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// WsHandler ...
// @Title WebSocket
// @Description WebSocket
// @router /ws [get]
func (c *LbSysWebSocketController) WsHandler() {
	// ws 连接用户标识
	uname := c.GetString("uname")

	var (
		//websocket 长连接
		ws  *websocket.Conn
		err error
	)

	//升级协议
	ws, err = upgrade.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		c.Data["json"] = server.RequestFail("不是websocket握手")
		c.ServeJSON()
		return
	} else if err != nil {
		c.Data["json"] = server.RequestFail("Cannot setup WebSocket connection")
		c.ServeJSON()
		return
	}

	wsocket.InitConnection(uname, ws)

	// 发送连接成功消息
	wsocket.Publish(wsocket.Message{
		Name: uname,
		Message: wsocket.ResponseMessage{
			Type:       "connect",
			Subscriber: uname,
			Data:       "hello " + uname,
		},
	})
}
