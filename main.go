package main

import (
	_ "es-3d-editor-go-back/routers"
	"es-3d-editor-go-back/utils"
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.Debug = true

	if err1 := orm.RegisterDriver("mysql", orm.DRMySQL); err1 != nil {
		logs.Error(err1.Error())
	}

	var sqlConn, _ = beego.AppConfig.String("sql::conn")
	if err2 := orm.RegisterDataBase("default", "mysql", sqlConn); err2 != nil {
		logs.Error(err2.Error())
		panic(err2.Error())
	}

	// 注册各模块
	utils.InitLogger() //调用logger初始化
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//var appPath, _ = beego.AppConfig.String("currentAbPath")
	//fmt.Println("当前项目路径:", appPath)
	beego.Run()
}
