package system

import (
	"es-3d-editor-go-back/server"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/upyun/go-sdk/v3/upyun"
	"io"
	"sync"
	"time"
)

// LbSysUpYunController 通用-又拍云
type LbSysUpYunController struct {
	beego.Controller
}

// URLMapping ...
func (c *LbSysUpYunController) URLMapping() {
	c.Mapping("DoUpload", c.DoUpload)
	c.Mapping("DOBlob", c.DOBlob)
	c.Mapping("DoRemove", c.DoRemove)
}

var up *upyun.UpYun

// init 初始化UpYun-服务
func init() {
	var bucket, _ = beego.AppConfig.String("upyun::bucket")
	var operator, _ = beego.AppConfig.String("upyun::operator")
	var password, _ = beego.AppConfig.String("upyun::password")
	up = upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   bucket,
		Operator: operator,
		Password: password,
	})
}

// UpYunUpload 文件上传至又拍云
func UpYunUpload(savePath string, file io.ReadCloser) error {
	defer file.Close()

	err := up.Put(&upyun.PutObjectConfig{
		Path:   savePath,
		Reader: file,
	})
	if err != nil {
		return err
	}
	return nil
}

// DoUpload ...
// @Title 文件上传
// @Description 文件上传至又拍云
// @Param	biz	 body 	string	true	"文件上传的业务路径"
// @Success 200 {string} filePath
// @Failure 403 biz is empty
// @router /upyun/upload [post]
func (c *LbSysUpYunController) DoUpload() {
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

	//存储的文件夹
	saveFolder := biz + "/" + currentDate
	//存储的路径
	savePath := saveFolder + "/" + fileHeader.Filename

	/* 保存在 static/upload */
	err = up.Put(&upyun.PutObjectConfig{
		Path:   "/static/upload/" + savePath,
		Reader: file.(io.Reader),
	})

	if err != nil {
		c.Data["json"] = server.RequestFail("上传文件保存失败！error=" + err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess("static/upload/" + savePath)
	}
	c.ServeJSON()
}

// UpYunGetBlob 获取文件流
func UpYunGetBlob(ctx *context.Context, filePath string) (err error) {
	fileInfo, err := up.Get(&upyun.GetObjectConfig{
		Path:   filePath,
		Writer: ctx.ResponseWriter,
	})
	if err != nil {
		return err
	}

	ctx.Output.Header("Content-Type", "application/octet-stream;charset=utf-8")
	ctx.Output.Header("Content-Disposition", "attachment; filename="+fileInfo.Name)
	return nil
}

// DOBlob ...
// @Title 获取文件流
// @Description 获取文件流
// @Param	filePath	 query 	string	true 	"文件路径"
// @Success 200 {stream} file stream
// @Failure 403 filePath is empty
// @router /upyun/blob [get]
func (c *LbSysUpYunController) DOBlob() {
	filePath := c.GetString("filepath")

	fileInfo, err := up.Get(&upyun.GetObjectConfig{
		Path:   filePath,
		Writer: c.Ctx.ResponseWriter,
	})
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
		c.ServeJSON()
	}

	c.Ctx.Output.Header("Content-Type", "application/octet-stream;charset=utf-8")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fileInfo.Name)
	c.ServeJSON()
}

// UpYunDelete 移除又拍云文件方法
func UpYunDelete(filePath string) (err error) {
	err = up.Delete(&upyun.DeleteObjectConfig{
		Path:  filePath, // 云存储中的路径
		Async: true,     // 是否异步删除
	})

	if err != nil {
		//删除失败
		return err
	}
	return nil
}

// UpYunDeleteAll 移除又拍云文件夹及其下所有文件方法
func UpYunDeleteAll(folderPath string) (err error) {
	// 列目录
	objsChan := make(chan *upyun.FileInfo, 100)
	go func() {
		err := up.List(&upyun.GetObjectsConfig{
			Path:           folderPath,
			ObjectsChan:    objsChan,
			MaxListObjects: 5000, // 最大列对象个数
			MaxListTries:   3,    // 列目录最大重试次数
			MaxListLevel:   3,    // 递归最大深度
		})
		if err != nil {
			return
		}
	}()

	var wg sync.WaitGroup

	// 控制最大并发数量
	ch := make(chan struct{}, 500)
	defer close(ch)

	for file := range objsChan {
		ch <- struct{}{}

		wg.Add(1)

		go func(fileInfo *upyun.FileInfo) {
			defer wg.Done()

			if fileInfo.IsDir {
				// 递归删除目录
				err = UpYunDeleteAll(folderPath + fileInfo.Name + "/")
				if err != nil {
					return
				}
			}

			// 删除文件
			err = UpYunDelete(folderPath + fileInfo.Name)
			if err != nil {
				return
			}
			<-ch
		}(file)
	}

	wg.Wait()

	// 删除目录
	err = UpYunDelete(folderPath)

	return nil
}

// DoRemove ...
// @Title 删除文件接口
// @Description 删除upyun下的文件
// @Param	filePath	 body 	string	true	"文件所在的业务路径"
// @Success 200 {string} "文件已删除！"
// @Failure 403 "文件删除失败"
// @router /upyun/remove [post]
func (c *LbSysUpYunController) DoRemove() {
	filePath := c.GetString("biz")

	err := UpYunDelete(filePath)

	if err != nil {
		//删除失败
		c.Data["json"] = server.RequestFail("文件删除失败！error=" + err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess("文件已删除！")
	}
	c.ServeJSON()
}
