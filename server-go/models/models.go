package models

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//gorm全局变量
var db *gorm.DB

func init() {
	RegisterGormDB()
}



//gorm创建
func RegisterGormDB() {
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqldb := beego.AppConfig.String("mysqldb")
	host := beego.AppConfig.String("host")
	port := beego.AppConfig.String("port")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqluser, mysqlpass, host, port, mysqldb)
	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if beego.AppConfig.DefaultBool("ormdebug", false) { // 启用Logger，显示详细日志
		db.LogMode(true)
	}
	// 全局禁用表名复数
	db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
	//您可以通过定义DefaultTableNameHandler对默认表名应用任何规则。
	/*gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return beego.AppConfig.DefaultString("prefix", "cod_") + defaultTableName
	}*/
}
