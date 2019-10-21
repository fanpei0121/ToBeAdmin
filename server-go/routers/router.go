package routers

import (
	"github.com/astaxie/beego"
	"server-go/controllers"
)

func init() {
	user := beego.NewNamespace("user",
		beego.NSRouter("/login", &controllers.AuthController{}, "post:DoLogin"),
		beego.NSRouter("/info", &controllers.AuthController{}, "get:Info"),
	)
	beego.Router("/auth/logout", &controllers.AuthController{}, "post:Logout")
	beego.Router("/upload", &controllers.BaseController{}, "post:Upload")
	beego.Router("/settings", &controllers.AuthController{}, "post:Setting") // 修改资料

	rbac := beego.NewNamespace("auth",
		beego.NSBefore(),
		// 规则
		beego.NSRouter("/rule", &controllers.RbacController{}, "get:Rules"),
		beego.NSRouter("/ruleAdd", &controllers.RbacController{}, "post:RuleAdd"),
		beego.NSRouter("/rule/:id", &controllers.RbacController{}, "put:RuleEdit"),
		beego.NSRouter("/ruleDelete/:id", &controllers.RbacController{}, "delete:RuleDelete"),
		beego.NSRouter("/tree", &controllers.RbacController{}, "get:Tree"),
		// 角色
		beego.NSRouter("/role", &controllers.RbacController{}, "get:Roles"),
		beego.NSRouter("/roleAdd", &controllers.RbacController{}, "post:RoleAdd"),
		beego.NSRouter("/role/:id", &controllers.RbacController{}, "put:RoleEdit"),
		beego.NSRouter("/roleDelete/:id", &controllers.RbacController{}, "delete:RoleDelete"),
		// 管理员
		beego.NSRouter("/user", &controllers.RbacController{}, "get:Users"),
		beego.NSRouter("/userAdd", &controllers.RbacController{}, "post:UserAdd"),
		beego.NSRouter("/userDelete/:id", &controllers.RbacController{}, "delete:UserDelete"),
		beego.NSRouter("/user/:id", &controllers.RbacController{}, "put:UserEdit"),
	)
	beego.AddNamespace(user, rbac)



}
