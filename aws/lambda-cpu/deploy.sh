#!/bin/bash
set -e
cd "$(dirname "$0")"
SCRIPT_DIR=$(pwd)

LAMBDA_SRC_DIR="${SCRIPT_DIR}/lambdasrc"
DIST_DIR="${LAMBDA_SRC_DIR}/dist"
mkdir -p "${DIST_DIR}"

cd "${LAMBDA_SRC_DIR}"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -trimpath -o "${DIST_DIR}/bootstrap-x86_64" main.go
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -trimpath -o "${DIST_DIR}/bootstrap-arm64" main.go

cd "${DIST_DIR}"
zip -j function-x86_64.zip bootstrap-x86_64
mv bootstrap-x86_64 bootstrap
zip -j function-x86_64.zip bootstrap
rm bootstrap

mv bootstrap-arm64 bootstrap
zip -j function-arm64.zip bootstrap
rm bootstrap

cd "${SCRIPT_DIR}/tf"
terraform init
terraform apply
