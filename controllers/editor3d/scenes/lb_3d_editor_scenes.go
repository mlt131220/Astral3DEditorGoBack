package scenes

import (
	"encoding/json"
	"errors"
	"es-3d-editor-go-back/controllers"
	"es-3d-editor-go-back/models/editor3d/scenes"
	"es-3d-editor-go-back/server"
	"github.com/google/uuid"
	"math"
	"strings"
)

// Lb3dEditorScenesController operations for Lb3dEditorScenes
type Lb3dEditorScenesController struct {
	controllers.BaseController
}

// URLMapping ...
func (c *Lb3dEditorScenesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title 新增场景
// @Description create Lb3dEditorScenes
// @Param	body	scenes.Lb3dEditorScenes	true		"body for Lb3dEditorScenes content"
// @Success 201 {int} scenes.Lb3dEditorScenes
// @Failure 403 body is empty
// @router /scenes/add [post]
func (c *Lb3dEditorScenesController) Post() {
	// 由于是开源项目，此处新增场景是先检查数据量，如果已经有500个场景，则不允许新增，否则允许新增
	total, _ := scenes.GetTotalLb3dEditorScenes(nil)
	if total >= 500 {
		c.Data["json"] = server.RequestFail("共享项目场景数量已达上限（500个），不允许新增")
		c.ResponseJson()
		return
	}

	var v scenes.Lb3dEditorScenes
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		// 使用UUID
		v.Id = uuid.New().String()

		if _, err := scenes.AddLb3dEditorScenes(&v); err == nil {
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
// @Title 查询场景
// @Description get Lb3dEditorScenes by id
// @Param	id		path 	string	true		"场景id"
// @Success 200 {object} scenes.Lb3dEditorScenes
// @Failure 403 :id is empty
// @router /scenes/get/:id [get]
func (c *Lb3dEditorScenesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	v, err := scenes.GetLb3dEditorScenesById(idStr)
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
// @Success 200 {object} scenes.Lb3dEditorScenes
// @Failure 403
// @router /scenes/getAll [get]
func (c *Lb3dEditorScenesController) GetAll() {
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

	l, err := scenes.GetAllLb3dEditorScenes(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		res := make(map[string]interface{})
		res["items"] = l
		res["current"] = offset + 1
		res["pageSize"] = limit
		if total, totalErr := scenes.GetTotalLb3dEditorScenes(query); totalErr == nil {
			res["pages"] = math.Ceil(float64(total) / float64(limit))
			res["total"] = total
		}

		c.Data["json"] = server.RequestSuccess(res)
	}
	c.ResponseJson()
}

// Put ...
// @Title 更新数据
// @Description update the Lb3dEditorScenes
// @Param	id		path 	string	true		"场景id"
// @Param	body		body 	scenes.Lb3dEditorScenes	true		"body for Lb3dEditorScenes content"
// @Success 200 {scenes.Lb3dEditorScenes} scenes.Lb3dEditorScenes
// @Failure 403 :id is not int
// @router /scenes/update/:id [put]
func (c *Lb3dEditorScenesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	v := scenes.Lb3dEditorScenes{Id: idStr}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := scenes.UpdateLb3dEditorScenesById(&v); err == nil {
			c.Data["json"] = server.RequestSuccess(v)
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
// @Description 通过Id删除场景
// @Param	id	 path 	string	true	"场景表Id"
// @Success 200 {string} delete success!
// @Failure 403 "删除失败"
// @router /scenes/del/:id [delete]
func (c *Lb3dEditorScenesController) Delete() {
	id := c.Ctx.Input.Param(":id")
	if err := scenes.DeleteLb3dEditorScenes(id); err == nil {
		c.Data["json"] = server.RequestSuccess("删除成功！")
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}
