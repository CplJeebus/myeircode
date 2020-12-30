#!/bin/bash 

set -e 

buildServer(){
  export GOOS=linux
  export GOARCH=amd64

  go build -o myeircode
}

getLatestTag(){
  i=$(docker image ls | grep gcr.io | grep myeircode | awk '{print $2}' | cut -d '.' -f 2 | head -n 1)
  echo $((i+1))
}

dockerStuff(){
  latest=$1
  base="${GCR_BASE}"
  
  docker build --tag ${base}:${latest} .
  docker push ${base}:${latest}
}

deploy(){
gcloud run services update myeircode --platform managed --image ${GCR_BASE}:${1} --region europe-west1
}

main(){
  buildServer
  tag=$(getLatestTag)
  dockerStuff $tag
  deploy $tag
}

main "@"
