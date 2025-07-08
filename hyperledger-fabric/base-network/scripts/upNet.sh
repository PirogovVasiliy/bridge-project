#!/usr/bin/env bash

: ${CONTAINER_CLI:="docker"}
if command -v ${CONTAINER_CLI}-compose > /dev/null 2>&1; then
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
else
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI} compose"}
fi

CRYPTO="cryptogen"

function checkPrereqs() {

  peer version > /dev/null 2>&1

  if [[ $? -ne 0 || ! -d "../config" ]]; then
    echo "Peer binary and configuration files not found.."
    exit 1
  fi

  LOCAL_VERSION=$(peer version | sed -ne 's/^ Version: //p')
  DOCKER_IMAGE_VERSION=$(${CONTAINER_CLI} run --rm hyperledger/fabric-peer:latest peer version | sed -ne 's/^ Version: //p')

  echo "LOCAL_VERSION=$LOCAL_VERSION"
  echo "DOCKER_IMAGE_VERSION=$DOCKER_IMAGE_VERSION"

  if [ "$LOCAL_VERSION" != "$DOCKER_IMAGE_VERSION" ]; then
    echo "Local fabric binaries and docker images are out of sync. This may cause problems."
    exit 1
  fi

  echo "✅ checkPrereqs"
}

function createOrgs(){
    if [ -d "organizations/peerOrganizations" ]; then
        rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
    fi

    which cryptogen
    if [ "$?" -ne 0 ]; then
      echo "cryptogen tool not found. exiting"
      exit 1
    fi

    echo "Creating Org1 Identities"
    set -x
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org1.yaml --output="organizations"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
        echo "Failed to generate certificates..."
        exit 1
    fi

    echo "Creating Org1 Identities"
    set -x
    cryptogen generate --config=./organizations/cryptogen/crypto-config-org2.yaml --output="organizations"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
      echo "Failed to generate certificates..."
      exit 1
    fi

    echo "Creating Orderer Org Identities"
    set -x
    cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output="organizations"
    res=$?
    { set +x; } 2>/dev/null
    if [ $res -ne 0 ]; then
      echo "Failed to generate certificates..."
      exit 1
    fi

    echo "Generating CCP files for Org1 and Org2"
    ./organizations/ccp-generate.sh
    echo "✅ createOrgs"
}

function networkUp(){
    checkPrereqs

    createOrgs

    COMPOSE_FILES="-f compose/${COMPOSE_FILE_BASE} -f compose/${CONTAINER_CLI}/${CONTAINER_CLI}-${COMPOSE_FILE_BASE}"

    DOCKER_SOCK="${DOCKER_SOCK}" ${CONTAINER_CLI_COMPOSE} ${COMPOSE_FILES} up -d 2>&1

    $CONTAINER_CLI ps -a
    if [ $? -ne 0 ]; then
        echo "Unable to start network"
        exit 1
    fi

    echo "✅ networkUp"
}

COMPOSE_FILE_BASE=compose-test-net.yaml

# Get docker sock path from environment variable
SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
DOCKER_SOCK="${SOCK##unix://}"

networkUp
echo "👏"