package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"news/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) GetRegister() {
	this.TplName = "register.html"
}
func (this *UserController) HandelRegister() {
	name := this.GetString("userName")
	pass := this.GetString("password")

	if name == "" || pass == "" {
		this.Data["errmsg"] = "用户名和密码不能为空!"
		this.TplName = "register.html"
		return
	}

	var u models.User
	u.Name = name
	u.Pass = pass
	n, err := orm.NewOrm().Insert(&u)
	if err != nil || n <= 0 {
		this.Data["errmsg"] = err.Error()
		this.TplName = "register.html"
		return
	}
	this.Redirect("/login", 302)

}

func (this *UserController) GetLogin() {
	name := this.Ctx.GetCookie("name")
	dec, _ := base64.StdEncoding.DecodeString(name)

	if name != "" {
		this.Data["name"] = string(dec)
		this.Data["checked"] = "checked"
	} else {
		this.Data["name"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"
}
func (this *UserController) HandelLogin() {
	name := this.GetString("userName")
	pass := this.GetString("password")

	if name == "" || pass == "" {
		this.Data["errmsg"] = "用户名和密码不能为空!"
		this.TplName = "login.html"
		return
	}
	var u models.User
	u.Name = name
	err := orm.NewOrm().Read(&u, "name")
	if err != nil {
		this.Data["errmsg"] = "用户名查找失败"
		this.TplName = "login.html"
		return
	}
	if u.Pass != pass {
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")
	if remember != "" {
		enc := base64.StdEncoding.EncodeToString([]byte(name))
		this.Ctx.SetCookie("name", enc, 60*1) //时间单位s
	} else {
		this.Ctx.SetCookie("name", "", -1)
	}

	this.SetSession("name", name)

	this.Redirect("/artical/mainPage", 302) //tpl渲染  redirect重定向跳转
}

func (this *UserController) ShowUserList() {
	name := this.GetSession("name")
	if name == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["name"] = name.(string)
	this.Layout = "layout.html"
	this.TplName = "userlist.html"
}
