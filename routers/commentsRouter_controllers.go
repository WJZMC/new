package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["news/controllers:ArticalController"] = append(beego.GlobalControllerRouter["news/controllers:ArticalController"],
        beego.ControllerComments{
            Method: "AddArtical",
            Router: `/AddArticalRouter`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["news/controllers:ArticalController"] = append(beego.GlobalControllerRouter["news/controllers:ArticalController"],
        beego.ControllerComments{
            Method: "ShowArticalList",
            Router: `/ShowArticalList`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
