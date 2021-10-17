#!/usr/bin/env sh
REPO_ROOT="${REPO_ROOT:-$(cd "$(dirname "$0")/.." && pwd)}"
PB_PATH="${REPO_ROOT}/api/v1/pb"
PROTO_FILE=${1:-"chat/chat.proto"}

echo "Generating pb files for ${PROTO_FILE} service"
protoc --go_out="${PB_PATH}" --go_opt=paths=source_relative \
    --go-grpc_out="${PB_PATH}" --go-grpc_opt=paths=source_relative \
    -I="${PB_PATH}" \
    "${PB_PATH}/${PROTO_FILE}"