package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

func main() {
	// 初始化配置文件
	err := initConfig("ini", "./logtransfer/conf/log_transfer.conf")
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(logConfig)

	// 初始化logger
	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
		return
	}
	logs.Debug("init logger succ")

	// 初始化kafka消费
	err = initKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	logs.Debug("init kafka succ")

	// 初始化elasticsearch
	err = initES(logConfig.ESAddr)
	if err != nil {
		logs.Error("init elasticsearch failed, err:%v", err)
		return
	}
	logs.Debug("init es succ")

	fmt.Println("sdfsdfdsfd")
	// 初始化配置文件
	err = run()
	if err != nil {
		logs.Error("server run failed, err:", err)
		return
	}

	logs.Warn("program exited")
}
