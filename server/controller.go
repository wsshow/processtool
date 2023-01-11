package server

import (
	"net/http"
	"processtool/core"

	"github.com/gin-gonic/gin"
)

var resp Response

func GetProcs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		procs, err := core.Get().Processes()
		if err != nil {
			ctx.JSON(http.StatusOK, resp.Failure().WithDesc(err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, resp.Success(procs))
	}
}

func NoRoute(err error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, resp.Failure().WithDesc(err.Error()))
	}
}
