package controllers

import (
	"server-go/models"
	"server-go/server/jwtx"
)

type AuthController struct {
	BaseController
}

// login
func (t *AuthController) DoLogin() {
	param := make(map[string]string)
	t.GetJsonParam(&param)
	var authUser models.AuthUser
	tokenOrMsg, err := authUser.Check(param)
	if err != nil {
		t.Error(err.Error())
	}
	data := make(map[string]string)
	data["token"] = tokenOrMsg
	t.Success("success", data)
}

func (t *AuthController) Logout() {
	t.Success("退出成功")
}

type permissionsParam struct {
	Action          string                   `json:"action"`
	ActionEntitySet []map[string]interface{} `json:"actionEntitySet"`
	ActionList      interface{}              `json:"actionList"`
	Actions         []map[string]string      `json:"actions"`
	DataAccess      interface{}              `json:"dataAccess"`
	Id              int                      `json:"id"`
	PermissionId    string                   `json:"permissionId"`
	Title           string                   `json:"title"`
}

func (t *AuthController) Info() {
	tokenString := t.GetString("token")
	claims, err := jwtx.ParseToken(tokenString)
	if err != nil {
		t.Error("参数错误")
	}
	uid := int(claims["uid"].(float64))
	var authUser models.AuthUser
	userInfo := authUser.Info(uid)
	retData := make(map[string]interface{})
	retData["name"] = userInfo.Nickname
	retData["avatar"] = userInfo.Headimg
	retData["status"] = userInfo.Status
	retData["nickname"] = userInfo.Nickname
	retData["remark"] = userInfo.Remark
	retData["email"] = userInfo.Email
	var authPermission models.AuthPermission
	parentPermission := authPermission.GetParent()         // 所有顶级权限
	permissions := authPermission.GetPermissionsByUid(uid) // 获取拥有的子权限
	var permissionsParams []*permissionsParam
	for _, v := range parentPermission {
		actions := authPermission.GetTree(permissions, v.Id, 0)
		if len(actions) > 0 {
			var actionEntity []map[string]interface{}
			for _, vv := range actions {
				actionEntity = append(actionEntity, map[string]interface{}{
					"action":       vv["action"],
					"defaultCheck": false,
					"describe":     vv["title"],
				})
			}
			permissionsParams = append(permissionsParams, &permissionsParam{
				Action:          v.Action,
				ActionEntitySet: actionEntity,
				ActionList:      nil,
				Actions:         actions,
				DataAccess:      nil,
				Id:              v.Id,
				PermissionId:    v.Action,
				Title:           v.Title,
			})
		}
	}
	per := make(map[string]interface{})
	per["permissions"] = permissionsParams
	retData["role"] = per
	t.Success("success", retData)
}

// 账户设置
func (t *AuthController) Setting() {
	param := make(map[string]string)
	t.GetJsonParam(&param)
	userId := t.GetUidByHead()
	var userModel models.AuthUser
	//if userId == userModel.GetSupperId() {
	//	t.Error("预览版，不能修改超级管理员资料")
	//}
	err := userModel.Setting(userId, param)
	if err != nil {
		t.Error(err.Error())
	}
	t.Success("设置成功")
}

