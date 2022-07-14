#!/bin/sh
docker push ruqqq/inkplate6-sg-next-bus-armv7:latest
docker push ruqqq/inkplate6-sg-next-bus-armv7:${jq --raw-output '.version' config.json}
