package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logSystem/tailf"
	"logSystem/conf"
	lsLogs "logSystem/logs"
	"logSystem/kafka"
)

func main() {

	// 初始化配置文件
	fileName := "./conf/logagent.conf"
	err := conf.InitConfig("ini", fileName)
	if err != nil {
		fmt.Printf("load conf failed, err:%v\n", err)
		panic("load conf failed")
		return
	}

	logs.Debug("init config succ")

	// 初始化日志文件
	err = lsLogs.InitLogger()
	if err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		panic("init logger failed")
		return
	}
	logs.Debug("init logger succ")

	// 初始化etcd  获取配置
	c, err := InitEtcd(conf.AppConfig.EtcdAddr, conf.AppConfig.EtcdKey)
	if err != nil {
		logs.Error("init etcd failed,err:", err)
	}
	logs.Debug("init etcd succ")

	// 初始化tail组件
	err = tailf.InitTail(c, conf.AppConfig.ChanSize)
	if err != nil {
		logs.Error("init tail failed, err:", err)
		return
	}
	logs.Debug("init tailf succ")

	// 初始化kafka
	err = kafka.InitKafka(conf.AppConfig.KafkaAddr)
	if err != nil {
		logs.Error("init kafka failed, err:", err)
		return
	}

	logs.Debug("init kafka succ")

	logs.Debug("initalize all succ")

	//go func() {
	//	var count int
	//	for {
	//		count++
	//		logs.Debug("test for logger %d", count)
	//		time.Sleep(time.Millisecond * 1000)
	//	}
	//}()

	// 启动服务
	serverRun()
	if err != nil {
		logs.Error("server run failed, err:", err)
	}

	logs.Info("program exited")

}
