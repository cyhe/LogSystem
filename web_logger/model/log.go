package model

import (
	"github.com/astaxie/beego/logs"
	etcdclient "github.com/coreos/etcd/clientv3"
	"fmt"
	"encoding/json"
	"time"
	"context"
)

var (
	etcdClient *etcdclient.Client
)


type CollectConf struct {
	LogPath string `json:"log_path"`
	Topic   string `json:"topic"`
}

type LogInfo struct {
	AppId      int    `db:"app_id"`
	AppName    string `db:"app_name"`
	LogId      int    `db:"log_id"`
	Status     int    `db:"status"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
}

func InitEtcd(client *etcdclient.Client) {
	etcdClient = client
}

func GetAllLogInfo() (loglist []LogInfo, err error) {
	err = DB.Select(&loglist,
		"select a.app_id, b.app_name, a.create_time, a.topic, a.log_id, a.status, a.log_path from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func CreateLog(info *LogInfo) (err error) {

	conn, err := DB.Begin()
	if err != nil {
		logs.Warn("CreateApp failed, Db.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}

		conn.Commit()
	}()

	var appId []int
	err = DB.Select(&appId, "select app_id from tbl_app_info where app_name=?", info.AppName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}

	info.AppId = appId[0]
	r, err := conn.Exec("insert into tbl_log_info(app_id, log_path, topic)values(?, ?, ?)",
		info.AppId, info.LogPath, info.Topic)

	if err != nil {
		logs.Warn("CreateApp failed, Db.Exec error:%v", err)
		return
	}

	sd, err := r.LastInsertId()
	if err != nil {
		logs.Warn("CreateApp failed, Db.LastInsertId error:%v", err)
		return
	}

	fmt.Printf("sunccess : %d",sd)

	return
}

func SetLogConfToEtcd(etcdKey string, info *LogInfo) (err error) {

	var logConfArr []CollectConf
	logConfArr = append(
		logConfArr,
		CollectConf{
			LogPath: info.LogPath,
			Topic:   info.Topic,
		},
	)

	data, err := json.Marshal(logConfArr)
	if err != nil {
		logs.Warn("marshal failed, err:%v", err)
		return
	}


	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//cli.Delete(ctx, EtcdKey)
	//return
	_, err = etcdClient.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		logs.Warn("Put failed, err:%v", err)
		return
	}

	logs.Debug("put etcd succ, data:%v", string(data))
	return
}
