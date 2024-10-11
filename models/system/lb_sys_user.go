package system

import (
	"errors"
	"es-3d-editor-go-back/struct"
	"es-3d-editor-go-back/utils/jwt"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
)

type LbSysUser struct {
	Id            int       `orm:"column(id);auto"`
	Username      string    `orm:"column(username);size(255)" description:"用户名"`
	Nickname      string    `orm:"column(nickname);size(255)" description:"用户昵称"`
	Mobile        string    `orm:"column(mobile);size(20)" description:"用户手机号"`
	Password      string    `orm:"column(password);size(255)" description:"用户密码"`
	Sex           int8      `orm:"column(sex)" description:"性别， 0 表示女， 1 表示男"`
	Avatar        string    `orm:"column(avatar);size(255);null" description:"头像"`
	Email         string    `orm:"column(email);size(100);null" description:"邮箱"`
	DelTag        int8      `orm:"column(delTag)" description:"删除标记，0 未删除 1 已删除"`
	Salt          string    `orm:"column(salt);size(255)" description:"jwt 鉴权 SALT值"`
	LastLoginTime time.Time `orm:"column(lastLoginTime);type(datetime);null" description:"最后登录时间"`
	LastLoginIp   string    `orm:"column(lastLoginIp);size(255);null" description:"最后登录ip"`
	RegisterIp    string    `orm:"column(registerIp);size(255);null" description:"注册时的ip地址"`
	CreateTime    time.Time `orm:"column(createTime);type(datetime);auto_now_add"`
	UpdateTime    time.Time `orm:"column(updateTime);type(datetime);auto_now_add"`
	DelTime       time.Time `orm:"column(delTime);type(datetime);null"`
}

func (t *LbSysUser) TableName() string {
	return "lb_sys_user"
}

func init() {
	orm.RegisterModel(new(LbSysUser))
}

// AddLbSysUser insert a new LbSysUser into database and returns
// last inserted Id on success.
func AddLbSysUser(m *LbSysUser) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLbSysUserById retrieves LbSysUser by Id. Returns error if
// Id doesn't exist
func GetLbSysUserById(id int) (v *LbSysUser, err error) {
	o := orm.NewOrm()
	v = &LbSysUser{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLbSysUser retrieves all LbSysUser matches certain condition. Returns empty list if
// no records exist
func GetAllLbSysUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(LbSysUser))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []LbSysUser
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateLbSysUser updates LbSysUser by Id and returns error if
// the record to be updated doesn't exist
func UpdateLbSysUserById(m *LbSysUser) (err error) {
	o := orm.NewOrm()
	v := LbSysUser{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLbSysUser deletes LbSysUser by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLbSysUser(id int) (err error) {
	o := orm.NewOrm()
	v := LbSysUser{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&LbSysUser{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// DoLogin "登录处理"
func DoLogin(lr *_struct.LoginRequest) (*_struct.LoginResponse, int, error) {
	fmt.Println(lr)

	// get username and password
	username := lr.UserName
	password := lr.Password

	// validate username and password is they are empty
	if len(username) == 0 || len(password) == 0 {
		return nil, http.StatusBadRequest, errors.New("错误: 用户名或密码不能为空！")
	}

	o := orm.NewOrm()

	// check if the username exists
	user := &LbSysUser{Username: username}
	err := o.Read(user, "Username")
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("错误: 用户名不存在！")
	}

	// generate the password hash
	hash, err := jwt.GeneratePassHash(password, user.Salt)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if hash != user.Password {
		return nil, http.StatusBadRequest, errors.New("错误: 密码错误！")
	}

	// 生成token
	tokenString, err := jwt.GenerateToken(lr, user.Id, 0)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &_struct.LoginResponse{
		UserName: user.Username,
		UserID:   user.Id,
		Token:    tokenString,
	}, http.StatusOK, nil
}

// DoRegister 用户注册
func DoRegister(cr *_struct.RegisterRequest) (*_struct.RegisterResponse, int, error) {
	o := orm.NewOrm()

	// 检查用户名是否已存在
	userNameCheck := LbSysUser{Username: cr.UserName}
	err := o.Read(&userNameCheck, "Username")
	if err == nil {
		return nil, http.StatusBadRequest, errors.New("用户名已存在！")
	}

	//生成SALT值
	saltKey, err := jwt.GenerateSalt()
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// generate password hash
	hash, err := jwt.GeneratePassHash(cr.Password, saltKey)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// create user
	user := LbSysUser{}
	user.Username = cr.UserName
	user.Password = hash
	user.Salt = saltKey

	_, err = o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	return &_struct.RegisterResponse{
		UserID:   user.Id,
		UserName: user.Username,
	}, http.StatusCreated, nil
}
