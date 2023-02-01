#!/bin/bash

if [[ -z $1 ]]; then
	echo "error: please set operation"
	exit 1
fi

./sh/${@:1}