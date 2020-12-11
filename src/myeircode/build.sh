#!/bin/bash


build_local(){
go build -o ../../myeircode main.go
}

docker_build(){
docker build --tag myeircode:00-alpha .
}

main(){
  case $1 in
  "local")
      build_local
      ;;
  "docker")
    export GOOS=linux
    export GOARCH=amd64
    build_local
    docker_build

  esac
}

main "$@"
