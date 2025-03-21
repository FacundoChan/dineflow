#!/opt/homebrew/bin/bash

set -euo pipefail || echo "Warning: Failed to set options"

source ./scripts/lib.sh

function install_if_not_exist() {
  TOOL_NAME=$1
  INSTALL_URL=$2
  if command -v $TOOL_NAME &> /dev/null
  then
    log_callout "$TOOL_NAME is already installed."
  else
    log_cmd "$TOOL_NAME is not installed. Installing..."
    run go install "$INSTALL_URL"
  fi
}

install_if_not_exist go-cleanarch github.com/roblaszczak/go-cleanarch@latest

readonly LINT_VERSION="1.64.8"
NEED_INSTALL=false
if command -v golangci-lint >/dev/null 2>&1; then
  CURRENT_VERSION=$(golangci-lint --version | awk '{print $4}' | sed 's/^v//')
  log_callout "golangci-lint v$CURRENT_VERSION already installed."
  if [ "$CURRENT_VERSION" == "$LINT_VERSION" ]; then
    NEED_INSTALL=false
  else
    NEED_INSTALL=true
  fi
else
  NEED_INSTALL=true
fi

if [ "$NEED_INSTALL" == true ]; then
  run curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin "v$(LINT_VERSION)"
fi

run go-cleanarch

log_info "lint modules:"
log_info "$(modules)"

run goimports -w -l .

while read -r module; do
  run cd ./internal/"$module"
  run golangci-lint run  --config "$ROOT_DIR/.golangci.yaml"
  run cd -
done < <(modules)