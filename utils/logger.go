package utils

import (
	"fmt"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

func InitLogger() (err error) {
	BConfig, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		fmt.Println("config init error:", err)
		return
	}
	maxLines, lerr := BConfig.Int64("log::maxLines")
	if lerr != nil {
		maxLines = 1000
	}

	filename, _ := BConfig.String("log::logPath")
	level, _ := BConfig.Int("log::logLevel")
	//maxSize, _ := BConfig.Int("log::maxSize")
	confStr := fmt.Sprintf(`{"filename":"%s","level":%d,"maxLines":%d,"separate":["error", "warning", "notice", "info"]}`, filename, level, maxLines)

	// 自定义输出格式
	//f := &logs.PatternLogFormatter{
	//	Pattern:    "[%T] %w  [%F:%n] >> %m",
	//	WhenFormat: "2006-01-02",
	//}
	//logs.RegisterFormatter("pattern", f)
	//_ = logs.SetGlobalFormatter("pattern")

	logs.SetLogger(logs.AdapterMultiFile, string(confStr))
	// 输出文件名和行号
	logs.EnableFuncCallDepth(true)
	// 异步输出日志
	logs.Async(1e3)
	return
}
