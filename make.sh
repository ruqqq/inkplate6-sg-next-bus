#!/bin/sh
GOOS=linux GOARCH=arm GOARM=7 go build -o inkplate6-sg-next-bus

docker run --rm -ti --name hassio-builder --privileged \
  -v `pwd`:/data \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  homeassistant/amd64-builder -t /data --all --test \
  -i inkplate6-sg-next-bus-{arch} -d ruqqq
