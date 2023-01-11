package main

import (
	"processtool/config"
	"processtool/log"
	"processtool/server"
)

func main() {
	c, err := config.Get()
	if err != nil {
		log.Fatalln(err)
	}
	server.ServerRun(c)
}
