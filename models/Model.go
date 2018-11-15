package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id              int
	Name            string `orm:"unique"`
	Pass            string
	LookArticals    []*Artical `orm:"rel(m2m)"`
	WrittenAritcals []*Artical `orm:"reverse(many)"`
}

type Artical struct {
	Id          int          `orm:"kp;auto"`
	Title       string       `orm:"size(100)" form:"articleName"`
	Content     string       `orm:"size(500)" form:"content"`
	AddTime     time.Time    `orm:"type(datetime);auto_now_add"`
	ReadCount   int          `orm:"size(11);default(0)"`
	ImageUrl    string       `orm:"null"`
	ArticalType *ArticalType `orm:"rel(fk);on_delete(set_default);default(1)"`
	LookedUsers []*User      `orm:"reverse(many)"`
	Author      *User        `orm:"rel(fk);on_delete(set_default);default(1)"`
}

type ArticalType struct {
	Id       int
	TypeName string     `orm:"size(100)"`
	Aritcals []*Artical `orm:"reverse(many)"`
}

func init() {
	orm.RegisterDataBase("default", "mysql", "jack:123456@tcp(127.0.0.1:3306)/news?charset=utf8&loc=Local")
	orm.RegisterModel(new(User), new(Artical), new(ArticalType))
	orm.RunSyncdb("default", false, true)
}
