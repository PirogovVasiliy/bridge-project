#!/usr/bin/env bash

TEST_NETWORK_HOME=${TEST_NETWORK_HOME:-${PWD}}
. ${TEST_NETWORK_HOME}/scripts/internal/envVar.sh

fetchChannelConfig() {
  ORG=$1
  CHANNEL=$2
  OUTPUT=$3

  setGlobals $ORG

  echo "Fetching the most recent configuration block for the channel"
  set -x
  peer channel fetch config ${TEST_NETWORK_HOME}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL --tls --cafile "$ORDERER_CA"
  { set +x; } 2>/dev/null

  echo "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  configtxlator proto_decode --input ${TEST_NETWORK_HOME}/channel-artifacts/config_block.pb --type common.Block --output ${TEST_NETWORK_HOME}/channel-artifacts/config_block.json
  jq .data.data[0].payload.data.config ${TEST_NETWORK_HOME}/channel-artifacts/config_block.json >"${OUTPUT}"
  res=$?
  { set +x; } 2>/dev/null
  verifyResult $res "Failed to parse channel configuration, make sure you have jq installed"
}

createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4

  set -x
  configtxlator proto_encode --input "${ORIGINAL}" --type common.Config --output ${TEST_NETWORK_HOME}/channel-artifacts/original_config.pb
  configtxlator proto_encode --input "${MODIFIED}" --type common.Config --output ${TEST_NETWORK_HOME}/channel-artifacts/modified_config.pb
  configtxlator compute_update --channel_id "${CHANNEL}" --original ${TEST_NETWORK_HOME}/channel-artifacts/original_config.pb --updated ${TEST_NETWORK_HOME}/channel-artifacts/modified_config.pb --output ${TEST_NETWORK_HOME}/channel-artifacts/config_update.pb
  configtxlator proto_decode --input ${TEST_NETWORK_HOME}/channel-artifacts/config_update.pb --type common.ConfigUpdate --output ${TEST_NETWORK_HOME}/channel-artifacts/config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat ${TEST_NETWORK_HOME}/channel-artifacts/config_update.json)'}}}' | jq . > ${TEST_NETWORK_HOME}/channel-artifacts/config_update_in_envelope.json
  configtxlator proto_encode --input ${TEST_NETWORK_HOME}/channel-artifacts/config_update_in_envelope.json --type common.Envelope --output "${OUTPUT}"
  { set +x; } 2>/dev/null
}