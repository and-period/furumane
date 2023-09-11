#!/bin/sh

MOCKGEN_PATH=$(which mockgen)
if [ ! -f "${MOCKGEN_PATH}" ]; then
  echo "Not found mockgen. Please install this tool."
  echo "-> $ go install go.uber.org/mock/mockgen@latest"
  exit 1
fi

#############################
# Function
#############################
build_mock() {
  path=${1##$(PWD)/proto/}

  dirname=${path%%/*}
  filename=${path##*/}

  mockgen -source ./proto/${path} -destination mock/proto/${dirname}/${filename}.go
}

#############################
# Main
#############################
paths=$(find $(PWD)/proto/**/*service.pb.go -type f)
for path in ${paths}; do
  build_mock ${path}
done
