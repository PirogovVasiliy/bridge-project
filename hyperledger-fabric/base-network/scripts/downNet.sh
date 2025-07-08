#!/usr/bin/env bash

: ${CONTAINER_CLI:="docker"}
if command -v ${CONTAINER_CLI}-compose > /dev/null 2>&1; then
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
else
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI} compose"}
fi

COMPOSE_FILE_BASE=compose-test-net.yaml
COMPOSE_FILES="-f compose/${COMPOSE_FILE_BASE} -f compose/${CONTAINER_CLI}/${CONTAINER_CLI}-${COMPOSE_FILE_BASE}"

SOCK="${DOCKER_HOST:-/var/run/docker.sock}"
DOCKER_SOCK="${SOCK##unix://}"

echo "🛑 Stopping Fabric network..."
DOCKER_SOCK="${DOCKER_SOCK}" ${CONTAINER_CLI_COMPOSE} ${COMPOSE_FILES} down --volumes --remove-orphans

function clearContainers() {
  echo "🧹 Removing Fabric containers..."

  ${CONTAINER_CLI} rm -f $(${CONTAINER_CLI} ps -aq --filter label=service=hyperledger-fabric) 2>/dev/null || true
  ${CONTAINER_CLI} rm -f $(${CONTAINER_CLI} ps -aq --filter name='dev-peer*') 2>/dev/null || true
  ${CONTAINER_CLI} kill $(${CONTAINER_CLI} ps -q --filter name=ccaas) 2>/dev/null || true

  echo "✅ Containers cleared"
}

function removeUnwantedImages() {
  echo "🧼 Removing chaincode images..."

  ${CONTAINER_CLI} image rm -f $(${CONTAINER_CLI} images -aq --filter reference='dev-peer*') 2>/dev/null || true

  echo "✅ Chaincode images removed"
}

clearContainers
removeUnwantedImages

echo "🧹 Cleaning up artifacts..."
if [ -d "organizations/peerOrganizations" ]; then
  rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
fi
if [ -d "channel-artifacts" ]; then
  rm -Rf channel-artifacts
fi
rm log.txt
rm basic_1.0.tar.gz

${CONTAINER_CLI} volume rm docker_orderer.example.com docker_peer0.org1.example.com docker_peer0.org2.example.com > /dev/null 2>&1 || true

echo "👏 сеть успешно остановлена 👏"
