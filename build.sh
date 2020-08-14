#!/bin/bash

APP_NAME=ss
ICON_FILE=plane.ico

GO=/mnt/d/soft/go/bin/go.exe
UPX=/mnt/d/soft/upx-3.96-win64/upx.exe

$GO generate

$GO get github.com/akavel/rsrc
rsrc.exe -manifest ico.manifest -ico ${ICON_FILE} -o ${APP_NAME}.syso
rm -f *.manifest

rm -f bin/*.exe
$GO build -ldflags "-s -w" -o bin/tmp.exe
$UPX -9 -o bin/${APP_NAME}.exe bin/tmp.exe
rm -f *.syso bin/tmp.exe