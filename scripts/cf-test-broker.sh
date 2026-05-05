#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -e

SERVICE=csb-databricks-workspace
PLAN=default
NAME=databricks-test

cf create-service "${SERVICE}" "${PLAN}" "${NAME}"

cf service "${NAME}" | grep "create in progress"

set +e
while [ $? -eq 0 ]; do
    sleep 15
    cf service "${NAME}" | grep "create in progress"
done
set -e

APP=test-app

cf bind-service "${APP}" "${NAME}"

cf restart "${APP}"

cf unbind-service "${APP}" "${NAME}"

cf delete-service -f "${NAME}"
