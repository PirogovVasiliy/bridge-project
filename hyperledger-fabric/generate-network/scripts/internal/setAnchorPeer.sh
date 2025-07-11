#!/usr/bin/env bash

TEST_NETWORK_HOME=${TEST_NETWORK_HOME:-${PWD}}
. "${TEST_NETWORK_HOME}/scripts/internal/configUpdate.sh" 

createAnchorPeerUpdate() {
  echo "Fetching channel config for channel ${CHANNEL_NAME}"
  fetchChannelConfig $ORG $CHANNEL_NAME ${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}config.json

  echo "Generating anchor peer update transaction for ${ORG_DOMAINS[IDX]} on channel ${CHANNEL_NAME}"

  IDX=$(( ORG - 1 ))
  HOST="peer0.${ORG_DOMAINS[IDX]}"
  PORT="${PEER_PORTS[IDX]}"

  set -x
  jq '.channel_group.groups.Application.groups.'${CORE_PEER_LOCALMSPID}'.values +=
      {"AnchorPeers":{"mod_policy":"Admins","value":{"anchor_peers":[{"host":"'${HOST}'","port":'${PORT}'}]},"version":"0"}}' \
      "${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}config.json" \
      > "${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}modified_config.json"
  res=$?
  { set +x; } 2>/dev/null
  verifyResult $res "Channel configuration update for anchor peer failed (jq?)"

  createConfigUpdate "${CHANNEL_NAME}" \
      "${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}config.json" \
      "${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}modified_config.json" \
      "${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx"
}

updateAnchorPeer() {
  peer channel update \
       -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com \
       -c "${CHANNEL_NAME}" \
       -f "${TEST_NETWORK_HOME}/channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx" \
       --tls --cafile "$ORDERER_CA" >&log.txt
  res=$?
  cat log.txt
  verifyResult $res "Anchor peer update failed"
  echo "Anchor peer set for org '${CORE_PEER_LOCALMSPID}' on channel '${CHANNEL_NAME}'"
}

ORG=$1
CHANNEL_NAME=$2

setGlobals $ORG
createAnchorPeerUpdate
updateAnchorPeer
