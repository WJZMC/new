package controllers

import (
	"github.com/astaxie/beego"
	//"news/models"
)

type ArticalController struct {
	beego.Controller
}

//@router  /ShowArticalList [get]
func (this *ArticalController) ShowArticalList() {
	//this.TplName = "add.html"
}

//@router /AddArticalRouter [post]
func (this *ArticalController) AddArtical() {
	//artical := models.Aritcal{}
	//if err := this.ParseForm(&artical); err != nil {
	//	beego.Error(err)
	//	return
	//}
	//beego.Info(artical)
}
