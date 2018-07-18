package logs

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"fmt"
	"logSystem/conf"
)

func convertLogLevel(level string) int {
	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "trace":
		return logs.LevelTrace
	case "info":
		return logs.LevelInformational
	default:
		return logs.LevelDebug
	}
}

func InitLogger() (err error) {
	// 配置log组件
	config := make(map[string]interface{})
	// 日志的路径,文件名
	config["filename"] = conf.AppConfig.LogPath
	// 日志级别(开发环境)
	config["level"] = convertLogLevel(conf.AppConfig.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("init logger failed, Marshal, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
