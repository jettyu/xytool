package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	//"github.com/astaxie/beego/context"
)

type Controller struct {
	beego.Controller
}

type CtrHandler interface {
	Handler(*Controller) bool
        AddHandler(path string, h func(*Controller))
        DelHandler(path string)
}

var globalSessions *session.Manager
var globalClientMap map[string]CtrHandler

func SetStaticPath(k, v string) {
    beego.SetStaticPath(k,v)
}

func Router(rootpath string, h CtrHandler, c beego.ControllerInterface, mappingMethods ...string) *beego.App {
    if len(rootpath) != 0 {
        globalClientMap[rootpath] = h
    }
    return beego.Router(rootpath, c, mappingMethods...)
}

func Run(params ...string) {
    beego.Run(params...)
}

//func AddCtrHandler(s string, h CtrHandler) {
//	globalClientMap[s] = h
//}
//
//func DelCtrHandler(s string) {
//	delete(globalClientMap, s)
//}

func Init(sess *session.Manager) {
	globalSessions = sess
	go globalSessions.GC()
}

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
