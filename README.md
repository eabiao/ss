## ss client

config.json
```
{
    "listen": "127.0.0.1:7777",
    "server": "SERVER_IP:SERVER_PORT",
    "method": "aes-256-cfb",
    "passwd": "PASSWORD!"
}
```

build
```
./build.sh
```

png to ico
```
yum install icoutils
icotool -c aa.png -o aa.ico
```