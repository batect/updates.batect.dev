#! /usr/bin/env bash

set -euo pipefail

existing=$(gcloud config configurations list --format "value(name)" --filter "name=$CLOUDSDK_ACTIVE_CONFIG_NAME")

if [[ "$existing" != "$CLOUDSDK_ACTIVE_CONFIG_NAME" ]]; then
  gcloud config configurations create "$CLOUDSDK_ACTIVE_CONFIG_NAME" --no-activate --quiet
fi
