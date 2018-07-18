package router

import (
	"github.com/astaxie/beego"
	"web_logger/controller/appcontroller"
	"web_logger/controller/logcontroller"

)

func init() {
	beego.Router("/index", &appcontroller.AppController{}, `*:AppList`)
	beego.Router("/app/list", &appcontroller.AppController{}, `*:AppList`)
	beego.Router("/app/apply", &appcontroller.AppController{}, `*:AppApply`)
	beego.Router("/app/create", &appcontroller.AppController{}, `*:AppCreate`)





	beego.Router("/log/apply", &logcontroller.LogController{},`*:LogApply`)
	beego.Router("/log/list", &logcontroller.LogController{}, `*:LogList`)
	beego.Router("/log/create", &logcontroller.LogController{}, `*:LogCreate`)
}
