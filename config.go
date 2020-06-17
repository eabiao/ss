package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	config = initConfig()
)

// 配置
type Config struct {
	Listen string
	Server string
	Method string
	Passwd string
}

// 配置读取
func initConfig() *Config {
	data, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Fatal("missing config file")
	}

	var config = &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal("read config fail")
	}

	return config
}
