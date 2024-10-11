package controllers

import (
	"es-3d-editor-go-back/server"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
	Data map[interface{}]server.ResultJson
}

// NestPreparer 约定：如果子controller 存在NestPrepare()方法，就实现了该接口，
// 用于子controller自身的通用逻辑
type NestPreparer interface {
	NestPrepare()
}

// Prepare 所有继承了BaseController的子controller 每次请求都会调用BaseController的Prepare()方法
func (c *BaseController) Prepare() {
	c.Data = make(map[interface{}]server.ResultJson)

	// 判断子类是否实现了NestPreparer接口，如果实现了就调用接口方法。
	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

// ResponseJson 请求的响应 (比c.ServeJSON()多了SetStatus)
func (c *BaseController) ResponseJson(encoding ...bool) error {
	c.Ctx.Output.SetStatus(c.Data["json"].Code)

	var (
		hasIndent   = beego.BConfig.RunMode != beego.PROD
		hasEncoding = len(encoding) > 0 && encoding[0]
	)

	return c.Ctx.Output.JSON(c.Data["json"], hasIndent, hasEncoding)
}
