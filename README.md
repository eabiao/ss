# ss
根据IP地理位置进行分流的ss代理

config.json
```
{
    "listen": "127.0.0.1:7777",
    "server": "ip:port",
    "method": "AEAD_CHACHA20_POLY1305",
    "passwd": "password"
}
```

地理位置数据文件 IP2LOCATION-LITE-DB1.BIN  
下载地址 https://lite.ip2location.com/database/ip-country
