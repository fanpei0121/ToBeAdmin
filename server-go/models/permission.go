package models

import "github.com/astaxie/beego"

type AuthPermission struct {
	Id     int    `gorm:"primary_key" json:"id"`
	Action string `json:"action"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Pid    int    `json:"pid"`
	Status int    `json:"status"`
}

//根据角色ids获取子权限ids
func (t *AuthPermission) GetPermissionsByRoleIds(roleIds []int) []int {
	var permissionIds []int
	db.Table("auth_role_permission_access").Where("role_id in (?)", roleIds).Pluck("permission_id", &permissionIds)
	return permissionIds
}

//根据用户id获取子权限
func (t *AuthPermission) GetPermissionsByUid(uid int) []*AuthPermission {
	var authRole AuthRole
	var authUser AuthUser
	var authPermission AuthPermission
	var permissionIds []int
	if authUser.GetSupperId() == uid { //超级管理员
		permissionIds = authPermission.GetAllPermissions() //获取所有子权限
	} else {
		roles := authRole.GetRolesByUid(uid)                          // 角色列表
		permissionIds = authPermission.GetPermissionsByRoleIds(roles) // 权限ids
	}
	permissions := authPermission.GetPermissions(permissionIds) // 权限列表
	return permissions
}

// 获取所有子权限
func (t *AuthPermission) GetAllPermissions() []int {
	var permissionIds []int
	db.Table("auth_permission").Where("pid <> ?", 0).Pluck("id", &permissionIds)
	return permissionIds
}

//根据权限id获取权限
func (t *AuthPermission) GetPermissions(permissionIds []int) []*AuthPermission {
	var permissions []*AuthPermission
	db.Where("id in (?)", permissionIds).Find(&permissions)
	return permissions
}

// 获取所有顶级权限
func (t *AuthPermission) GetParent() []*AuthPermission {
	var permissions []*AuthPermission
	db.Where("pid = ?", 0).Find(&permissions)
	return permissions
}

//获取父级权限下的子权限
func (t *AuthPermission) GetTree(permissions []*AuthPermission, id int, step int) []map[string]string {
	var data []map[string]string
	step++
	for _, v := range permissions {
		if step > 2 { //最多2层递归 避免栈内存溢出
			break
		}
		if v.Pid == id {
			data = append(data, map[string]string{
				"action": v.Action,
				"title":  v.Title,
			})
			data = append(data, t.GetTree(permissions, v.Id, step)...)
		}
	}
	return data
}

// 获取所有规则并格式化
func (t *AuthPermission) GetAllPermissionsTree(page int, pageSize int) map[string]interface{} {
	var permissions []*AuthPermission
	retData := make(map[string]interface{})
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Where("pid = ?", 0).Order("id desc").Find(&permissions) // 顶级权限
	pagination := make(map[string]int)
	pagination["current"] = page
	pagination["pageSize"] = pageSize
	var total int
	db.Model(&AuthPermission{}).Where("pid = ?", 0).Count(&total)
	pagination["total"] = total
	retData["data"] = t.TreeNode(permissions)
	retData["pagination"] = pagination
	return retData
}

type perData struct {
	AuthPermission
	Actions      []*NodeChildrenParam `json:"actions"`
	PermissionId string               `json:"permissionId"`
}

// 权限列表组装数据
func (t *AuthPermission) TreeNode(permissions []*AuthPermission) []perData {
	var perDatas []perData
	for _, v := range permissions {
		actions := t.TreeNodeChildren(v.Id)
		perDatas = append(perDatas, perData{
			AuthPermission: AuthPermission{
				Id:     v.Id,
				Action: v.Action,
				Name:   v.Name,
				Title:  v.Title,
				Pid:    v.Pid,
				Status: v.Status,
			},
			Actions:      actions,
			PermissionId: v.Action,
		})
	}
	return perDatas
}

type NodeChildrenParam struct {
	AuthPermission
	Label string `json:"label"`
	Value int    `json:"value"`
}

func (t *AuthPermission) TreeNodeChildren(pid int) []*NodeChildrenParam {
	var data []*NodeChildrenParam
	db.Table("auth_permission").Where("pid = ?", pid).Find(&data)
	for k, v := range data {
		data[k].Label = v.Title
		data[k].Value = v.Id
	}
	return data
}

// 添加权限
func (t *AuthPermission) Add(param map[string]interface{}) error {
	authPermission := AuthPermission{
		Action: param["action"].(string),
		Name:   param["name"].(string),
		Title:  param["title"].(string),
		Pid:    int(param["pid"].(float64)),
		Status: int(param["status"].(float64)),
	}
	if err := db.Create(&authPermission).Error; err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

// 编辑
func (t *AuthPermission) Edit(id int, param map[string]interface{}) error {
	var authPermission AuthPermission
	db.Where("id = ? ", id).First(&authPermission)
	err := db.Model(&authPermission).Updates(param).Error;
	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func (t *AuthPermission) Delete(id int) error {
	var authPermission AuthPermission
	db.Where("id = ? ", id).First(&authPermission)
	err := db.Delete(&authPermission).Error;
	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}
