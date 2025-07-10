#!/usr/bin/env bash
set -e

TEST_NETWORK_HOME=${TEST_NETWORK_HOME:-${PWD}}
export CORE_PEER_TLS_ENABLED=true
export ORDERER_CA=${TEST_NETWORK_HOME}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem

source "${TEST_NETWORK_HOME}/scripts/internal/envVar.env"


setGlobals () {
  local idx="${OVERRIDE_ORG:-$1}"
  idx=$((idx - 1))

  if [[ -z "${ORG_DOMAINS[idx]-}" ]]; then
    echo "Unknown organization index: $1"
    return 1
  fi

  export CORE_PEER_LOCALMSPID="${MSP_IDS[idx]}"
  export CORE_PEER_TLS_ROOTCERT_FILE="${TEST_NETWORK_HOME}/${TLS_CA_PATHS[idx]}"
  export CORE_PEER_MSPCONFIGPATH="${TEST_NETWORK_HOME}/organizations/peerOrganizations/${ORG_DOMAINS[idx]}/users/Admin@${ORG_DOMAINS[idx]}/msp"
  export CORE_PEER_ADDRESS="localhost:${PEER_PORTS[idx]}"

  echo "👉 Using ${ORG_DOMAINS[idx]} (MSP=${CORE_PEER_LOCALMSPID})"
  env | grep CORE
}

parsePeerConnectionParameters () {
  PEER_CONN_PARMS=()
  PEERS=""

  while [[ $# -gt 0 ]]; do
    setGlobals "$1"
    local i=$(( $1 - 1 ))
    local PEER="peer0.${ORG_DOMAINS[i]}"

    [[ -z "$PEERS" ]] && PEERS="$PEER" || PEERS="$PEERS $PEER"
    PEER_CONN_PARMS+=("--peerAddresses" "$CORE_PEER_ADDRESS")
    PEER_CONN_PARMS+=("--tlsRootCertFiles" "${TEST_NETWORK_HOME}/${TLS_CA_PATHS[i]}")

    shift
  done
}

verifyResult () {
  if [[ "$1" -ne 0 ]]; then
    echo "$2"
    exit 1
  fi
}


