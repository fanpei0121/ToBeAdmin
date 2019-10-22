package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"server-go/models"
	"server-go/server/jwtx"
	"strconv"
	"time"
)

type BaseController struct {
	beego.Controller
}

func (t *BaseController) Success(message string, data ...interface{}) {
	ret := map[string]interface{}{"code": 20000, "message": message}
	if len(data) > 0 {
		ret["data"] = data[0]
	} else {
		ret["data"] = ""
	}
	t.Data["json"] = ret
	t.ServeJSON()
	t.StopRun()
}

func (t *BaseController) Error(message string) {
	ret := map[string]interface{}{"code": 50015, "message": message}
	t.Data["json"] = ret
	t.ServeJSON()
	t.StopRun()
}

//获取请求json参数
func (t *BaseController) GetJsonParam(param interface{}) {
	err := json.Unmarshal(t.Ctx.Input.RequestBody, param)
	if err != nil {
		beego.Error(err)
		t.Error("ParametersError")
	}
}

//权限验证
func (t *BaseController) FilterPermission(permission string) beego.FilterFunc {
	return func(ctx *context.Context) {
		tokenString := ctx.Input.Header("Authorization")
		claims, err := jwtx.ParseToken(tokenString)
		retData := make(map[string]interface{})
		retData["code"] = 50015
		if err != nil {
			beego.Error(err)
			retData["message"] = "ParametersError"
			ctx.Output.JSON(retData, true, true)
		}
		uid := int(claims["uid"].(float64))
		var authPermission models.AuthPermission
		permissions := authPermission.GetPermissionsByUid(uid)
		for _, v := range permissions {
			if v.Name == permission {
				return
			}
		}
		retData["message"] = "没有权限"
		ctx.Output.JSON(retData, true, true)
	}
}

// 上传
func (t *BaseController) Upload() {
	f, h, err := t.GetFile("file")
	if err != nil {
		beego.Error(err)
		return
	}
	defer f.Close()
	fileName := strconv.Itoa(int(time.Now().Unix())) + h.Filename
	path := "static/upload/" + fileName
	err = t.SaveToFile("file", path)
	if err != nil {
		beego.Error(err)
		return
	}
	param := make(map[string]string)
	param["url"] = path
	t.Success("上传成功", param)
}

// 根据head获取user_Id
func (t *BaseController) GetUidByHead() int {
	tokenString := t.Ctx.Input.Header("Authorization")
	if tokenString == "" {
		return 0
	}
	claims, err := jwtx.ParseToken(tokenString)
	if err != nil {
		beego.Error(err)
		t.Error("参数错误")
	}
	uid := int(claims["uid"].(float64))
	return uid
}
