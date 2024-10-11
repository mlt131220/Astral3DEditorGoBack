package system

import (
	"es-3d-editor-go-back/server"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"os"
	"strconv"
	"strings"
	"time"
)

// LbSysUploadController 通用-文件上传和获取
type LbSysUploadController struct {
	beego.Controller
}

// URLMapping ...
func (c *LbSysUploadController) URLMapping() {
	c.Mapping("DoUpload", c.DoUpload)
	c.Mapping("DOBlob", c.DOBlob)
	c.Mapping("DoDownload", c.DoDownload)
	c.Mapping("DoRemoveFile", c.DoRemoveFile)
}

// DoUpload ...
// @Title 文件上传
// @Description 处理文件上传
// @Param	biz	 body 	string	true	"文件上传的业务路径"
// @Success 200 {string} filePath
// @Failure 403 biz is empty
// @router /upload [post]
func (c *LbSysUploadController) DoUpload() {
	file, fileHeader, err := c.GetFile("file")

	// 关闭上传的文件，不然的话会出现临时文件不能清除的情况
	defer file.Close()

	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
		c.ServeJSON()
		return
	}

	biz := c.GetString("biz")

	currentDate := time.Now().Format("20060102")
	var fileNameMap []string = strings.Split(fileHeader.Filename, ".")

	//存储的文件夹
	saveFolder := biz + "/" + currentDate
	//存储的路径
	savePath := saveFolder + "/" + fileNameMap[0] + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "." + fileNameMap[1]

	/* 保存位置在 static/upload, 没有文件夹要先创建 */
	//先判断创建目录
	if _, err = os.Stat("static/upload/" + saveFolder); os.IsNotExist(err) {
		err := os.MkdirAll("static/upload/"+saveFolder, os.ModePerm)
		if err != nil {
			c.Data["json"] = server.RequestFail("服务端创建文件夹" + biz + "失败！error=" + err.Error())
			c.ServeJSON()
			return
		}
	}

	err = c.SaveToFile("file", "static/upload/"+savePath)
	if err != nil {
		c.Data["json"] = server.RequestFail("服务器端文件保存失败！error=" + err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess("static/upload/" + savePath)
	}
	c.ServeJSON()
}

// DOBlob ...
// @Title 获取文件流
// @Description 获取文件流
// @Success 200 {stream} file stream
// @Param	filepath	query 	string	true		"文件路径"
// @Failure 403 file is empty
// @router /blob [get]
func (c *LbSysUploadController) DOBlob() {
	filePath := c.GetString("filepath")

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
		c.ServeJSON()
	}

	c.Ctx.Output.Header("Content-Type", "application/octet-stream;charset=utf-8")
	filePathArr := strings.Split(filePath, "/")
	filename := filePathArr[len(filePathArr)-1]
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+filename)
	err = c.Ctx.Output.Body(fileBytes)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	}

	c.ServeJSON()
}

// DoDownload ...
// @Title 下载文件
// @Description 下载文件
// @Param	filePath	 query 	string	true	"下载文件的业务路径"
// @Success 200 {string} filePath
// @Failure 403 filePath is empty
// @router /downloadFile [get]
func (c *LbSysUploadController) DoDownload() {
	v := c.GetString("filePath")

	fmt.Println(v)

	//第一个参数是文件的地址，第二个参数是下载显示的文件的名称
	c.Ctx.Output.Download("static/upload/"+v, strings.Split(v, "/")[2])
}

// DoRemoveFile ...
// @Title 删除文件
// @Description 删除static/upload文件夹下的文件
// @Param	filePath	 body 	string	true	"文件所在的业务路径"
// @Success 200 {string} "文件已删除！"
// @Failure 403 "文件删除失败"
// @router /removeFile [post]
func (c *LbSysUploadController) DoRemoveFile() {
	filePath := c.GetString("biz")
	err := os.Remove(filePath)
	if err != nil {
		//删除失败
		c.Data["json"] = server.RequestFail("文件删除失败！error=" + err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess("文件已删除！")
	}
	c.ServeJSON()
}
