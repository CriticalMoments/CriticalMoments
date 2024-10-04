#!/bin/bash

set -e

JSON_CONFIG='./ios/sample_app/SampleApp/cmDevConfig.json'
SIGNED_NAME='./docs/sample_app_config.cmconfig'

echo "Signing config file..."
mkdir -p docs 
status_code=$(curl -s -X POST --data-binary @"${JSON_CONFIG}" -w "%{response_code}" --header "Content-Type: application/json" https://criticalmoments.io/account/api/sign_config -o "${SIGNED_NAME}") 
echo "Config signed with status code: $status_code"
if [ $status_code != "200" ]; then 
  echo "Error signing config: $status_code"
  cat ${SIGNED_NAME}
  rm ${SIGNED_NAME}
  exit 1
else
  echo "Config signed successfully: ${SIGNED_NAME}"
fi