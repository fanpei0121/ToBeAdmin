package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"server-go/helper"
	"server-go/server/jwtx"
)

type AuthUser struct {
	Id       int    `gorm:"primary_key" json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	Headimg  string `json:"headimg"`
	Remark   string `json:"remark"`
	Email    string `json:"email"`
}

// 获得超级管理员id
func (t *AuthUser) GetSupperId() int {
	return 1
}

// 账号检测
func (t *AuthUser) Check(param map[string]string) (string, error) {
	var user AuthUser
	if err := db.Where("name = ?", param["username"]).First(&user).Error; err != nil {
		beego.Error(err)
		return "", err
	}
	if user.Status != 1 {
		return "", errors.New("账号被禁用")
	}
	if !helper.PasswordVerify(param["password"], user.Password) {
		return "", errors.New("账号或密码不正确")
	}
	jwtParam := make(map[string]interface{})
	jwtParam["uid"] = user.Id
	tokenString, err := jwtx.GenToken(jwtParam)
	if err != nil {
		beego.Error(err)
		return "", err
	}
	return tokenString, nil
}

// 用户信息
func (t *AuthUser) Info(id int) AuthUser {
	var user AuthUser
	db.Where("id = ?", id).First(&user)
	return user
}

type userListParam struct {
	AuthUser
	Rules []int `json:"rules"`
	Roles []int `json:"roles"`
}

// 管理员列表
func (t *AuthUser) Users(page int, pageSize int) (map[string]interface{}, error) {
	var users []AuthUser
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Where("id <> ?", t.GetSupperId()).Order("id desc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	var listParams []userListParam
	for _, v := range users {
		var roleIds []int
		db.Table("auth_user_role_access").Where("user_id = ?", v.Id).Pluck("role_id", &roleIds)
		var ruleIds []int
		db.Table("auth_role_permission_access").Where("role_id in (?)", roleIds).Pluck("permission_id", &ruleIds)
		listParams = append(listParams, userListParam{
			AuthUser: AuthUser{
				Id:       v.Id,
				Name:     v.Name,
				Nickname: v.Nickname,
				Status:   v.Status,
			},
			Rules: ruleIds,
			Roles: roleIds,
		})
	}
	pagination := make(map[string]int)
	pagination["current"] = page
	pagination["pageSize"] = pageSize
	var total int
	db.Model(&AuthUser{}).Where("id <> ?", t.GetSupperId()).Count(&total)
	pagination["total"] = total
	retData := make(map[string]interface{})
	retData["data"] = listParams
	retData["pagination"] = pagination
	return retData, nil
}

// 管理员添加
func (t *AuthUser) Add(param map[string]interface{}) error {
	newPass, err := helper.PasswordHash(param["password"].(string))
	roles := param["roles"].([]interface{})
	if err != nil {
		beego.Error(err)
		return err
	}
	user := AuthUser{
		Name:     param["name"].(string),
		Nickname: param["nickname"].(string),
		Password: newPass,
		Status:   int(param["status"].(float64)),
	}
	err = db.Create(&user).Error
	if err != nil {
		beego.Error(err)
		return err
	}
	sql := "INSERT INTO `auth_user_role_access` (`user_id`,`role_id`) VALUES"
	for key, value := range roles {
		if len(roles)-1 == key {
			sql += fmt.Sprintf("(%d, %s);", user.Id, value.(string))
		} else {
			sql += fmt.Sprintf("(%d, %s),", user.Id, value.(string))
		}
	}
	err = db.Exec(sql).Error
	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

type AuthUserRoleAccess struct {
	UserId int `json:"user_id"`
	RoleId int `json:"role_id"`
}

// 删除
func (t *AuthUser) Delete(userId int) error {
	tx := db.Begin()
	err := db.Where("id = ?", userId).Delete(&AuthUser{}).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	err = db.Where("user_id = ?", userId).Delete(&AuthUserRoleAccess{}).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 管理员修改
func (t *AuthUser) Edit(userId int, param map[string]interface{}) error {
	tx := db.Begin()
	var user AuthUser
	err := tx.Where("id = ?", userId).First(&user).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	editUser := AuthUser{
		Name:     param["name"].(string),
		Nickname: param["nickname"].(string),
		Status:   int(param["status"].(float64)),
	}
	password, ok := param["password"].(string)
	if ok {
		newPass, err := helper.PasswordHash(password)
		if err != nil {
			beego.Error(err)
			return err
		}
		editUser.Password = newPass
	}
	err = tx.Model(&user).UpdateColumns(editUser).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	err = tx.Where("user_id = ?", userId).Delete(&AuthUserRoleAccess{}).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	roleIds := param["roles"].([]interface{})
	sql := "INSERT INTO `auth_user_role_access` (`user_id`,`role_id`) VALUES"
	for key, value := range roleIds {
		if len(roleIds)-1 == key {
			//最后一条数据以分号结尾
			sql += fmt.Sprintf("(%d, %s);", userId, value.(string))
		} else {
			sql += fmt.Sprintf("(%d, %s),", userId, value.(string))
		}
	}
	err = tx.Exec(sql).Error
	if err != nil {
		beego.Error(err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 账户设置
func (t *AuthUser) Setting(userId int, param map[string]string) error {
	var user AuthUser
	db.Where("id = ?", userId).First(&user)
	updateUser := AuthUser{
		Nickname: param["nickname"],
		Headimg:  param["headimg"],
		Remark:   param["remark"],
		Email:    param["email"],
	}
	if param["password"] != "" {
		if param["password"] != param["password2"] {
			return errors.New("密码输入不一致")
		}
		jmPass, err := helper.PasswordHash(param["password"])
		if err != nil {
			beego.Error(err)
			return err
		}
		updateUser.Password = jmPass
	}
	err := db.Model(&user).Updates(updateUser).Error
	if err != nil {
		return err
	}
	return nil
}
