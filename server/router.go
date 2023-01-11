package server

import (
	"embed"
	"io/fs"
	"net/http"
	"processtool/server/middleware"

	"github.com/gin-gonic/gin"
)

//go:embed ui/dist
var tmpl embed.FS

func initApiServerRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/procs", GetProcs())
	return r
}

func initUIServerRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())
	fads, err := fs.Sub(tmpl, "ui/dist")
	if err != nil {
		r.NoRoute(NoRoute(err))
		return r
	}
	r.StaticFS("/", http.FS(fads))
	return r
}
