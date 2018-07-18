package main

import (
	"fmt"
	"time"
	"encoding/json"
	etcd_client "github.com/coreos/etcd/clientv3"
	"context"
	"logSystem/conf"
)

// 获取ip可以搞一个基础库,
// 或者 web管理页面,也不需要获取ip,实际情况下, 日志搜集web页面,可以调用运维的系统,从中获取ip(运维应该有这个系统)
const (
	etcdKey = "/Users/cyhe/go/src/logSystem/conf/192.168.31.62"
)

//// 来描述一个配置
//type LogConf struct {
//	Path  string `json:"path"`  // 文件路径
//	Topic string `json:"topic"` // 所属topic
//	//sendQps int  // 1秒发多少个,高峰期可以发少点,低峰期发少一点
//}


func SetLogConfToEtcd() {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect faied,err:", err)
		return
	}
	fmt.Println("connect succ \n")
	defer client.Close()

	var logConArr []conf.CollectConf

	logConArr = append(
		logConArr,
		conf.CollectConf{
			LogPath:  "/usr/local/var/log/nginx/access.log",
			Topic: "nginx_log",
		},
	)

	logConArr = append(
		logConArr,
		conf.CollectConf{
			LogPath:  "/Users/esirnus/Documents/logs/error.log",
			Topic: "nginx_log_err",
		},
	)

	data, err := json.Marshal(logConArr)
	if err != nil {
		fmt.Println("json failed, err:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx,etcdKey,string(data))

	cancel()

	if err != nil {
		fmt.Println("put failed, err", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, etcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err", err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s \n", ev.Key, ev.Value)
	}
}

func main() {
	SetLogConfToEtcd()
}