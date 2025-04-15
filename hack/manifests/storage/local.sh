#!/bin/bash

localPath=${1:-/opt/quickon/storage/local}
sourcePath=$(dirname "$0")/local.yaml.tpl
targetPath=$(dirname "$0")/local.deploy.yaml

mkdir -p "${localPath}"

cp -a "${sourcePath}" "${targetPath}"

sed -i "s|__LOCAL_PATH__|${localPath}|g" "${targetPath}"
