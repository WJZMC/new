package main

import (
	"github.com/astaxie/beego"
	_ "news/models"
	_ "news/routers"
)

func main() {
	beego.AddFuncMap("prePage", prePage)
	beego.AddFuncMap("NextPage", nextPage)
	beego.Run()
}

func prePage(pageIndex int) int {
	pageIndex -= 1
	if pageIndex < 1 {
		pageIndex = 1
	}
	return pageIndex
}

func nextPage(pageIndex int, pageCount float64) int {
	pageIndex++
	if pageIndex > int(pageCount) {
		pageIndex = int(pageCount)
	}
	return pageIndex
}
