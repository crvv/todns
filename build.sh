#!/usr/bin/env bash

PROJECT="$(basename `pwd`)"

checkErr () {
    if [ "$1" != "0" ]; then
        echo "Build failed"
        exit $1
    fi
}

# Prepare files
cd data
if [ "$1" = "debug" -o "$2" = "debug" ]; then
    ${GOPATH}/bin/go-bindata -debug -pkg="data" -o data.go -ignore=\.\*\\.go  ./...
else
    ${GOPATH}/bin/go-bindata -pkg="data" -o data.go -ignore=\.\*\\.go ./...
fi
cd ..

if [ "$1" = "amd64" -o "$1" = "386" ]; then
    echo $1
    export GOARCH=$1
fi

GOOS=`go env GOOS`

build_cmd='go build'
if [ "${GOOS}" = "windows" ]; then
    PROJECT=${PROJECT}.exe
    build_cmd='go build -ldflags -H=windowsgui'
fi
${build_cmd} -o ${PROJECT}
checkErr $?

echo "Build succeeded"
exit 0
