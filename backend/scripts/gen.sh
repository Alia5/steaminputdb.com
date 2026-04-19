#!/bin/sh
set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
PROTO_DIR="$(cd "${BACKEND_DIR}/../steam_protobufs" && pwd)"

cd "${BACKEND_DIR}"

echo "Generating Go types from protobufs..."
protoc \
	--go_out=steamapi \
	--go_opt=paths=source_relative \
	--go_opt=Mcommon_base.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--go_opt=Mcommon.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--go_opt=Menums.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--go_opt=Mservice_steaminputmanager.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--go_opt=Mservice_publishedfile.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--go_opt=Mservice_storequery.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--go_opt=Mservice_player.proto=github.com/Alia5/steaminputdb.com/steamapi \
	--proto_path="${PROTO_DIR}/webui" \
	--proto_path="${PROTO_DIR}/steam" \
	"${PROTO_DIR}/steam/enums.proto" \
	"${PROTO_DIR}/webui/common_base.proto" \
	"${PROTO_DIR}/webui/common.proto" \
	"${PROTO_DIR}/webui/service_publishedfile.proto" \
	"${PROTO_DIR}/webui/service_storequery.proto" \
    "${PROTO_DIR}/webui/service_player.proto"

mkdir -p steam/client
protoc \
	--go_out=steam/client \
	--go_opt=paths=source_relative \
	--go_opt=Msteammessages_base.proto=github.com/Alia5/steaminputdb.com/steam/client \
	--go_opt=Msteammessages_clientserver_login.proto=github.com/Alia5/steaminputdb.com/steam/client \
        --go_opt=Msteammessages_clientserver_appinfo.proto=github.com/Alia5/steaminputdb.com/steam/client \
        --go_opt=Menums_clientserver.proto=github.com/Alia5/steaminputdb.com/steam/client \
        --proto_path="${PROTO_DIR}/steam" \
        --proto_path="${PROTO_DIR}" \
        "${PROTO_DIR}/steam/steammessages_base.proto" \
        "${PROTO_DIR}/steam/steammessages_clientserver_login.proto" \
        "${PROTO_DIR}/steam/steammessages_clientserver_appinfo.proto" \
        "${PROTO_DIR}/steam/enums_clientserver.proto"
