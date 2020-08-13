#!/bin/bash

APP_NAME=ss

GO=/mnt/d/soft/go/bin/go.exe
UPX=/mnt/d/soft/upx-3.96-win64/upx.exe

$GO generate

$GO get github.com/akavel/rsrc
rsrc.exe -manifest manifest.txt -ico plane.ico -o main.syso
rm -f manifest.txt

rm -f bin/*.exe
$GO build -ldflags "-s -w" -o bin/tmp.exe
$UPX -9 -o bin/${APP_NAME}.exe bin/tmp.exe
rm -f bin/tmp.exe
rm -f main.syso