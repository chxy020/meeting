package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["seater/controllers:MeetingController"] = append(beego.GlobalControllerRouter["seater/controllers:MeetingController"],
        beego.ControllerComments{
            Method: "MeetingShow",
            Router: `/sort`,
            AllowHTTPMethods: []string{"Post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
