#!/bin/bash

mkdir -p ./builds

program="j2es"
tag="$1"

# build Linux 64bit program
env GOOS=linux GOARCH=amd64 go build -o $program main/main.go
linux64=$(printf "%s_%s_linux_amd64.tar.gz" "$program" "$tag")
tar -cvzf ./builds/$linux64 $program
