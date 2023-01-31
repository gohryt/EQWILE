#!/bin/bash

if [[ -z $1 ]]; then
	echo "error: please set operation"
	exit 1
fi

echo "$1 : ${@:2}"
exit
./scripts/$1.sh "${@:2}"