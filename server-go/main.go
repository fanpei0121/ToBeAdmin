package main

import (
	"github.com/astaxie/beego"
	_ "server-go/routers"
)

func main() {
	beego.SetLogger("file", `{"filename":"logs/logs.log"}`) // 日志记录
	beego.SetLevel(beego.LevelWarning) // 日志级别
	beego.BeeLogger.DelLogger("console") // 取消日志console输出
	beego.Run()
}
