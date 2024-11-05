package scenes

import (
	"errors"
	"es-3d-editor-go-back/controllers/system"
	"fmt"
	"github.com/beego/beego/v2/adapter/logs"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Lb3dEditorScenes struct {
	Id                string    `orm:"column(id);pk" description:"主键ID,UUID"  json:"id"`
	SceneType         string    `orm:"column(sceneType);size(24);null" description:"场景类型" json:"sceneType"`
	HasDrawing        int       `orm:"column(hasDrawing)" description:"场景是否包含图纸 0:false  1:true" json:"hasDrawing"`
	SceneIntroduction string    `orm:"column(sceneIntroduction);size(255);null" description:"场景描述" json:"sceneIntroduction"`
	SceneName         string    `orm:"column(sceneName);size(255);null" description:"场景名称" json:"sceneName"`
	SceneVersion      int       `orm:"column(sceneVersion);" description:"场景版本" json:"sceneVersion"`
	Zip               string    `orm:"column(zip);size(128)" description:"场景zip包" json:"zip"`
	ZipSize           string    `orm:"column(zipSize);size(32)" description:"场景zip包大小" json:"zipSize"`
	CoverPicture      string    `orm:"column(coverPicture);size(40000)" description:"保存场景时自动生成的封面图url" json:"coverPicture"`
	ExampleSceneId    string    `orm:"column(exampleSceneId);null" description:"创建项目时来源于哪一个示例模板项目，null代表从空项目创建" json:"exampleSceneId"`
	ProjectType       int       `orm:"column(projectType);" description:"项目类型。0：Web3D-THREE  1：WebGIS-Cesium" json:"projectType"`
	CesiumConfig      string    `orm:"column(cesiumConfig);size(1000);null" description:"WebGIS-Cesium 类型项目的基础Cesium配置" json:"cesiumConfig"`
	UpdateTime        time.Time `orm:"column(updateTime);type(datetime);auto_now" json:"updateTime"`
	CreateTime        time.Time `orm:"column(createTime);type(datetime);null;auto_now_add" json:"createTime"`
	// 不使用关联查询，意义不大
	//ExampleScene      *Lb3dEditorScenesExample `orm:"rel(fk);null" json:"exampleScene"` // 外键关联
}

func (t *Lb3dEditorScenes) TableName() string {
	return "lb_3d_editor_scenes"
}

func init() {
	orm.RegisterModel(new(Lb3dEditorScenes))
}

// AddLb3dEditorScenes insert a new Lb3dEditorScenes into database and returns
// last inserted Id on success.
func AddLb3dEditorScenes(m *Lb3dEditorScenes) (id int64, err error) {
	o := orm.NewOrm()

	// 原生sql插入
	//sql := "INSERT INTO `lb_3d_editor_scenes` (`id`,`sceneType`, `hasDrawing`, `sceneIntroduction`, `sceneName`, `sceneVersion`, `zip`, `zipSize`, `coverPicture`, `projectType`, `example_scene_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	//res, err := o.Raw(sql, m.Id, m.SceneType, m.HasDrawing, m.SceneIntroduction, m.SceneName, m.SceneVersion, m.Zip, m.ZipSize, m.CoverPicture, m.ProjectType, m.ExampleSceneId).Exec()
	//id, err = res.LastInsertId()

	id, err = o.Insert(m)
	return
}

// GetLb3dEditorScenesById retrieves Lb3dEditorScenes by Id. Returns error if
// Id doesn't exist
func GetLb3dEditorScenesById(id string) (v *Lb3dEditorScenes, err error) {
	o := orm.NewOrm()
	v = &Lb3dEditorScenes{Id: id}
	if err = o.Read(v); err == nil {
		// 处理zip和CoverPicture字段，如果该字段为空，则从关联的ExampleScene中获取zip字段的值
		if (v.Zip == "" || v.CoverPicture == "") && v.ExampleSceneId != "" {
			exampleScene, err := GetLb3dEditorScenesExampleById(v.ExampleSceneId)
			if err == nil {
				if v.Zip == "" {
					v.Zip = exampleScene.Zip
				}
				if v.CoverPicture == "" {
					v.CoverPicture = exampleScene.CoverPicture
				}
			}
		}
		return v, nil
	}

	// 不使用关联查询，意义不大
	//if err = o.QueryTable(new(Lb3dEditorScenes)).Filter("Id", id).RelatedSel().One(v); err == nil {
	//	fmt.Println("example_scene:", v.ExampleScene)
	//	// 处理zip字段，如果该字段为空，则从关联的ExampleScene中获取zip字段的值
	//	if v.Zip == "" && v.ExampleScene != nil {
	//		v.Zip = v.ExampleScene.Zip
	//	}
	//
	//	// 返回值中不包含ExampleScene字段，避免返回过多数据
	//	v.ExampleScene = nil
	//
	//	return v, nil
	//}

	return nil, err
}

// GetAllLb3dEditorScenes retrieves all Lb3dEditorScenes matches certain condition. Returns empty list if
// no records exist
func GetAllLb3dEditorScenes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorScenes))
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

	var l []Lb3dEditorScenes
	qs = qs.OrderBy(sortFields...)
	_, err = qs.Limit(limit, offset).All(&l, fields...)

	// 处理CoverPicture字段，如果该字段为空，则从关联的ExampleScene中获取CoverPicture字段的值
	for i, v := range l {
		if v.CoverPicture == "" && v.ExampleSceneId != "" {
			v.CoverPicture = GetCoverPicture(v.ExampleSceneId)
			l[i] = v
		}
	}

	if err == nil {
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

// GetTotalLb3dEditorScenes 获取查询总数量
func GetTotalLb3dEditorScenes(query map[string]string) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Lb3dEditorScenes))

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

// UpdateLb3dEditorScenesById updates Lb3dEditorScenes by Id and returns error if
// the record to be updated doesn't exist
func UpdateLb3dEditorScenesById(m *Lb3dEditorScenes) (err error) {
	o := orm.NewOrm()
	v := Lb3dEditorScenes{Id: m.Id}
	// 确定数据库中是否存在ID
	if err = o.Read(&v); err == nil {
		// 先移除 zip 字段对应的文件
		if v.Zip != "" {
			// 从又拍云删除文件父级的文件夹及其下所有文件
			var folder = v.Zip[:strings.LastIndex(v.Zip, "/")] + "/"
			if err := system.UpYunDeleteAll(folder); err != nil {
				fmt.Println(folder, " 删除文件失败: ", err)
				return err
			}
		}

		// 判断封面图是否变更
		if m.CoverPicture != v.CoverPicture {
			// 从又拍云删除原封面图文件
			if v.CoverPicture != "" {
				if err := system.UpYunDelete(v.CoverPicture); err != nil {
					fmt.Println(v.CoverPicture, " 删除文件失败: ", err)
					return err
				}
			}
		}

		if _, err = o.Update(m); err == nil {
			logs.Info("lb_3d_editor_scenes 更新数据:", v.Id)
		}
	}
	return
}

// DeleteLb3dEditorScenes 通过表Id删除场景
func DeleteLb3dEditorScenes(id string) (err error) {
	o := orm.NewOrm()
	v := Lb3dEditorScenes{Id: id}
	// 数据库中存在确定id
	if err := o.Read(&v); err == nil {
		// 先移除 zip 字段对应的文件
		if v.Zip != "" && strings.Contains(v.Zip, "/") {
			// 从又拍云删除文件父级的文件夹及其下所有文件
			var folder = v.Zip[:strings.LastIndex(v.Zip, "/")] + "/"
			if err := system.UpYunDeleteAll(folder); err != nil {
				fmt.Println(folder, " 删除文件失败: ", err)
				return err
			}
		}

		// 从又拍云删除原封面图文件
		if v.CoverPicture != "" {
			if err := system.UpYunDelete(v.CoverPicture); err != nil {
				fmt.Println(v.CoverPicture, " 删除文件失败: ", err)
				return err
			}
		}

		if _, err = o.Delete(&Lb3dEditorScenes{Id: id}); err == nil {
			logs.Info("lb_3d_editor_scenes 删除数据:", v)
		} else {
			return err
		}
	}
	return
}
