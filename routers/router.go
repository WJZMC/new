package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"news/controllers"
)

func init() {
	beego.InsertFilter("/artical/*", beego.BeforeExec, fileterLogin)
	//beego.Include(&controllers.ArticalController{})
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:GetRegister;post:HandelRegister")
	beego.Router("/login", &controllers.UserController{}, "get:GetLogin;post:HandelLogin")
	beego.Router("/artical/mainPage", &controllers.MainPageController{}, "get:GetMainPage")
	beego.Router("/artical/add_Artical", &controllers.MainPageController{}, "get:ShowAddArtical;post:HandelAddArtical")
	beego.Router("/artical/aritcalDetail", &controllers.MainPageController{}, "get:ShowAritcalDetail")
	beego.Router("/artical/updateArtical", &controllers.MainPageController{}, "get:ShowEditArtical;post:HandelEditArtical")
	beego.Router("/artical/delArtical", &controllers.MainPageController{}, "get:HandelDeleteArtical")
	beego.Router("/artical/addArticalType", &controllers.MainPageController{}, "get:ShowAddArticalType;post:HandelAddArticalType")
	beego.Router("/artical/deleteType", &controllers.MainPageController{}, "get:HandelDeleteArticalType")
	beego.Router("/artical/logout", &controllers.MainPageController{}, "get:HandelLogout")
	beego.Router("artical/userList", &controllers.UserController{}, "get:ShowUserList")
}

func fileterLogin(ctx *context.Context) {
	name := ctx.Input.Session("name")
	if name == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
