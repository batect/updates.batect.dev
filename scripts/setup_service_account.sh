#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CREDS_DIR="$SCRIPT_DIR/../.creds"
CREDS_PATH="$CREDS_DIR/gcp_service_account_${CLOUDSDK_ACTIVE_CONFIG_NAME}.json"

function main() {
  copyCredsIntoPlace
  createConfiguration
  activateServiceAccount
}

function copyCredsIntoPlace() {
  mkdir -p "$CREDS_DIR"
  echo "$GCP_SERVICE_ACCOUNT_KEY" > "$CREDS_PATH"
}

function createConfiguration() {
  "$SCRIPT_DIR/ensure_gcloud_configuration.sh"
}

function activateServiceAccount() {
  gcloud auth activate-service-account "$GCP_SERVICE_ACCOUNT_EMAIL" --key-file "$CREDS_PATH"
}

main
