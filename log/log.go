package log

import (
	"log"
	"processtool/config"
)

func init() {
	log.SetPrefix("[processtool] ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
}

func Debug(v ...interface{}) {
	if config.IsDebug() {
		Println(v...)
	}
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}
