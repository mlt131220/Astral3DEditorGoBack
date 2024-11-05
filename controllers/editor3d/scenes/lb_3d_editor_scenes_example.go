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

// Lb3dEditorScenesExampleController operations for Lb3dEditorScenesExample
type Lb3dEditorScenesExampleController struct {
	controllers.BaseController
}

// URLMapping ...
func (c *Lb3dEditorScenesExampleController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title 新增
// @Description create Lb3dEditorScenesExample
// @Param	body		body 	models.Lb3dEditorScenesExample	true		"body for Lb3dEditorScenesExample content"
// @Success 201 {int} models.Lb3dEditorScenesExample
// @Failure 403 body is empty
// @router /sceneExample [post]
func (c *Lb3dEditorScenesExampleController) Post() {
	var v scenes.Lb3dEditorScenesExample
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		// 使用UUID
		v.Id = uuid.New().String()

		if _, err := scenes.AddLb3dEditorScenesExample(&v); err == nil {
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
// @Title 按id查询
// @Description get Lb3dEditorScenesExample by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Lb3dEditorScenesExample
// @Failure 403 :id is empty
// @router /sceneExample/:id [get]
func (c *Lb3dEditorScenesExampleController) GetOne() {
	id := c.Ctx.Input.Param(":id")
	v, err := scenes.GetLb3dEditorScenesExampleById(id)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess(v)
	}
	c.ResponseJson()
}

// GetAll ...
// @Title 查询全部
// @Description get Lb3dEditorScenesExample
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Lb3dEditorScenesExample
// @Failure 403
// @router /sceneExample [get]
func (c *Lb3dEditorScenesExampleController) GetAll() {
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

	l, err := scenes.GetAllLb3dEditorScenesExample(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		res := make(map[string]interface{})
		res["items"] = l
		res["current"] = offset + 1
		res["pageSize"] = limit
		if total, totalErr := scenes.GetTotalLb3dEditorSceneExample(query); totalErr == nil {
			res["pages"] = math.Ceil(float64(total) / float64(limit))
			res["total"] = total
		}

		c.Data["json"] = server.RequestSuccess(res)
	}
	c.ResponseJson()
}

// Put ...
// @Title Put
// @Description update the Lb3dEditorScenesExample
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Lb3dEditorScenesExample	true		"body for Lb3dEditorScenesExample content"
// @Success 200 {object} models.Lb3dEditorScenesExample
// @Failure 403 :id is not int
// @router /sceneExample/:id [put]
func (c *Lb3dEditorScenesExampleController) Put() {
	id := c.Ctx.Input.Param(":id")
	v := scenes.Lb3dEditorScenesExample{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := scenes.UpdateLb3dEditorScenesExampleById(&v); err == nil {
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
// @Description delete the Lb3dEditorScenesExample
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /sceneExample/:id [delete]
func (c *Lb3dEditorScenesExampleController) Delete() {
	id := c.Ctx.Input.Param(":id")
	if err := scenes.DeleteLb3dEditorScenesExample(id); err == nil {
		c.Data["json"] = server.RequestSuccess("delete success!")
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ResponseJson()
}
