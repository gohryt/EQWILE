#!/bin/bash

if [[ -z $1 ]]; then
	echo "error: please set operation"
	exit 1
fi

./scripts/$1 "${@:2}"