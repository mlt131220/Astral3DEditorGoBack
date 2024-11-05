package scenes

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Lb3dEditorScenesExample struct {
	Id                string    `orm:"column(id);pk" description:"主键ID,UUID" json:"id"`
	SceneName         string    `orm:"column(sceneName);size(255);null" description:"场景名称" json:"sceneName"`
	SceneType         string    `orm:"column(sceneType);size(24);null" description:"场景类型" json:"sceneType"`
	SceneVersion      int       `orm:"column(sceneVersion);null" description:"场景版本" json:"sceneVersion"`
	SceneIntroduction string    `orm:"column(sceneIntroduction);size(255);null" description:"场景描述" json:"sceneIntroduction"`
	CoverPicture      string    `orm:"column(coverPicture);size(4000)" description:"保存场景时自动生成的封面图url" json:"coverPicture"`
	HasDrawing        int       `orm:"column(hasDrawing)" description:"场景是否包含图纸 0:false  1:true" json:"hasDrawing"`
	ProjectType       int       `orm:"column(projectType)" description:"示例项目类型。0：Web3D-THREE  1：WebGIS-Cesium" json:"projectType"`
	CesiumConfig      string    `orm:"column(cesiumConfig);size(1000);null" description:"WebGIS-Cesium 类型项目的基础Cesium配置" json:"cesiumConfig"`
	Zip               string    `orm:"column(zip);size(128)" description:"场景zip包" json:"zip"`
	ZipSize           string    `orm:"column(zipSize);size(32)" description:"场景zip包大小" json:"zipSize"`
	CreateTime        time.Time `orm:"column(createTime);type(datetime);null;auto_now_add" json:"createTime"`
	UpdateTime        time.Time `orm:"column(updateTime);type(datetime);auto_now" json:"updateTime"`
	DelTag            int8      `orm:"column(delTag)" description:"删除标记，0 未删除 1 已删除"`
	DelTime           time.Time `orm:"column(delTime);type(datetime);null" description:"删除时间"`
}

func (t *Lb3dEditorScenesExample) TableName() string {
	return "lb_3d_editor_scenes_example"
}

// CoverPictureMap 存储封面图片的map
var CoverPictureMap = make(map[string]string)

// 用于确保初始化只发生一次的标志
var isInitialized = false
var mu sync.Mutex

func init() {
	orm.RegisterModel(new(Lb3dEditorScenesExample))

	//SetCoverPictureMap()
}

func SetCoverPictureMap() {
	mu.Lock()
	defer mu.Unlock()

	if isInitialized {
		return // 如果已经初始化过，则直接返回
	}

	l, err := Lb3dEditorScenesExampleGetAll()
	if err == nil {
		for _, v := range l {
			CoverPictureMap[v.Id] = v.CoverPicture
		}

		isInitialized = true // 标记为已初始化
		fmt.Println("GetCoverPicture id:", CoverPictureMap)
	}
}

func GetCoverPicture(id string) (coverPicture string) {
	if value, exists := CoverPictureMap[id]; exists {
		return value
	} else {
		SetCoverPictureMap()

		if value, exists = CoverPictureMap[id]; exists {
			return value
		} else {
			return ""
		}
	}
}

func Lb3dEditorScenesExampleGetAll() (vm []*Lb3dEditorScenesExample, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorScenesExample))
	var l []*Lb3dEditorScenesExample
	_, err = qs.All(&l)
	if err == nil {
		vm = l
	}
	return
}

// AddLb3dEditorScenesExample insert a new Lb3dEditorScenesExample into database and returns
// last inserted Id on success.
func AddLb3dEditorScenesExample(m *Lb3dEditorScenesExample) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLb3dEditorScenesExampleById retrieves Lb3dEditorScenesExample by Id. Returns error if
// Id doesn't exist
func GetLb3dEditorScenesExampleById(id string) (v *Lb3dEditorScenesExample, err error) {
	o := orm.NewOrm()
	v = &Lb3dEditorScenesExample{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLb3dEditorScenesExample retrieves all Lb3dEditorScenesExample matches certain condition. Returns empty list if
// no records exist
func GetAllLb3dEditorScenesExample(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorScenesExample))
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

	var l []Lb3dEditorScenesExample
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

// GetTotalLb3dEditorSceneExample 获取查询总数量
func GetTotalLb3dEditorSceneExample(query map[string]string) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorScenesExample))

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

// UpdateLb3dEditorScenesExample updates Lb3dEditorScenesExample by Id and returns error if
// the record to be updated doesn't exist
func UpdateLb3dEditorScenesExampleById(m *Lb3dEditorScenesExample) (err error) {
	o := orm.NewOrm()
	v := Lb3dEditorScenesExample{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLb3dEditorScenesExample deletes Lb3dEditorScenesExample by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLb3dEditorScenesExample(id string) (err error) {
	o := orm.NewOrm()
	v := Lb3dEditorScenesExample{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Lb3dEditorScenesExample{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
