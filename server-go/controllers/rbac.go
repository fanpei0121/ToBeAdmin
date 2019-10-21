package controllers

import (
	"server-go/models"
)

type RbacController struct {
	BaseController
}

var authPermissionModel models.AuthPermission

// 规则列表
func (t *RbacController) Rules() {
	page, _ := t.GetInt("page", 1)
	pageSize, _ := t.GetInt("pageSize", 10)
	var authPermission models.AuthPermission
	list := authPermission.GetAllPermissionsTree(page, pageSize)
	t.Success("操作成功", list)
}

// 添加规则
func (t *RbacController) RuleAdd() {
	param := make(map[string]interface{})
	t.GetJsonParam(&param)
	var authPermission models.AuthPermission
	err := authPermission.Add(param)
	if err != nil {
		t.Error("新增失败")
	} else {
		t.Success("新增成功")
	}

}

// 修改规则
func (t *RbacController) RuleEdit() {
	param := make(map[string]interface{})
	t.GetJsonParam(&param)
	id, err := t.GetInt(":id", 0)
	if err != nil {
		t.Error("paramError")
	}
	var authPermission models.AuthPermission
	err = authPermission.Edit(id, param)
	if err != nil {
		t.Error("编辑失败")
	} else {
		t.Success("编辑成功")
	}
}

// 删除
func (t *RbacController) RuleDelete() {
	id, err := t.GetInt(":id", 0)
	if err != nil {
		t.Error("paramError")
	}
	err = authPermissionModel.Delete(id)
	if err != nil {
		t.Error("删除成功")
	}
	t.Success("删除成功")
}

type treeParam struct {
	models.AuthPermission
	Cname string `json:"cname"`
}

// 顶级权限
func (t *RbacController) Tree() {
	var permission models.AuthPermission
	permissions := permission.GetParent()
	var treeParams []treeParam
	for _, v := range permissions {
		treeParams = append(treeParams, treeParam{
			AuthPermission: models.AuthPermission{
				Id:     v.Id,
				Action: v.Action,
				Name:   v.Name,
				Pid:    v.Pid,
				Status: v.Status,
				Title:  v.Title,
			},
			Cname: v.Title,
		})
	}
	t.Success("success", treeParams)
}

var authRoleModel models.AuthRole

//角色列表
func (t *RbacController) Roles() {
	var err error
	page, err := t.GetInt("page", 1)
	pageSize, err := t.GetInt("pageSize", 10)
	if err != nil {
		t.Error("参数错误")
	}
	roles, err := authRoleModel.Roles(page, pageSize)
	if err != nil {
		panic(err)
	}
	t.Success("success", roles)
}

// 角色添加
func (t *RbacController) RoleAdd() {
	param := make(map[string]interface{})
	t.GetJsonParam(&param)
	var roleModel models.AuthRole
	err := roleModel.Add(param)
	if err == nil {
		t.Success("新增成功")
	} else {
		t.Error("新增失败, 唯一标识重复")
	}
}

// 角色修改
func (t *RbacController) RoleEdit() {
	roleId, err := t.GetInt(":id", 0)
	if err != nil || roleId == 0 {
		t.Error("paramError")
	}
	param := make(map[string]interface{})
	t.GetJsonParam(&param)
	err = authRoleModel.Edit(roleId, param)
	if err != nil {
		t.Error("修改失败, 唯一标识重复")
	}
	t.Success("修改成功")
}

// 角色删除
func (t *RbacController) RoleDelete() {
	roleId, err := t.GetInt(":id", 0)
	if err != nil || roleId == 0 {
		t.Error("paramError")
	}
	err = authRoleModel.Delete(roleId)
	if err != nil {
		t.Error("删除失败")
	}
	t.Success("删除成功")
}

var userModel models.AuthUser

// 管理员列表
func (t *RbacController) Users() {
	page, err := t.GetInt("page", 1)
	pageSize, err := t.GetInt("pageSize", 10)
	if err != nil {
		t.Error("paramError")
	}
	retData := make(map[string]interface{})
	data, err := authRoleModel.Roles(1, 1000)
	if err != nil {
		panic(err)
	}
	retData["roles"] = data["roles"]
	retData["rules"] = data["rules"]
	users, err := userModel.Users(page, pageSize)
	if err != nil {
		panic(err)
	}
	retData["users"] = users
	t.Success("success", retData)
}

// 管理员添加
func (t *RbacController) UserAdd() {
	param := make(map[string]interface{})
	t.GetJsonParam(&param)
	err := userModel.Add(param)
	if err != nil {
		panic(err)
		t.Error("新增失败")
	}
	t.Success("新增成功")
}

// 管理员删除
func (t *RbacController) UserDelete() {
	userId, err := t.GetInt(":id", 0)
	if err != nil || userId == 0 {
		t.Error("paramError")
	}
	err = userModel.Delete(userId)
	if err != nil {
		t.Error("删除失败")
	}
	t.Success("删除成功")
}

// 管理员修改
func (t *RbacController) UserEdit() {
	userId, err := t.GetInt(":id", 0)
	if err != nil || userId == 0 {
		t.Error("paramError")
	}
	param := make(map[string]interface{})
	t.GetJsonParam(&param)
	err = userModel.Edit(userId, param)
	if err != nil {
		panic(err)
		t.Error("修改失败")
	}
	t.Success("修改成功")
}
