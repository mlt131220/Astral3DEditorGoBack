package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/bim2gltf/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "AddAndConversion",
            Router: `/bim2gltf/addAndConversion`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/bim2gltf/del/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/bim2gltf/get/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/bim2gltf/getAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/bim2gltf/update/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/bim:Lb3dEditorBimToGltfController"],
        beego.ControllerComments{
            Method: "DoRvtUpload",
            Router: `/bim2gltf/uploadRvt`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/cad/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/cad/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/cad/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/cad/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"],
        beego.ControllerComments{
            Method: "Dwg2Dxf",
            Router: `/cad/dwg2dxf`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/cad:Lb3dEditorCadController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/cad/getAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/scenes/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/scenes/del/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/scenes/get/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/scenes/getAll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/scenes/update/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/sceneExample`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/sceneExample`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/sceneExample/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/sceneExample/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/editor3d/scenes:Lb3dEditorScenesExampleController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/sceneExample/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUpYunController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUpYunController"],
        beego.ControllerComments{
            Method: "DOBlob",
            Router: `/upyun/blob`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUpYunController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUpYunController"],
        beego.ControllerComments{
            Method: "DoRemove",
            Router: `/upyun/remove`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUpYunController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUpYunController"],
        beego.ControllerComments{
            Method: "DoUpload",
            Router: `/upyun/upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"],
        beego.ControllerComments{
            Method: "DOBlob",
            Router: `/blob`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"],
        beego.ControllerComments{
            Method: "DoDownload",
            Router: `/downloadFile`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"],
        beego.ControllerComments{
            Method: "DoRemoveFile",
            Router: `/removeFile`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUploadController"],
        beego.ControllerComments{
            Method: "DoUpload",
            Router: `/upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "Register",
            Router: `/register`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/user/createUser`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/user/delUser/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/user/getAllUser`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/user/getUserInfo/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysUserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/user/updateUser/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysWebSocketController"] = append(beego.GlobalControllerRouter["es-3d-editor-go-back/controllers/system:LbSysWebSocketController"],
        beego.ControllerComments{
            Method: "WsHandler",
            Router: `/ws`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
