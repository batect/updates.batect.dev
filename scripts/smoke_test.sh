#! /usr/bin/env bash

set -euo pipefail

BASE_URL=${1:-https://updates.batect.dev}

function main() {
  checkPing

  echoGreenText "Smoke test completed successfully."
}

function checkPing() {
  echoBlueText "Checking /ping..."

  RESPONSE=$(curl \
    --fail \
    --silent \
    --verbose \
    --show-error \
    "$BASE_URL/ping"
  )

  echo "Response:"
  echo "$RESPONSE"
  echo

  diff -U 9999 <(echo "$RESPONSE") <(echo "pong") || { echo; echoRedText "Response was not as expected. See diff above. '-' represents what was expected, '+' represents what was returned by the API."; exit 1; }

  echo "/ping check passed."
  echo
}

function echoBlueText() {
  echo "$(tput setaf 4)$1$(tput sgr0)"
}

function echoGreenText() {
  echo "$(tput setaf 2)$1$(tput sgr0)"
}

function echoRedText() {
  echo "$(tput setaf 1)$1$(tput sgr0)"
}

main
