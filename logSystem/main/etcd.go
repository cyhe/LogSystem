package main

import (
	"time"
	"fmt"
	etcd_client "github.com/coreos/etcd/clientv3"
	"github.com/astaxie/beego/logs"
	"context"
	"strings"
	"encoding/json"
	"logSystem/conf"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"logSystem/tailf"
)

type EtcdClient struct {
	client *etcd_client.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
)

func InitEtcd(addr string, key string) (collectConf []conf.CollectConf, err error) {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:2379", "localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}
	fmt.Printf("connect etcd succ")

	etcdClient = &EtcdClient{
		client: client,
	}

	// 从etcd获取配置

	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}

	for _, ip := range localIPArray {
		etcdkey := fmt.Sprintf("%s%s", key, ip)
		etcdClient.keys = append(etcdClient.keys, etcdkey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := client.Get(ctx, etcdkey)
		if err != nil {
			logs.Error("client get from etcd failed, err:", err)
			continue
		}
		cancel()
		logs.Debug("resp from etcd:%v", resp.Kvs)
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdkey {
				err := json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed, err:%v", err)
					continue
				}
				logs.Debug("log congig is %v", collectConf)
			}
		}
	}

	initEtcdWatch()

	return
}

/*
	监听配置的节点有没有变化
*/
func initEtcdWatch() {
	for _, key := range etcdClient.keys {
		go watchKey(key)
	}
}

func watchKey(key string) {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:2379", "localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	logs.Debug("begin watch etcdkey : %s", key)

	for {
		rch := client.Watch(context.Background(), key)
		// 最新的配置
		var collectConf []conf.CollectConf
		var getConfSucc = true
		for wresp := range rch {
			for _, ev := range wresp.Events {

				if ev.Type == mvccpb.DELETE { //delete
					logs.Warn("key[%s] `s config deleted", key)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key [%s ], unmarshal[%s], err:%v", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				// 把最新的配置给tailf
				tailf.UpdateConfig(collectConf)
			}
		}

	}
}
