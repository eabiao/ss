## ss代理

配置文件config.json
```
{
    "listen": "127.0.0.1:7777",
    "server": "SERVER_IP:SERVER_PORT",
    "method": "aes-256-cfb",
    "passwd": "PASSWORD!"
}
```

构建
```
./build.sh
```

ico图标文件生成
```
yum install icoutils
icotool -c aa.png -o aa.ico
```