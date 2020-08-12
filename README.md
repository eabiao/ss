## ss代理

配置文件 config.json
```
{
    "listen": "127.0.0.1:7777",
    "server": "ip:port",
    "method": "AEAD_CHACHA20_POLY1305",
    "passwd": "password"
}
```

### 图标添加方式

生成图标文件
```
yum install icoutils
icotool -c ss.png -o assets/ss.ico
```

图标配置文件 assets/manifest.txt
```
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
<assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="x86"
    name="controls"
    type="win32"
></assemblyIdentity>
<dependency>
    <dependentAssembly>
        <assemblyIdentity
            type="win32"
            name="Microsoft.Windows.Common-Controls"
            version="6.0.0.0"
            processorArchitecture="*"
            publicKeyToken="6595b64144ccf1df"
            language="*"
        ></assemblyIdentity>
    </dependentAssembly>
</dependency>
</assembly>
```

生成图标工程文件
```
go get github.com/akavel/rsrc
rsrc.exe -manifest assets/manifest.txt -ico assets/ss.ico -o main.syso
```

build
```
./build.sh
```