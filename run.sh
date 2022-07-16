#!/usr/bin/with-contenv bashio
CONFIG_PATH=/data/options.json

export DATAMALL_ACCOUNT_KEY="$(jq --raw-output '.datamall_account_key' $CONFIG_PATH)"

./inkplate6-sg-next-bus
