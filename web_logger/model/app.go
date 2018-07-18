package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/astaxie/beego/logs"
)

type AppInfo struct {
	AppId       int    `db:"app_id"`
	AppName     string `db:"app_name"`
	AppType     string `db:"app_type"`
	CreateTime  string `db:"create_time"`
	DevelopPath string `db:"develop_path"`
	IP          []string
}

var (
	DB *sqlx.DB
)

func InitDb(db *sqlx.DB) {
	DB = db
}

func GetAllAppInfo() (appinfoList []AppInfo, err error) {

	err = DB.Select(&appinfoList, "select app_id, app_name,app_type,create_time,develop_path from tbl_app_info")
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
	}
	return
}

func GetIPInfoById(appId int) (iplist []string, err error) {
	err = DB.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId)
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func GetIPInfoByName(appName string) (iplist []string, err error) {

	var appId []int
	err = DB.Select(&appId, "select app_id from tbl_app_info where app_name=?", appName)
	if err != nil || len(appId) == 0 {
		logs.Warn("select app_id failed, Db.Exec error:%v", err)
		return
	}

	err = DB.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId[0])
	if err != nil {
		logs.Warn("Get All App Info failed, err:%v", err)
		return
	}
	return
}

func CreateApp(info *AppInfo) (err error) {

	conn, err := DB.Begin()
	if err != nil {
		logs.Warn("CreateApp failed, DB.Begin error:%v", err)
		return
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}
		conn.Commit()
	}()

	r, err := conn.Exec("insert into tbl_app_info(app_name,app_type,develop_path)values (?,?,?)",
		info.AppName, info.AppType, info.DevelopPath)
	if err != nil {
		logs.Warn("CreateApp insert failed, DB.Exec error:%v", err)
		return
	}

	id, err := r.LastInsertId()
	if err != nil {
		logs.Warn("CreateApp GetLastInsertId failed, DB.LastInsertId error:%v", err)
		return
	}

	for _, ip := range info.IP {
		_, err = conn.Exec("insert into tbl_app_ip(app_id,ip)values (?,?)", id, ip)
		if err != nil {
			logs.Warn("CreateApp IP failed, DB.Exec error:%v", err)
			return
		}
	}
	return
}
