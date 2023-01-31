#!/bin/bash

if [[ -z $1 ]]; then
	echo "error: please set name"
	exit 1
fi

if [[ -z $2 ]]; then
	echo "error: please set host"
	exit 1
fi

SERVICE="
[Unit]
Description=$1
After=network.target
[Service]
Type=simple
User=root
WorkingDirectory=/root/$1-data
ExecStart=/root/$1/main
[Install]
WantedBy=multi-user.target"

INSTALL="
rm -r /root/$1
rm -r /root/$1-data

mv -f /root/update /root/$1
mv -f /root/update-data /root/$1-data

mv -f /root/$1/$1.service /etc/systemd/system/$1.service

systemctl --now enable $1"

mkdir update
mkdir update-data

CGO_ENABLED="0" GOARCH="amd64" GOOS="linux" go build -ldflags='-s -w' -trimpath -o update/main ./$1
echo "$SERVICE" > update/"$1".service

cp $1-data/.configuration update-data/.configuration

scp -pr update root@$2:/root/update
scp -pr update-data root@$2:/root/update-data

rm -r update
rm -r update-data

ssh root@$2 "$INSTALL"
