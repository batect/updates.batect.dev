#! /usr/bin/env bash

set -euo pipefail

BASE_URL="${1:-https://$DOMAIN}"
USER_AGENT="UpdatesServiceSmokeTest/1.0.0"

function main() {
  checkPing
  checkLatest
  checkDownload

  echoGreenText "Smoke test completed successfully."
}

function checkPing() {
  echoBlueText "Checking /ping..."

  RESPONSE=$(
    curl \
      --fail \
      --silent \
      --verbose \
      --show-error \
      --user-agent "$USER_AGENT" \
      "$BASE_URL/ping"
  )

  echo
  echo "Response:"
  echo "$RESPONSE"
  echo

  diff -U 9999 <(echo "$RESPONSE") <(echo "pong") || {
    echo
    echoRedText "Response was not as expected. See diff above. '-' represents what was expected, '+' represents what was returned by the API."
    exit 1
  }

  echo "/ping check passed."
  echo
}

function checkLatest() {
  echoBlueText "Checking /v1/latest..."

  RESPONSE=$(
    curl \
      --fail \
      --silent \
      --verbose \
      --show-error \
      --user-agent "$USER_AGENT" \
      "$BASE_URL/v1/latest"
  )

  echo
  echo "Response:"
  echo "$RESPONSE"
  echo

  # FIXME: this is a bit of a hack - this checks that the response is well-formed JSON and has a `url` key.
  URL=$(echo "$RESPONSE" | jq -r '.url')
  echo "$URL" | grep -q 'https://github.com/batect/batect/releases/tag' || {
    echo
    echoRedText "Response was not as expected. See response above."
    exit 1
  }

  echo "/v1/latest check passed."
  echo
}

function checkDownload() {
  echoBlueText "Checking /v1/files/:version/:filename..."

  RESPONSE=$(
    curl \
      --fail \
      --silent \
      --verbose \
      --show-error \
      --user-agent "$USER_AGENT" \
      --include \
      "$BASE_URL/v1/files/0.0.0/batect-0.0.0.jar"
  )

  echo
  echo "Response headers:"
  echo "$RESPONSE"

  echo "$RESPONSE" | grep -q 'HTTP/2 302' || {
    echo
    echoRedText "Response was not a HTTP 302 response. See response headers above."
    exit 1
  }

  echo "$RESPONSE" | grep -q 'location: https://github.com/batect/batect/releases/download/0.0.0/batect-0.0.0.jar' || {
    echo
    echoRedText "Response did not contain a Location header with the expected URL. See response headers above."
    exit 1
  }

  echo "/v1/files/:version/:filename check passed."
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
