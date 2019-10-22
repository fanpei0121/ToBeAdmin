package models

import (
	"fmt"
	"github.com/astaxie/beego"
)

type AuthRole struct {
	Id     int    `gorm:"primary_key" json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}

// 根据管理员id获取角色
func (t *AuthRole) GetRolesByUid(uid int) []int {
	var roles []int
	err := db.Table("auth_user_role_access").Where("user_id = ?", uid).Pluck("role_id", &roles).Error
	if err != nil {
		panic(err)
	}
	return roles
}

type listParam struct {
	AuthRole
	Permissions []int `json:"permissions"`
}

// 角色列表
func (t *AuthRole) Roles(page int, pageSize int) (map[string]interface{}, error) {
	var authRoles []AuthRole
	retData := make(map[string]interface{})
	offset := (page - 1) * pageSize
	var err error
	db.Table("auth_role").Offset(offset).Limit(pageSize).Order("id desc").Find(&authRoles)
	var listParams []listParam
	for _, v := range authRoles {
		var ids []int
		err = db.Table("auth_role_permission_access").Where("role_id = ?", v.Id).Pluck("permission_id", &ids).Error
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		listParams = append(listParams, listParam{
			AuthRole: AuthRole{
				Id:     v.Id,
				Name:   v.Name,
				Status: v.Status,
				Title:  v.Title,
			},
			Permissions: ids,
		})
	}
	pagination := make(map[string]int)
	pagination["current"] = page
	pagination["pageSize"] = pageSize
	var total int
	db.Model(&AuthRole{}).Count(&total)
	pagination["total"] = total
	retData["data"] = listParams
	retData["pagination"] = pagination
	warm := make(map[string]interface{})
	warm["roles"] = retData
	// 规则
	var permissionModel AuthPermission
	permissions := permissionModel.GetAllPermissionsTree(1, 1000)
	warm["rules"] = permissions
	return warm, nil
}

// 角色添加
func (t *AuthRole) Add(param map[string]interface{}) error {
	role := AuthRole{
		Name:   param["name"].(string),
		Title:  param["title"].(string),
		Status: int(param["status"].(float64)),
	}
	if err := db.Create(&role).Error; err != nil {
		beego.Error(err)
		return err
	}
	permissionIds := param["rules"].([]interface{})
	sql := "INSERT INTO `auth_role_permission_access` (`role_id`,`permission_id`) VALUES"
	for key, value := range permissionIds {
		if len(permissionIds)-1 == key {
			//最后一条数据以分号结尾
			sql += fmt.Sprintf("(%d, %d);", role.Id, int(value.(float64)))
		} else {
			sql += fmt.Sprintf("(%d,%d),",  role.Id, int(value.(float64)))
		}
	}
	db.Exec(sql)
	return nil
}

// 角色权限关联模型
type AuthRolePermissionAccess struct {
	RoleId       int `gorm:"primary_key"`
	PermissionId int `json:"permission_id"`
}

// 角色修改
func (t *AuthRole) Edit(roleId int, param map[string]interface{}) error {
	var role AuthRole
	db.Where("id = ?", roleId).First(&role)
	err := db.Model(&role).UpdateColumns(AuthRole{
		Name:   param["name"].(string),
		Title:  param["title"].(string),
		Status: int(param["status"].(float64)),
	}).Error
	if err != nil {
		beego.Error(err)
		return err
	}
	// 删除关联
	db.Where("role_id = ?", roleId).Delete(&AuthRolePermissionAccess{})
	//新增关联
	permissionIds := param["rules"].([]interface{})
	sql := "INSERT INTO `auth_role_permission_access` (`role_id`,`permission_id`) VALUES"
	for key, value := range permissionIds {
		if len(permissionIds)-1 == key {
			//最后一条数据以分号结尾
			sql += fmt.Sprintf("(%d, %d);", roleId, int(value.(float64)))
		} else {
			sql += fmt.Sprintf("(%d,%d),", roleId, int(value.(float64)))
		}
	}
	db.Exec(sql)
	return nil
}

// 角色删除
func (t *AuthRole) Delete(roleId int) error {
	tx := db.Begin()

	err := db.Where("id = ?", roleId).Delete(&AuthRole{}).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	err = db.Where("role_id = ?", roleId).Delete(&AuthRolePermissionAccess{}).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
