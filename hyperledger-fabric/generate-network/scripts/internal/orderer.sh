#!/usr/bin/env bash

channel_name="$1"
ROOT_DIR=${ROOT_DIR:-${PWD}}

export ORDERER_ADMIN_TLS_SIGN_CERT="${ROOT_DIR}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt"
export ORDERER_ADMIN_TLS_PRIVATE_KEY="${ROOT_DIR}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key"

: "${ORDERER_CA:?ORDERER_CA is not set}"

osnadmin channel join \
  --channelID     "${channel_name}" \
  --config-block  "./channel-artifacts/${channel_name}.block" \
  -o              localhost:7053 \
  --ca-file       "${ORDERER_CA}" \
  --client-cert   "${ORDERER_ADMIN_TLS_SIGN_CERT}" \
  --client-key    "${ORDERER_ADMIN_TLS_PRIVATE_KEY}" \
  >> log.txt 2>&1

return 0
