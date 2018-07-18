package main

import (
	"github.com/astaxie/beego"
	_ "web_logger/router"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"web_logger/model"
	"time"
	etvd_client "github.com/coreos/etcd/clientv3"
)

func initDb() (err error) {

	database, err := sqlx.Open("mysql", "root:12345678@tcp(localhost:3306)/web_logger")
	if err != nil {
		logs.Warn("open mysql failed,", err)
	}

	// 传给model
	model.InitDb(database)

	return
}

func initEtcd() (err error) {
	cli, err := etvd_client.New(etvd_client.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logs.Warn("conn etcd failed,", err)
		return
	}
	//传给model
	model.InitEtcd(cli)
	return
}

func main() {

	err := initDb()
	if err != nil {
		logs.Warn("init db failed, err:%v", err)
		return
	}

	err = initEtcd()
	if err != nil {
		logs.Warn("init etcd failed, err:%v", err)
		return
	}

	beego.Run()
}
