#!/bin/bash 

set -e 

buildServer(){
  export GOOS=linux
  export GOARCH=amd64

  go build -o myeircode
}

getLatestTag(){
  i=$(docker image ls | grep gcr.io | grep redirect | awk '{print $2}' | cut -d '.' -f 2 | head -n 1)
  echo $((i+1))
}

dockerStuff(){
  latest=$(getLatestTag)
  base="eu.gcr.io/amplified-album-259521/myeircode"
  
  docker build --tag ${base}:${latest} .
  docker push ${base}:${latest}
}

main(){
  buildServer
  dockerStuff
}

main "@"
