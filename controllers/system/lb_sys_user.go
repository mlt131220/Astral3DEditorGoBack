package system

import (
	"encoding/json"
	"es-3d-editor-go-back/models/system"
	"es-3d-editor-go-back/server"
	"es-3d-editor-go-back/struct"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// LbSysUserController 通用-用户
type LbSysUserController struct {
	beego.Controller
}

// URLMapping ...
func (c *LbSysUserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Login", c.Login)
	c.Mapping("Register", c.Register)
}

// Post ...
// @Title 新增用户
// @Description create 新增用户
// @Param	body	body 	system.LbSysUser	true	"body for LbSysUser content"
// @Success 201 {object} system.LbSysUser
// @Failure 403 body is empty
// @router /user/createUser [post]
func (c *LbSysUserController) Post() {
	var v system.LbSysUser
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := system.AddLbSysUser(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = server.RequestSuccess(v)
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ServeJSON()
}

// GetOne ...
// @Title 获取用户
// @Description 通过id获取用户
// @Param	id		query 	string	true		"用户id"
// @Success 200 {object} system.LbSysUser
// @Failure 403 :id is empty
// @router /user/getUserInfo/:id [get]
func (c *LbSysUserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	//idStr := c.GetString("id")
	id, _ := strconv.Atoi(idStr)
	v, err := system.GetLbSysUserById(id)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess(v)
	}
	c.ServeJSON()
}

// GetAll ...
// @Title 获取全部用户
// @Description 获取全部用户
// @Param	query	query	string	false	"过滤字段. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"返回字段. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"排序字段. e.g. col1,col2 ..."
// @Param	order	query	string	false	"每个排序字段对应的顺序，如果是单个值，则应用于所有排序字段. e.g. desc,asc ..."
// @Param	limit	query	string	false	"限制结果集的大小。必须是整数"
// @Param	offset	query	string	false	"结果集的起始位置。必须是整数"
// @Success 200 {object} system.LbSysUser
// @Failure 403
// @router /user/getAllUser [get]
func (c *LbSysUserController) GetAll() {
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
				c.Data["json"] = server.RequestFail("error:无效的查询键/值对")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := system.GetAllLbSysUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = server.RequestFail(err.Error())
	} else {
		c.Data["json"] = server.RequestSuccess(l)
	}
	c.ServeJSON()
}

// Put ...
// @Title 更新用户信息
// @Description 通过id更新用户信息
// @Param	id		path 	string	true		"用户id"
// @Param	body		body 	system.LbSysUser	true		"body for LbSysUser content"
// @Success 200 {object} system.LbSysUser
// @Failure 403 :id is not int
// @router /user/updateUser/:id [put]
func (c *LbSysUserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := system.LbSysUser{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := system.UpdateLbSysUserById(&v); err == nil {
			c.Data["json"] = server.RequestSuccess("OK")
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ServeJSON()
}

// Delete ...
// @Title 删除用户
// @Description 通过id删除用户
// @Param	id		path 	string	true		"用户id"
// @Success 200 {object} system.LbSysUser
// @Failure 403 id is empty
// @router /user/delUser/:id [delete]
func (c *LbSysUserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := system.DeleteLbSysUser(id); err == nil {
		c.Data["json"] = server.RequestSuccess("OK")
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ServeJSON()
}

// Login ...
// @Title 登录
// @Description 用户登录接口
// @Param	body	body 	models.LoginRequest	 true	"登陆字段结构体"
// @Success 200 {object} models.LoginResponse
// @Failure 403 body is empty
// @Failure 500 "服务器错误"
// @router /login [post]
func (c *LbSysUserController) Login() {
	var v _struct.LoginRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if res, status, err := system.DoLogin(&v); err == nil {
			c.Ctx.Output.SetStatus(status)
			c.Data["json"] = server.RequestSuccess(res)
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ServeJSON()
}

// Register ...
// @Title 注册
// @Description 用户注册接口
// @Param	body	body 	models.RegisterRequest	true	"注册字段结构体"
// @Success 200 {object} models.RegisterResponse
// @Failure 403 body is empty
// @router /register [post]
func (c *LbSysUserController) Register() {
	var v _struct.RegisterRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if res, status, err := system.DoRegister(&v); err == nil {
			c.Ctx.Output.SetStatus(status)
			c.Data["json"] = server.RequestSuccess(res)
		} else {
			c.Data["json"] = server.RequestFail(err.Error())
		}
	} else {
		c.Data["json"] = server.RequestFail(err.Error())
	}
	c.ServeJSON()
}
