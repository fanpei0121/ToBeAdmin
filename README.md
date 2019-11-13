## beego Ant Design Pro Vue RBAC

[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)
beego 与 Ant Design Pro Vue 基础权限系统  

后端: http://beego.me 

ORM: gorm https://gorm.io/zh_CN/docs/

Ant Design Pro Vue文档: https://pro.loacg.com/docs/getting-started

**预览:**

http://www.micro123.cn



![1.png](https://image.jnemall.com/uploads/attachment/20191021/4487711810264ac99177ba5eb9f7e6b3.jpg)

![2.png](https://image.jnemall.com/uploads/attachment/20191021/676656cd0bd2536dbca48c8bf8a1c988.png)

![3.png](https://image.jnemall.com/uploads/attachment/20191021/3d8028a2f9f9c0606445cfdf8572da09.png)

### 账号
* 超级管理员 `admin, 1234` 
* 普通管理员 `test, 1234`

### 前端部署
#### 安装
```
cd client
npm install
```
#### 预览
```
npm run serve

浏览器访问 localhost:8000
```
#### 打包
```
npm run build
```
#### 刷新404问题
```
nginx 加上几行配置

location / {
    try_files $uri $uri/ /index.html last;
}
```
#### 后端API配置
/client/src/config/defaultSettings.js  修改 baseURL

### go部署

#### mysql
* mysql导入sql文件  /server-go/sql/table.sql 
* 修改 /server-go/conf/app.conf 数据库配置

#### 运行

```
cd server-go
bee run 
```

#### 编译
```
bee pack -be GOOS=linux
```

#### 设置权限
```
/routers/filter.go 文件设置权限

beego.InsertFilter("/auth/rule", beego.BeforeRouter, base.FilterPermission("rule-view"))

/auth/rule 为路由
rule-view 为后台添加的权限规则
```

