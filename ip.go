package main

import (
	"github.com/ip2location/ip2location-go"
	"log"
	"net"
)

var (
	ipDB = initIPDataBase()
)

// 初始化ip地理位置数据库
func initIPDataBase() *ip2location.DB {
	db, err := ip2location.OpenDB("IP2LOCATION-LITE-DB1.BIN")
	if err != nil {
		log.Fatal("missing ip location database file")
	}
	return db
}

// 解析IP地址
func getIPFromHost(host string) (string, error) {
	if net.ParseIP(host) != nil {
		return host, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}

// 解析地理位置
func getIPLocation(ip string) (ip2location.IP2Locationrecord, error) {
	return ipDB.Get_all(ip)
}
