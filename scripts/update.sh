#!/bin/bash

if [[ -z $1 ]]; then
	echo "error: please set name"
	exit 1
fi

if [[ -z $2 ]]; then
	echo "error: please set host"
	exit 1
fi

INSTALL="
rm -r /root/$1

cp -a /root/update/. /root/$1/
cp -a /root/update-data/. /root/$1-data/

rm -r /root/update
rm -r /root/update-data

systemctl restart $1"

mkdir update
mkdir update-data

CGO_ENABLED="0" GOARCH="amd64" GOOS="linux" go build -ldflags='-s -w' -trimpath -o update/main ./$1

cp $1-data/.configuration update-data/.configuration

scp -pr update root@$2:/root/update
scp -pr update-data root@$2:/root/update-data

rm -r update
rm -r update-data

ssh root@$2 "$INSTALL"
