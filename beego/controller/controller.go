package controller

import (
        "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
        "github.com/xiying/xytool/simini"
	//"github.com/astaxie/beego/context"
)

type Controller struct {
	beego.Controller
}

type CtrHandler interface {
	Handler(*Controller) bool
        AddHandler(path string, h func(*Controller)int)
        DelHandler(path string)
}

var globalSessions *session.Manager
var globalClientMap map[string]CtrHandler

func SetStaticPath(k, v string) {
    beego.SetStaticPath(k,v)
}

func Router(rootpath string, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
    fmt.Println("controler::Router")
    return beego.Router(rootpath, c, mappingMethods...)
}

func Run(ini *simini.SimIni, params ...string) {
    fmt.Println("controler::Run")
    sess := ini.GetSession("session")
    for k,v := range sess {
        fmt.Println("sess|k="+k+"|v="+v)
        s,e := session.NewManager(k, v)
        if e != nil {
            beego.Error(e.Error())
        }
	globalSessions = s
	go globalSessions.GC()
    }
    beego.Run(params...)
}

//func AddCtrHandler(s string, h CtrHandler) {
//	globalClientMap[s] = h
//}
//
//func DelCtrHandler(s string) {
//	delete(globalClientMap, s)
//}

func (this *Controller) Get() {
	h, ok := globalClientMap[this.Path()]
	if ok {
		h.Handler(this)
	}
}

func (this *Controller) Path() string {
    return this.Ctx.Request.URL.Path
}

func (this *Controller) FormValue(key string) string {
	return this.Ctx.Request.FormValue(key)
}

func (this *Controller) WriteString(s string) {
	this.Ctx.WriteString(s)
}

func (this *Controller) Write(s []byte) {
        this.Ctx.ResponseWriter.Write(s)
}

func (this *Controller) SessionStart() (session session.SessionStore, err error) {
	return globalSessions.SessionStart(this.Ctx.ResponseWriter,
		this.Ctx.Request)
}

func (this *Controller) SessionRelease(session session.SessionStore) {
	session.SessionRelease(this.Ctx.ResponseWriter)
}

func init() {
	globalClientMap = make(map[string]CtrHandler)
}
