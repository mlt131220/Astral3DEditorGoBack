# Astral3DEditorGoBack

🌍
*[简体中文](README.md)*

基于`Golang 1.20.3`和`Beego 2.0.0`的 [Astral3DEditor](https://github.com/mlt131220/Astral3DEditor) 项目后端代码.

![Static Badge](https://img.shields.io/badge/go-1.20.3-green)
![Static Badge](https://img.shields.io/badge/beego-2.0.0-8732D7)
![Static Badge](https://img.shields.io/badge/license-MIT-blue)

## 快速开始
```shell
   git clone https://github.com/mlt131220/Astral3DEditorGoBack
```
Tips:
* 数据库使用 MySQL，表结构文件位于 `static/sql/astral-3d-editor.sql`；
* 修改 `conf/app.conf` 下 `sql::conn` 项为自己的数据库连接，格式为`用户名:密码@tcp(地址:端口)/数据库?charset=utf8&loc=Local`；
* 修改 `conf/app.conf` 下 `dev::cadDwgConverterAbPath` 项为本地CAD Dwg转换器执行程序文件夹绝对路径，转换程序使用 libreDWG(已包含在项目static/lib/libredwg文件夹)；
* `conf/app.conf` 下其它配置项(Revit转换服务配置/又拍云配置...)按需求改动；
* 如果开发环境下运行接口均报错404，请运行 `bee generate routers` 重新生成注解路由文件；

## 运行
```
    bee run -downdoc=true -gendoc=true
```
Tips:
* [bee 工具安装](https://beego.gocn.vip/beego/zh/developing/bee/#bee-%E5%B7%A5%E5%85%B7%E7%9A%84%E5%AE%89%E8%A3%85)
* `-downdoc=true` 启用swagger文件自动下载；`-gendoc=true` 启用自动生成文档;

## 打包
```shell
    # Linux
    bee pack -be GOOS=linux -be GOARCH=amd64
    # Windows
    bee pack -be GOOS=windows
```
Tips:
* 修改`conf/app.conf` -> `runmode = prod`;

## 感谢🌹🌹🌹
如果本项目帮助到了你，请在[这里](https://github.com/mlt131220/Astral3DEditorGoBack/issues/1)留下你的网址，让更多的人看到。您的回复将会是我继续更新维护下去的动力。

## Star History
[![Star History Chart](https://api.star-history.com/svg?repos=mlt131220/Astral3DEditorGoBack&type=Date)](https://star-history.com/#mlt131220/Astral3DEditorGoBack&Date)