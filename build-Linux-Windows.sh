#!/usr/bin/bash
operating_systems=(linux windows)
version=0.9-20240418

for os in ${operating_systems[@]}
do
    env GOOS=${os} GOARCH=amd64 go build -o bin/workday-${version}-${os}-amd64.bin workday.go
    
    if [ "windows" == "${os}" ]; then 
        mv bin/workday-${version}-${os}-amd64.bin bin/workday-${version}-${os}-amd64.exe
    fi 

done

ls -alh bin/ 

