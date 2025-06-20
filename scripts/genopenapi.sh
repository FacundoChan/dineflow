#!/opt/homebrew/bin/bash

set -euo pipefail

# enables globstar, using `**`.
shopt -s globstar

if ! [[ "$0" =~ scripts/genopenapi.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

source ./scripts/lib.sh

OPENAPI_ROOT="./api/openapi"

GEN_SERVER=(
  # "chi-server"
  # "echo-server"
  "gin-server"
)

if [ "${#GEN_SERVER[@]}" -ne 1 ]; then
  log_error "GEN_SERVER enables more than 1 server, please check."
  exit 255
fi

log_callout "USING ${GEN_SERVER[0]}"

function openapi_files() {
  openapi_files=$(ls ${OPENAPI_ROOT})
  echo "${openapi_files[@]}"
}

# output_dir, package_name, service_name
function gen() {
  local output_dir=$1
  local package=$2
  local service=$3

  run mkdir -p "$output_dir"
  run find "$output_dir" -type f -name "*.gen.go" -delete

  prepare_dir "internal/common/client/$service"

  run go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
    -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" -config api/openapi/cfg.yaml "api/openapi/$service.yml"
  run go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
    -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" -config api/openapi/cfg.yaml "api/openapi/$service.yml"

  run go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
    -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" -config api/openapi/cfg.yaml "api/openapi/$service.yml"
  run go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen \
    -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" -config api/openapi/cfg.yaml "api/openapi/$service.yml"
}

gen internal/order/ports ports order

log_success "openapi generate successfully!"
