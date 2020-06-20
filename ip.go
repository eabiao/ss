package main

import (
	"github.com/ip2location/ip2location-go"
	"log"
	"net"
)

var (
	locationDB = initIPLocationDataBase()
)

// 初始化ip地理位置数据库
func initIPLocationDataBase() *ip2location.DB {
	db, err := ip2location.OpenDB("IP2LOCATION-LITE-DB1.BIN")
	if err != nil {
		log.Fatal("missing ip location database file")
	}
	return db
}

// 解析域名地理位置
func getHostLocation(host string) string {
	//判断是否为IP
	if net.ParseIP(host) != nil {
		return "-"
	}

	// 解析IP
	ips, err := net.LookupIP(host)
	if err != nil {
		return "NONE"
	}

	// 解析地理位置
	location, err := locationDB.Get_all(ips[0].String())
	if err != nil {
		return "NONE"
	}
	return location.Country_short
}
