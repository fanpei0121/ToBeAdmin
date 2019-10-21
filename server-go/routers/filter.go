package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"server-go/controllers"
	"time"
)

func init() {

	//允许跨域
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{}, //允许websocket跨域
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           time.Second * 600,
	}))
	// 权限验证
	var base controllers.BaseController
	beego.InsertFilter("/auth/rule", beego.BeforeRouter, base.FilterPermission("rule-view"))
	beego.InsertFilter("/auth/rule/:id", beego.BeforeRouter, base.FilterPermission("rule-update"))
	beego.InsertFilter("/auth/ruleAdd", beego.BeforeRouter, base.FilterPermission("rule-add"))
	beego.InsertFilter("/ruleDelete/:id", beego.BeforeRouter, base.FilterPermission("rule-delete"))
	// 角色
	beego.InsertFilter("/auth/role", beego.BeforeRouter, base.FilterPermission("role-view"))
	beego.InsertFilter("/auth/roleAdd", beego.BeforeRouter, base.FilterPermission("role-add"))
	beego.InsertFilter("/auth/role/:id", beego.BeforeRouter, base.FilterPermission("role-update"))
	beego.InsertFilter("/auth/roleDelete/:id", beego.BeforeRouter, base.FilterPermission("role-delete"))
	// 管理员
	beego.InsertFilter("/auth/user", beego.BeforeRouter, base.FilterPermission("account-view"))
	beego.InsertFilter("/auth/userAdd", beego.BeforeRouter, base.FilterPermission("account-add"))
	beego.InsertFilter("/auth/userDelete/:id", beego.BeforeRouter, base.FilterPermission("account-delete"))
	beego.InsertFilter("/auth/user/:id", beego.BeforeRouter, base.FilterPermission("account-update"))


}
