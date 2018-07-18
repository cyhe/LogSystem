package appcontroller

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/logs"
	"web_logger/model"
	"strings"
)

type AppController struct {
	beego.Controller
}

func (p *AppController) AppList() {

	p.Layout = "layout/layout.html"

	applist, err := model.GetAllAppInfo()
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("服务器繁忙")
		p.TplName = "app/error.html"
		logs.Warn("get app list failed, err:%v",err)
		return
	}

	logs.Debug("get app list succ, data:%v",applist)
	p.Data["applist"] = applist

	p.TplName = "app/index.html"
}

func (p *AppController) AppApply() {
	p.Layout = "layout/layout.html"
	p.TplName = "app/apply.html"
}

func (p *AppController) AppCreate() {
	appName := p.GetString("app_name")
	appTyoe := p.GetString("app_type")
	devLogPath := p.GetString("develop_path")
	ipstr := p.GetString("iplist")

	p.Layout = "layout/layout.html"

	if len(appName) == 0 || len(appTyoe) == 0 || len(devLogPath) == 0 || len(ipstr) == 0 {
		p.Data["Error"] = fmt.Sprintf("非法参数")
		p.TplName = "app/error.html"
		logs.Warn("invalid parameter")
		return
	}

	appinfo := &model.AppInfo{
		AppName:     appName,
		AppType:     appTyoe,
		DevelopPath: devLogPath,
		IP:          strings.Split(ipstr, ","),
	}

	err := model.CreateApp(appinfo)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("项目创建失败,数据库繁忙")
		p.TplName = "app/error.html"
		logs.Warn("invalid parameter")
		return
	}

	p.Layout = "layout/layout.html"

	p.Redirect("/app/list", 302)
}
