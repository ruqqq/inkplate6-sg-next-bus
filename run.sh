#!/usr/bin/with-contenv bashio
CONFIG_PATH=/data/options.json

# export REMINDER_OFFSET="$(jq --raw-output '.reminder_offset' $CONFIG_PATH)"

./inkplate6-sg-next-bus
