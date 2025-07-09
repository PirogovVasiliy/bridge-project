#!/usr/bin/env bash

export PATH=$PATH:~/go/src/hl/bridgeproject/hyperledger-fabric/bin

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

    scripts/internal/run-cryptogen.sh
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

}

COMPOSE_FILE_BASE=compose-test-net.yaml
SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
DOCKER_SOCK="${SOCK##unix://}"

networkUp

echo "👏 сеть успешно запущена 👏"