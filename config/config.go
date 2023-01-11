package config

import (
	"flag"
	"fmt"
	"sync"
)

type Config struct {
	UIServerPort  int
	ApiServerPort int
	ShowUI        bool
	Debug         bool
}

var (
	conf *Config
	once sync.Once
	err  error
)

func Get() (*Config, error) {
	once.Do(func() {
		conf = new(Config)
		flag.IntVar(&conf.UIServerPort, "up", 9090, "ui server port")
		flag.IntVar(&conf.ApiServerPort, "ap", 9091, "api server port")
		flag.BoolVar(&conf.Debug, "debug", false, "debug mode, show more info")
		flag.BoolVar(&conf.ShowUI, "ui", true, "show web ui")
		flag.Parse()
		err = conf.Check()
	})
	return conf, err
}

func IsDebug() bool {
	return conf.Debug
}

func (c *Config) Check() error {
	if c.UIServerPort == c.ApiServerPort {
		return fmt.Errorf("port set error, ui port cannot be the same as the server port")
	}
	if c.ShowUI && (c.UIServerPort < 1 || c.UIServerPort > 65535) {
		return fmt.Errorf("port set error, range is [1, 65535], current port is %d", c.UIServerPort)
	}
	if c.ApiServerPort < 1 || c.ApiServerPort > 65535 {
		return fmt.Errorf("port set error, range is [1, 65535], current port is %d", c.ApiServerPort)
	}
	return nil
}
