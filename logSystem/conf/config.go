package conf

import (
	"github.com/astaxie/beego/config"
	"fmt"
	"errors"
)

var (
	AppConfig *Config
)

type Config struct {
	LogLevel    string
	LogPath     string
	CollectConf []CollectConf
	ChanSize    int
	KafkaAddr   string
	EtcdAddr    string
	EtcdKey     string
}

type CollectConf struct {
	LogPath string `json:"log_path"`
	Topic   string `json:"topic"`
}

func loadCollectConf(conf config.Configer) (err error) {
	var cc CollectConf
	cc.LogPath = conf.String("collect::log_path")
	if len(cc.LogPath) == 0 {
		err = errors.New("invalid collect::log_path")
		return
	}
	cc.Topic = conf.String("collect::topic")
	if len(cc.Topic) == 0 {
		err = errors.New("invalid collect::topic")
		return
	}
	AppConfig.CollectConf = append(AppConfig.CollectConf, cc)
	return
}

func InitConfig(confType, fileName string) (err error) {

	// 初始化配置库
	conf, err := config.NewConfig(confType, fileName)
	if err != nil {
		fmt.Println("new config failed,err :", err)
		return
	}

	// 读取 写日志 配置项
	AppConfig = &Config{
		LogLevel: conf.String("logs::log_level"),
		LogPath:  conf.String("logs::log_path"),
	}

	AppConfig.ChanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		AppConfig.ChanSize = 100
	}

	if len(AppConfig.LogLevel) == 0 {
		AppConfig.LogLevel = "debug"
	}
	if len(AppConfig.LogPath) == 0 {
		AppConfig.LogPath = "./logs"
	}

	AppConfig.KafkaAddr = conf.String("kafka::server_addr")
	if len(AppConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka addr")
		return
	}

	AppConfig.EtcdAddr = conf.String("etcd::etcd_addr")
	if len(AppConfig.EtcdAddr) == 0 {
		err = fmt.Errorf("invalid etcd addr")
		return
	}

	AppConfig.EtcdKey = conf.String("etcd::config_key")
	if len(AppConfig.EtcdKey) == 0 {
		err = fmt.Errorf("invalid etcd Key")
		return
	}

	// 读取 收集日志 配置项
	err = loadCollectConf(conf)
	if err != nil {
		fmt.Printf("load collect conf failed, err:%v\n", err)
		return
	}

	return
}
