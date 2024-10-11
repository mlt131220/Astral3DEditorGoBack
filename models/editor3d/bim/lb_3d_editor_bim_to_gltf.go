package bim

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Lb3dEditorBimToGltf struct {
	Id                 int       `orm:"column(id);pk;auto" json:"id"`
	BimFilePath        string    `orm:"column(bim_file_path);size(255)" description:"bim源文件路径" json:"bimFilePath"`
	BimFileSize        float64   `orm:"column(bim_file_size);null" description:"bim源文件大小" json:"bimFileSize"`
	ConversionDuration float64   `orm:"column(conversion_duration);null" description:"转换时长（s）" json:"conversionDuration"`
	ConversionStatus   int       `orm:"column(conversion_status);" description:"0 转换中 1 转换完成 2 转换失败" json:"conversionStatus"`
	FileName           string    `orm:"column(file_name);size(255)" description:"文件名" json:"fileName"`
	Thumbnail          string    `orm:"column(thumbnail);size(255)" description:"缩略图" json:"thumbnail"`
	GltfFilePath       string    `orm:"column(gltf_file_path);size(255)" description:"转换后的gltf文件路径" json:"gltfFilePath"`
	GltfFileSize       float64   `orm:"column(gltf_file_size);null" description:"转换后的gltf文件大小" json:"gltfFileSize"`
	DelTag             int8      `orm:"column(delTag)" description:"删除标记，0 未删除 1 已删除" json:"delTag"`
	CreateTime         time.Time `orm:"column(createTime);type(datetime);null;auto_now_add" json:"createTime"`
	UpdateTime         time.Time `orm:"column(updateTime);type(datetime);auto_now" json:"updateTime"`
	DelTime            time.Time `orm:"column(delTime);type(datetime);null" json:"delTime"`
}

func (t *Lb3dEditorBimToGltf) TableName() string {
	return "lb_3d_editor_bim_to_gltf"
}

func init() {
	orm.RegisterModel(new(Lb3dEditorBimToGltf))
}

// AddLb3dEditorBimToGltf insert a new Lb3dEditorBimToGltf into database and returns
// last inserted Id on success.
func AddLb3dEditorBimToGltf(m *Lb3dEditorBimToGltf) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLb3dEditorBimToGltfById retrieves Lb3dEditorBimToGltf by Id. Returns error if
// Id doesn't exist
func GetLb3dEditorBimToGltfById(id int) (v *Lb3dEditorBimToGltf, err error) {
	o := orm.NewOrm()
	v = &Lb3dEditorBimToGltf{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLb3dEditorBimToGltf retrieves all Lb3dEditorBimToGltf matches certain condition. Returns empty list if
// no records exist
func GetAllLb3dEditorBimToGltf(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorBimToGltf))
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

	var l []Lb3dEditorBimToGltf
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

// GetTotalLb3dEditorBimToGltf 获取总数量
func GetTotalLb3dEditorBimToGltf(query map[string]string) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorBimToGltf))

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			//大小写不敏感模糊查询 + "__icontains"
			qs = qs.Filter(k+"__icontains", v)
		}
	}
	if cnt, err := qs.Count(); err == nil {
		return cnt, nil
	}
	return 0, err
}

// UpdateLb3dEditorBimToGltf updates Lb3dEditorBimToGltf by Id and returns error if
// the record to be updated doesn't exist
func UpdateLb3dEditorBimToGltfById(m *Lb3dEditorBimToGltf) (err error) {
	o := orm.NewOrm()
	v := Lb3dEditorBimToGltf{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLb3dEditorBimToGltf deletes Lb3dEditorBimToGltf by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLb3dEditorBimToGltf(id int) (err error) {
	o := orm.NewOrm()
	v := Lb3dEditorBimToGltf{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Lb3dEditorBimToGltf{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
