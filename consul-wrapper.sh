#!/bin/bash

APP_PATH=$1
APP_CFG=${APP_PATH}-cfg.json
DEFAULT_CFG=$2
if [[ -z ${APP_PATH} || -z ${DEFAULT_CFG} ]]; then
	echo "Usage: ${0} [APPLICATION EXECUTABLE] [DEFAULT CONFIG FILE]"
	exit 1
fi

if [[ -z ${CONSUL_HTTP_ADDR} ]]; then
	echo "Please, specify consul address in CONSUL_HTTP_ADDR environment variable"
	exit 1
fi
APP_KEY=$(basename ${APP_PATH})
echo /usr/bin/curl -X GET "${CONSUL_HTTP_ADDR}/v1/kv/otn/${APP_KEY}?raw"
/usr/bin/curl -X GET "${CONSUL_HTTP_ADDR}/v1/kv/otn/${APP_KEY}?raw" > ${APP_CFG}
if [[ $? -ne 0 || ! -s ${APP_CFG} ]]; then
	echo "Can't get config from consul, using default"
	APP_CFG=${DEFAULT_CFG}
fi

./${APP_PATH} -cfg ${APP_CFG}
