#!/usr/bin/bash
operating_systems=(linux windows)
version=0.1-20240325

for os in ${operating_systems[@]}
do
    env GOOS=${os} GOARCH=amd64 go build -o workday-${version}-${os}-amd64.bin workday.go
    
    if [ "windows" == "${os}" ]; then 
        mv workday-${version}-${os}-amd64.bin workday-${version}-${os}-amd64.exe
    fi 

done

ls -alh 
