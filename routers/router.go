// Package routers
// @APIVersion 1.0.0
// @Title ES3DEditor API
// @Description ES3DEditor API
// @Contact mlt131220@163.com
// @TermsOfServiceUrl https://mhbdng.cn/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"es-3d-editor-go-back/controllers/editor3d/bim"
	"es-3d-editor-go-back/controllers/editor3d/cad"
	"es-3d-editor-go-back/controllers/editor3d/scenes"
	"es-3d-editor-go-back/controllers/system"
	"es-3d-editor-go-back/server"
	"es-3d-editor-go-back/utils/jwt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	ns := beego.NewNamespace("/api",
		//所有系统共用的系统级路由
		beego.NSNamespace("/sys",
			beego.NSInclude(
				&system.LbSysUserController{},
				&system.LbSysUploadController{},
				&system.LbSysWebSocketController{},
				&system.LbSysUpYunController{},
			),
		),
		//逆光三维编辑器路由
		beego.NSNamespace("/editor3d",
			beego.NSInclude(
				&bim.Lb3dEditorBimToGltfController{},
				&cad.Lb3dEditorCadController{},
				&scenes.Lb3dEditorScenesController{},
				&scenes.Lb3dEditorScenesExampleController{},
			),
		),
	)
	// 添加路由过滤，鉴权(暂无用，es3d暂时没有登录及权限控制)
	beego.InsertFilter("/api/admin/*", beego.BeforeRouter, func(ctx *context.Context) {
		token := ctx.Request.Header.Get("X-Token")
		if token == "" {
			OutJson(ctx, server.RequestNoPermission())
			return
		}
		_, err := jwt.ValidateToken(token)
		if err != nil {
			OutJson(ctx, server.RequestNoPermission())
			return
		}
	})
	beego.AddNamespace(ns)
}

// OutJson 这是输出给权限的
func OutJson(ctx *context.Context, OutData server.ResultJson) {
	ctx.Output.Status = 200
	ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	ctx.Output.Header("Access-Control-Allow-Headers", "x-token,X-Token")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
	ctx.Output.JSON(OutData, true, true)
}
