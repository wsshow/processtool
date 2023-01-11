package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"processtool/config"
	"processtool/log"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func initServer(port int, handler http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	return srv
}

func ServerRun(conf *config.Config) {
	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	var uiSrv *http.Server
	if conf.ShowUI {
		uiSrv = initServer(conf.UIServerPort, initUIServerRouter())
		log.Println("ui server run", uiSrv.Addr)
	}
	apiSrv := initServer(conf.ApiServerPort, initApiServerRouter())
	log.Println("api server run", apiSrv.Addr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c
	close(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if conf.ShowUI {
		if err := uiSrv.Shutdown(ctx); err != nil {
			log.Fatalln(err)
		}
	}
	if err := apiSrv.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
