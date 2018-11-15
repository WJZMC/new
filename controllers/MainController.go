package controllers

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"news/models"
	"path"
	"strconv"
	"time"
)

type MainPageController struct {
	beego.Controller
}

func (this *MainPageController) GetMainPage() {

	o := orm.NewOrm()

	//文章类型
	var ats []models.ArticalType
	o.QueryTable("ArticalType").All(&ats)

	selectName := this.GetString("select")
	if selectName == "" {
		selectName = this.GetString("selectName")
		if selectName == "" {
			if len(ats) > 0 {
				artType := ats[0]
				selectName = artType.TypeName
			}
		}
	}
	this.Data["selectName"] = selectName

	var articles []models.Artical

	qs := o.QueryTable("Artical")

	articalCount, _ := qs.RelatedSel("ArticalType").Filter("ArticalType__TypeName", selectName).Count()

	pageSize := 20

	pageCount := math.Ceil(float64(articalCount) / float64(pageSize))

	curentPage, _ := this.GetInt("pageIndex")
	if curentPage <= 0 {
		curentPage = 1
	}
	startIndex := pageSize * (curentPage - 1)

	//此处RelatedSel 中的参数是model中的一对多的“变量名”,如果“AritcalType *ArticalType `orm:"rel(fk)"`”
	qs.Limit(pageSize, startIndex).RelatedSel("ArticalType").Filter("ArticalType__TypeName", selectName).RelatedSel("Author").All(&articles)

	//beego.Info(articles)
	this.Data["curentPage"] = curentPage
	this.Data["pageCount"] = pageCount
	this.Data["articalCount"] = articalCount

	this.Data["articals"] = articles

	this.Data["errmsg"] = this.GetString("errmsg")

	this.Data["ats"] = ats

	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}
	this.Layout = "layout.html"

	this.TplName = "index.html"
}
func (this *MainPageController) ShowAddArtical() {
	//文章类型
	var ats []models.ArticalType
	orm.NewOrm().QueryTable("ArticalType").All(&ats)
	this.Data["ats"] = ats

	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

func (this *MainPageController) HandelAddArtical() {

	articleName := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	if articleName == "" || content == "" {
		errmsg := "文章标题和文章内容不能为空"
		this.Redirect("/artical/add_Artical?errmsg="+errmsg, 302)
		return
	}
	if typeName == "" {
		errmsg := "类型不能为空"
		this.Redirect("/artical/add_Artical?errmsg="+errmsg, 302)
		return
	}

	fileUrl, err := this.uploadFile("uploadname")
	if err != nil {
		errmsg := err.Error()
		this.Redirect("/artical/add_Artical?errmsg="+errmsg, 302)
		return
	}
	o := orm.NewOrm()

	var artical models.Artical
	artical.Title = articleName
	artical.Content = content
	artical.ImageUrl = fileUrl

	var articalType models.ArticalType
	articalType.TypeName = typeName
	o.Read(&articalType, "TypeName")
	artical.ArticalType = &articalType

	var user models.User
	user.Name = this.GetSession("name").(string)
	o.Read(&user, "Name")
	artical.Author = &user

	n, err := o.Insert(&artical)
	if err != nil || n <= 0 {
		beego.Error(err)
		errmsg := "文章添加失败"
		this.Redirect("/artical/add_Artical?errmsg="+errmsg, 302)
		return
	}
	this.Redirect("/artical/mainPage", 302)
}

func (this *MainPageController) ShowAritcalDetail() {
	id, err := this.GetInt("Id")

	if err != nil {
		errmsg := "获取文章详情参数错误"
		this.Redirect("/artical/mainPage?errmsg="+errmsg, 302)
		return
	}
	var artical models.Artical
	artical.Id = id
	err = orm.NewOrm().Read(&artical)
	if err != nil {
		errmsg := "获取文章详情失败"
		this.Redirect("/artical/mainPage?errmsg="+errmsg, 302)
		return
	}
	artical.ReadCount++

	o := orm.NewOrm()
	o.Update(&artical)

	session := this.GetSession("name")
	if session != nil {
		name := session.(string)
		//beego.Info("name:", name)

		var user models.User
		user.Name = name
		o.Read(&user, "Name")
		m2m := o.QueryM2M(&artical, "LookedUsers")
		m2m.Add(user)

		//o.LoadRelated(&artical, "Users")//查询文章所有看过的用户,无法去重
		//beego.Info("artical:", artical)

		var users []*models.User
		//查询用户中看过某一篇文章的用户，按照用户id去重
		o.QueryTable("User").Filter("LookArticals__Artical__Id", artical.Id).Distinct().All(&users)
		//beego.Info("users:", users)
		this.Data["users"] = users
	}

	this.Data["artical"] = artical
	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

func (this *MainPageController) ShowEditArtical() {
	id, err := this.GetInt("Id")
	if err != nil {
		errmsg := "获取文章详情参数错误"
		this.Redirect("/artical/mainPage?errmsg="+errmsg, 302)
		return
	}
	var artical models.Artical
	artical.Id = id
	err = orm.NewOrm().Read(&artical)
	if err != nil {
		errmsg := "获取文章详情失败"
		this.Redirect("/artical/mainPage?errmsg="+errmsg, 302)
		return
	}

	this.Data["artical"] = artical
	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

func (this *MainPageController) HandelEditArtical() {

	articleName := this.GetString("articleName")
	content := this.GetString("content")
	if articleName == "" || content == "" {
		errmsg := "文章标题和文章内容不能为空"
		id, _ := this.GetInt("Id")
		this.Redirect("/artical/updateArtical?Id="+strconv.Itoa(id)+"&errmsg="+errmsg, 302)
		return
	}

	fileUrl, err := this.uploadFile("uploadname")
	if err != nil {
		id, _ := this.GetInt("Id")
		this.Redirect("/artical/updateArtical?Id="+strconv.Itoa(id)+"&errmsg="+err.Error(), 302)
		return
	}

	var artical models.Artical
	id, err := this.GetInt("Id")
	if err != nil {
		errmsg := "参数错误"
		id, _ := this.GetInt("Id")
		this.Redirect("/artical/updateArtical?Id="+strconv.Itoa(id)+"&errmsg="+errmsg, 302)
		return
	}
	artical.Id = id

	o := orm.NewOrm()

	err = o.Read(&artical)
	if err != nil {
		errmsg := "文章不存在,修改失败"
		id, _ := this.GetInt("Id")
		this.Redirect("/artical/updateArtical?Id="+strconv.Itoa(id)+"&errmsg="+errmsg, 302)
		return
	}

	artical.Title = articleName
	artical.Content = content
	if fileUrl != "" {
		artical.ImageUrl = fileUrl
	}
	n, err := orm.NewOrm().Update(&artical)
	if err != nil || n <= 0 {
		errmsg := "文章更新失败"
		id, _ := this.GetInt("Id")
		this.Redirect("/artical/updateArtical?Id="+strconv.Itoa(id)+"&errmsg="+errmsg, 302)
		return
	}
	this.Redirect("/artical/mainPage", 302)
}
func (this *MainPageController) HandelDeleteArtical() {
	id, err := this.GetInt("Id")

	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}
	this.Layout = "layout.html"
	if err != nil {
		this.Data["errmsg"] = err.Error()
		this.TplName = "index.html"
		return
	}
	var artical models.Artical
	artical.Id = id
	_, err = orm.NewOrm().Delete(&artical)
	if err != nil {
		this.Data["errmsg"] = err.Error()
		this.TplName = "index.html"
		return
	}
	this.Redirect("/artical/mainPage", 302)

}

func (this *MainPageController) uploadFile(name string) (str string, err error) {
	fp, header, err := this.GetFile(name)
	if err != nil {
		//beego.Error(err)
		errmsg := "文件上传失败"
		return "", errors.New(errmsg)
	}
	defer fp.Close()

	var fileUrl string
	if header.Size > 0 {
		if header.Size > 1024*1024 {
			errmsg := "文件过大"
			return "", errors.New(errmsg)
		}

		fileExit := path.Ext(header.Filename)
		if fileExit != ".png" && fileExit != ".jpg" && fileExit != ".jpeg" {
			errmsg := "文件格式错误"
			return "", errors.New(errmsg)
		}

		fileName := time.Now().Format("2006.01.02 15:04:05") + fileExit
		this.SaveToFile(name, "./static/upload/"+fileName)
		fileUrl = "/static/upload/" + fileName
	}

	return fileUrl, nil
}

func (this *MainPageController) ShowAddArticalType() {

	errmsg := this.GetString("errmsg")
	if errmsg != "" {
		this.Data["errmsg"] = errmsg
	}

	name := this.GetSession("name")
	if name != nil {
		this.Data["name"] = name.(string)
	}
	this.Layout = "layout.html"

	var articalTypes []models.ArticalType
	_, err := orm.NewOrm().QueryTable("ArticalType").All(&articalTypes)
	if err != nil {
		this.Data["errmsg"] = "查询失败"
		this.TplName = "addType.html"
		return
	}

	this.Data["articalTypes"] = articalTypes
	this.TplName = "addType.html"

}

func (this *MainPageController) HandelAddArticalType() {
	typeName := this.GetString("typeName")
	if typeName == "" {
		errmsg := "类型不能为空"
		this.Redirect("/artical/addArticalType?errmsg="+errmsg, 302)
		return
	}

	var articleType models.ArticalType
	articleType.TypeName = typeName
	_, err := orm.NewOrm().Insert(&articleType)
	if err != nil {
		errmsg := "入库失败"
		this.Redirect("/artical/addArticalType?errmsg="+errmsg, 302)
		return
	}

	this.Redirect("/artical/addArticalType", 302)
}

func (this *MainPageController) HandelDeleteArticalType() {
	id, err := this.GetInt("id")
	if err != nil {
		errmsg := "参数错误"
		this.Redirect("/artical/addArticalType?errmsg="+errmsg, 302)
		return
	}
	var artType models.ArticalType
	artType.Id = id
	orm.NewOrm().Delete(&artType)
	this.Redirect("/artical/addArticalType", 302)

}

func (this *MainPageController) HandelLogout() {
	session := this.GetSession("name")
	if session != nil {
		name := session.(string)
		this.DelSession("name")
		this.Ctx.SetCookie("name", name, -1)
	}

	this.Redirect("/login", 302)

}
