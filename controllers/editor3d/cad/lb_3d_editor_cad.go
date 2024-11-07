package cad

import (
	"encoding/json"
	"errors"
	"es-3d-editor-go-back/controllers"
	"es-3d-editor-go-back/controllers/system"
	"es-3d-editor-go-back/models/editor3d/cad"
	"es-3d-editor-go-back/server"
	"es-3d-editor-go-back/utils/wsocket"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// Lb3dEditorCadController operations for Lb3dEditorCad
type Lb3dEditorCadController struct {
	controllers.BaseController
}

// URLMapping ...
func (c *Lb3dEditorCadController) URLMapping() {
	c.Mapping("Dwg2Dxf", c.Dwg2Dxf)
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// ConversionResult cad转换结果
type ConversionResult struct {
	ConversionStatus string            `json:"conversionStatus"`
	Item             cad.Lb3dEditorCad `json:"item"`
	Message          string            `json:"message"`
}

// Dwg2Dxf
// @Title dwg to dxf
// @Description dwg转dxf
// @Param file formData file true "dwg文件"
// @Success 200 {string} outputPath
// @Failure 400 {string} error
// @router /cad/dwg2dxf [post]
func (c *Lb3dEditorCadController) Dwg2Dxf() {
	// 获取上传的文件
	file, header, err := c.GetFile("file")
	if err != nil {
		c.Data["json"] = server.RequestFail("获取上传文件失败！")

		c.ResponseJson()
		return
	}
	defer file.Close()

	conversionStatus, _ := c.GetInt("conversionStatus")
	converterFilePath := ""
	var dataPath string
	// 判断上传的文件是否需要转换（需转换格式：[.dwg,...]）
	var ext = strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])

	if ext != "dwg" && ext != "dxf" {
		c.Data["json"] = server.RequestFail("上传文件格式不正确！")

		c.ResponseJson()
		return
	}

	//20241108:文件名去除空格，不然cmd指令会报错
	header.Filename = strings.Replace(header.Filename, " ", "", -1)

	if ext != "dwg" {
		conversionStatus = 1
		converterFilePath = "/static/upload/cad/" + time.Now().Format("20060102") + "/" + header.Filename
		dataPath = converterFilePath

		// 无需解析，直接上传至upyun
		err = system.UpYunUpload(converterFilePath, file)
		if err != nil {
			c.Data["json"] = server.RequestFail(header.Filename + "上传至又拍云失败! error:" + err.Error())
			logs.Error(header.Filename+"上传至又拍云失败，Error: ", err.Error())

			c.ResponseJson()
			return
		}
	} else {
		conversionStatus = 0

		// 需要解析，保存文件到tmp目录
		dataPath, _ = beego.AppConfig.String("temporaryFolder")
		dataPath += header.Filename
		err = c.SaveToFile("file", dataPath)
		if err != nil {
			c.Data["json"] = server.RequestFail("保存文件失败！")
			logs.Error("cad上传文件保存失败，Error: ", err.Error())

			c.ResponseJson()
			return
		}
	}

	// 更新数据库
	var v = cad.Lb3dEditorCad{
		FilePath:          dataPath,
		Thumbnail:         c.GetString("thumbnail"),
		FileName:          c.GetString("fileName"),
		ConversionStatus:  conversionStatus,
		ConverterFilePath: converterFilePath,
	}
	if _, err := cad.AddLb3dEditorCad(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = server.RequestSuccess(v)

		// 转换dwg文件
		if conversionStatus == 0 {
			go func() {
				// 获取当前用户的ws连接
				var uName = c.GetString("uname")

				//拼接执行程序及命令
				var exePath, _ = beego.AppConfig.String("cadDwgConverterAbPath")
				var outputPath, _ = beego.AppConfig.String("temporaryFolder")

				// 拼接输出文件名
				var outputFile = outputPath + "dxf/" + strings.Replace(header.Filename, ".dwg", ".dxf", -1)

				// exePath + "dwg2dxf.exe " + datapath + " -o " + outputFile
				var builder strings.Builder
				builder.WriteString(exePath)
				builder.WriteString("dwg2dxf.exe ")
				builder.WriteString(dataPath)
				builder.WriteString(" -o ")
				builder.WriteString(outputFile)
				var command = builder.String()

				fmt.Println("[cad] run cmd：", command)
				logs.Notice("[cad] run cmd：", command)
				cmd := exec.Command("cmd", "/C", command)
				if err = cmd.Run(); err != nil {
					logs.Error("failed to call cmd.Run(): ", err.Error())

					v.ConversionStatus = 2

					wsocket.Publish(wsocket.Message{
						Name: uName,
						Message: wsocket.ResponseMessage{
							Type:       "cad",
							Subscriber: uName,
							Data: ConversionResult{
								ConversionStatus: "failed",
								Item:             v,
								Message:          err.Error(),
							},
						},
					})
				} else {
					// 读取转换后的文件流，上传至upyun
					dxfFile, err := os.Open(outputFile)
					defer dxfFile.Close()
					if err != nil {
						logs.Error("failed to open file: ", err.Error())

						v.ConversionStatus = 2

						wsocket.Publish(wsocket.Message{
							Name: uName,
							Message: wsocket.ResponseMessage{
								Type:       "cad",
								Subscriber: uName,
								Data: ConversionResult{
									ConversionStatus: "failed",
									Item:             v,
									Message:          "转换成功，但读取结果失败！Error: " + err.Error(),
								},
							},
						})
					} else {
						v.ConverterFilePath = "/static/upload/cad/" + time.Now().Format("20060102") + "/" + filepath.Base(dxfFile.Name())
						fmt.Println("[cad] 转换成功，准备上传至又拍云: ", dxfFile.Name())
						err = system.UpYunUpload(v.ConverterFilePath, dxfFile)
						if err != nil {
							logs.Error(header.Filename+"上传至又拍云失败，Error: ", err.Error())

							v.ConversionStatus = 2

							wsocket.Publish(wsocket.Message{
								Name: uName,
								Message: wsocket.ResponseMessage{
									Type:       "cad",
									Subscriber: uName,
									Data: ConversionResult{
										ConversionStatus: "failed",
										Item:             v,
										Message:          "转换成功，但上传结果至又拍云失败！Error: " + err.Error(),
									},
								},
							})
						} else {
							// 转换-打开-上传成功
							v.ConversionStatus = 1
							v.FilePath = v.ConverterFilePath
						}
					}
				}

				fmt.Println("[cad] 转换完成，准备更新数据库: ", v.Id)
				err := cad.UpdateLb3dEditorCadById(&v)
				if v.ConversionStatus == 1 {
					if err == nil {
						fmt.Println("[cad] 转换成功，更新数据库成功！发送成功消息")
						wsocket.Publish(wsocket.Message{
							Name: uName,
							Message: wsocket.ResponseMessage{
								Type:       "cad",
								Subscriber: uName,
								Data: ConversionResult{
									ConversionStatus: "completed",
									Item:             v,
									Message:          "转换成功！",
								},
							},
						})
					} else {
						wsocket.Publish(wsocket.Message{
							Name: uName,
							Message: wsocket.ResponseMessage{
								Type:       "cad",
								Subscriber: uName,
								Data: ConversionResult{
									ConversionStatus: "failed",
									Item:             v,
									Message:          "转换成功，但更新数据库失败！Error: " + err.Error(),
								},
							},
						})
					}
				}

				cmd.ProcessState.ExitCode()

				// 删除上传的临时文件和转换后的文件
				os.Remove(dataPath)
				os.Remove(outputFile)
			}()
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}

	c.ResponseJson()
}

// Post ...
// @Title 新增cad文件数据
// @Description create Lb3dEditorCad
// @Param	body	models.Lb3dEditorCad	true		"body for Lb3dEditorCad content"
// @Success 201 {object} models.Lb3dEditorCad
// @Failure 403 body is empty
// @router /cad/add [post]
func (c *Lb3dEditorCadController) Post() {
	var v cad.Lb3dEditorCad
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := cad.AddLb3dEditorCad(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = server.RequestSuccess(v)
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}

// GetOne ...
// @Title 获取一条数据
// @Description get Lb3dEditorCad by id
// @Param	id		path 	string	true		"获取一条数据"
// @Success 200 {object} models.Lb3dEditorCad
// @Failure 403 :id is empty
// @router /cad/:id [get]
func (c *Lb3dEditorCadController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := cad.GetLb3dEditorCadById(id)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess(v)
	}
	c.ResponseJson()
}

// GetAll ...
// @Title Get All
// @Description 按分页信息获取全部
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Lb3dEditorCad
// @Failure 403
// @router /cad/getAll [get]
func (c *Lb3dEditorCadController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = server.RequestFail(errors.New("error:无效的查询键/值对").Error())
				c.ResponseJson()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := cad.GetAllLb3dEditorCad(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		res := make(map[string]interface{})
		res["items"] = l
		res["current"] = offset + 1
		res["pageSize"] = limit
		if total, totalErr := cad.GetTotalLb3dEditorCad(query); totalErr == nil {
			res["pages"] = math.Ceil(float64(total) / float64(limit))
			res["total"] = total
		}

		c.Data["json"] = server.RequestSuccess(res)
	}
	c.ResponseJson()
}

// Put ...
// @Title 更新数据
// @Description update the Lb3dEditorCad
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Lb3dEditorCad	true		"body for Lb3dEditorCad content"
// @Success 200 {string} "OK"
// @Failure 403 :id is not int
// @router /cad/:id [put]
func (c *Lb3dEditorCadController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := cad.Lb3dEditorCad{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := cad.UpdateLb3dEditorCadById(&v); err == nil {
			c.Data["json"] = server.RequestSuccess("OK")
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}

// Delete ...
// @Title Delete
// @Description 通过id删除
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /cad/:id [delete]
func (c *Lb3dEditorCadController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := cad.DeleteLb3dEditorCad(id); err == nil {
		c.Data["json"] = server.RequestSuccess("delete success!")
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}
