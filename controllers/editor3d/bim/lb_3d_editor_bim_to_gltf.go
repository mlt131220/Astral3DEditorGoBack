package bim

import (
	"encoding/json"
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"es-3d-editor-go-back/controllers"
	"es-3d-editor-go-back/models/editor3d/bim"
	"es-3d-editor-go-back/server"
	"es-3d-editor-go-back/utils/wsocket"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// Lb3dEditorBimToGltfController BimToGltf bim轻量化
type Lb3dEditorBimToGltfController struct {
	controllers.BaseController
}

// RvtConversionRequestOptions rvt转换请求结构配置
type RvtConversionRequestOptions struct {
	UseDraco            bool   `json:"useDraco"`
	Optimize            bool   `json:"optimize"`
	ExportProperty      bool   `json:"exportProperty"`
	View                string `json:"view"`
	ViewName            string `json:"viewName"`
	DisplayStyle        string `json:"displayStyle"`
	CoordinateReference string `json:"coordinateReference"`
}

// RvtConversionRequest rvt转换请求结构
type RvtConversionRequest struct {
	FileId   int                         `json:"fileId"`
	FilePath string                      `json:"filePath"`
	Options  RvtConversionRequestOptions `json:"options"`
}

// RvtConversionResult rvt转换结果
type RvtConversionResult struct {
	ConversionStatus string                  `json:"conversionStatus"`
	Item             bim.Lb3dEditorBimToGltf `json:"item"`
	Process          float64                 `json:"process"`
}

// URLMapping ...
func (c *Lb3dEditorBimToGltfController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("AddAndConversion", c.AddAndConversion)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("DoRvtUpload", c.DoRvtUpload)
}

var revitWs *wsocket.WebsocketClientManager
var _appPath, _ = beego.AppConfig.String("currentAbPath")

func init() {
	var address, _ = beego.AppConfig.String("revit::address")
	var port, _ = beego.AppConfig.String("revit::port")
	var path, _ = beego.AppConfig.String("revit::path")
	var timeout, _ = beego.AppConfig.Int("revit::timeout")
	// 启动websocket连接revit,处理相关转换 (8小时自动重连一次)
	revitWs = wsocket.RevitWsClient(address, port, path, timeout)
	revitWs.Start()
}

// Post ...
// @Title 新增转换数据
// @Description create Lb3dEditorBimToGltf
// @Param	body 	bim.Lb3dEditorBimToGltf	true		"body for Lb3dEditorBimToGltf content"
// @Success 201 {object} bim.Lb3dEditorBimToGltf
// @Failure 403 body is empty
// @router /bim2gltf/add [post]
func (c *Lb3dEditorBimToGltfController) Post() {
	var v bim.Lb3dEditorBimToGltf
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := bim.AddLb3dEditorBimToGltf(&v); err == nil {
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

// AddAndConversion ...
// @Title 新增数据并启动socket连接revit转换
// @Description create Lb3dEditorBimToGltf and start socket connection
// @Param	body 	bim.Lb3dEditorBimToGltf	true		"body for Lb3dEditorBimToGltf content"
// @Param	query 	uname	true		"websocket unique identifier: uname"
// @Success 201 {object} bim.Lb3dEditorBimToGltf
// @Failure 403 body is empty
// @router /bim2gltf/addAndConversion [post]
func (c *Lb3dEditorBimToGltfController) AddAndConversion() {
	var op RvtConversionRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &op)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
		c.ResponseJson()
		return
	}

	var v bim.Lb3dEditorBimToGltf
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := bim.AddLb3dEditorBimToGltf(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = server.RequestSuccess(v)

			// 转换线程
			go func() {
				// 获取当前用户的ws连接
				var uName = c.GetString("uname")

				var bimFilePath = _appPath + v.BimFilePath

				if !revitWs.IsConnected() {
					wsocket.Publish(wsocket.Message{
						Name: uName,
						Message: wsocket.ResponseMessage{
							Type: "error",
							Data: "rvt 轻量化服务未连接",
						},
					})
					return
				} else {
					// 调用revit开始转换
					m := RvtConversionRequest{
						FileId:   v.Id,
						FilePath: bimFilePath,
						Options: RvtConversionRequestOptions{
							UseDraco:            true,
							Optimize:            op.Options.Optimize,
							ExportProperty:      op.Options.ExportProperty,
							View:                op.Options.View,
							ViewName:            op.Options.ViewName,
							DisplayStyle:        op.Options.DisplayStyle,
							CoordinateReference: op.Options.CoordinateReference,
						},
					}
					sm, err := json.Marshal(m)
					if err != nil {
						return
					}
					//revitWs.SendMsg(bimFilePath)
					revitWs.SendMsg(string(sm))

					// 读取revit消息
					for {
						msg := revitWs.ReadMsg()
						// 读取到消息是json字符串，转为map
						var msgMap map[string]interface{}
						if err := json.Unmarshal([]byte(msg), &msgMap); err == nil {
							if msgMap["bimFilePath"] != bimFilePath {
								continue
							}

							// 读取转换进度
							if msgMap["type"] == "progress" {
								// 读取到进度信息
								if msgMap["progress"] != nil {
									progress := msgMap["progress"].(float64)
									// 转换进度
									wsocket.Publish(wsocket.Message{
										Name: uName,
										Message: wsocket.ResponseMessage{
											Type: "bim2gltf",
											Data: RvtConversionResult{
												ConversionStatus: "progress",
												Item:             v,
												Process:          progress,
											},
										},
									})
								}
								continue
							} else if msgMap["type"] == "completed" { // 读取到转换完成信息
								var vc bim.Lb3dEditorBimToGltf
								vc.Id = msgMap["fileId"].(int)
								vc.FileName = v.FileName
								vc.GltfFilePath = msgMap["gltfPath"].(string)
								vc.ConversionDuration = msgMap["runSeconds"].(float64)
								vc.ConversionStatus = 1
								// 获取v.GltfFilePath对应文件的大小，存入v.GltfFileSize字段
								if fileInfo, err := os.Stat(_appPath + v.GltfFilePath); err == nil {
									v.GltfFileSize = float64(fileInfo.Size())
								}

								// TODO 上传至upyun

								// 返回给前端的json结构
								var resJson = wsocket.ResponseMessage{
									Type:       "bim2gltf",
									Subscriber: uName,
									Data: RvtConversionResult{
										ConversionStatus: "completed",
										Item:             v,
										Process:          100,
									},
								}

								// 更新数据库
								bim.UpdateLb3dEditorBimToGltfById(&vc)

								wsocket.Publish(wsocket.Message{
									Name:    uName,
									Message: resJson,
								})

								// 退出for循环，并退出线程
								break
							} else if msgMap["type"] == "failed" { // 读取到转换失败信息
								v.GltfFilePath = ""
								v.ConversionDuration = 0
								v.ConversionStatus = 2
								v.GltfFileSize = 0

								// 返回给前端的json结构
								var resJson = wsocket.ResponseMessage{
									Type:       "bim2gltf",
									Subscriber: uName,
									Data: RvtConversionResult{
										ConversionStatus: "failed",
										Item:             v,
									},
								}

								// 转换失败
								wsocket.Publish(wsocket.Message{
									Name:    uName,
									Message: resJson,
								})

								// 更新数据库
								// 也可以更改为删除rvt文件 并删除数据库数据，但需要对应更改前端页面
								bim.UpdateLb3dEditorBimToGltfById(&v)

								break
							}
						}
					}
				}
			}()
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}

// GetOne ...
// @Title 获取一条转换数据
// @Description get Lb3dEditorBimToGltf by id
// @Param	id		path 	string	true		"转换数据id"
// @Success 200 {object} bim.Lb3dEditorBimToGltf
// @Failure 403 :id is empty
// @router /bim2gltf/get/:id [get]
func (c *Lb3dEditorBimToGltfController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := bim.GetLb3dEditorBimToGltfById(id)
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
// @Success 200 {object} bim.Lb3dEditorBimToGltf
// @Failure 403
// @router /bim2gltf/getAll [get]
func (c *Lb3dEditorBimToGltfController) GetAll() {
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

	l, err := bim.GetAllLb3dEditorBimToGltf(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		res := make(map[string]interface{})
		res["items"] = l
		res["current"] = offset + 1
		res["pageSize"] = limit
		if total, totalErr := bim.GetTotalLb3dEditorBimToGltf(query); totalErr == nil {
			res["pages"] = math.Ceil(float64(total) / float64(limit))
			res["total"] = total
		}

		c.Data["json"] = server.RequestSuccess(res)
	}
	c.ResponseJson()
}

// Put ...
// @Title 更新数据
// @Description update the Lb3dEditorBimToGltf
// @Param	id		path 	string	true		"转换数据id"
// @Param	body		body 	bim.Lb3dEditorBimToGltf	true		"body for Lb3dEditorBimToGltf content"
// @Success 200 {string} "OK"
// @Failure 403 :id is not int
// @router /bim2gltf/update/:id [put]
func (c *Lb3dEditorBimToGltfController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := bim.Lb3dEditorBimToGltf{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := bim.UpdateLb3dEditorBimToGltfById(&v); err == nil {
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
// @router /bim2gltf/del/:id [delete]
func (c *Lb3dEditorBimToGltfController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := bim.DeleteLb3dEditorBimToGltf(id); err == nil {
		c.Data["json"] = server.RequestSuccess("delete success!")
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}

// DoRvtUpload ...
// @Title bim-rvt文件上传
// @Description 处理rvt文件上传
// @Param	file	body 	File	true	"rvt文件"
// @Success 200 {string} filePath
// @Failure 403 file is empty
// @router /bim2gltf/uploadRvt [post]
func (c *Lb3dEditorBimToGltfController) DoRvtUpload() {
	file, fileHeader, err := c.GetFile("file")

	// 关闭上传的文件，不然的话会出现临时文件不能清除的情况
	defer file.Close()

	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
		c.ResponseJson()
		return
	}

	var fileExt = path.Ext(fileHeader.Filename)
	var fileName = strings.TrimSuffix(fileHeader.Filename, fileExt)

	if strings.ToLower(fileExt) != ".rvt" {
		c.Data["json"] = server.RequestFail("只能上传rvt文件！")
		c.ResponseJson()
		return
	}

	currentDate := time.Now().Format("20060102")

	//存储的文件夹
	saveFolder := "bim/rvt/" + currentDate
	//存储的路径
	savePath := saveFolder + "/" + fileName + "-" + strconv.FormatInt(time.Now().Unix(), 10) + fileExt

	/* 保存位置在 static/upload, 没有文件夹要先创建 */
	//先判断创建目录
	if _, err = os.Stat("static/upload/" + saveFolder); os.IsNotExist(err) {
		err := os.MkdirAll("static/upload/"+saveFolder, os.ModePerm)
		if err != nil {
			c.Data["json"] = server.RequestFail("服务端创建文件夹" + currentDate + "失败！error=" + err.Error())
			c.ResponseJson()
			return
		}
	}

	err = c.SaveToFile("file", "static/upload/"+savePath)
	if err != nil {
		c.Data["json"] = server.RequestFail("服务器端文件保存失败！error=" + err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess("static/upload/" + savePath)
	}
	c.ResponseJson()
}
