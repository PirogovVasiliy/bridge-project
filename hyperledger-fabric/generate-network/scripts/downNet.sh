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
CFG_DIR="organizations/cryptogen"
find "$CFG_DIR" -maxdepth 1 -type f -name 'crypto-config-*.yaml' \
     ! -name 'crypto-config-orderer.yaml' \
     -exec rm -v {} \;

if [ -d "organizations/peerOrganizations" ]; then
    rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
fi

rm compose/compose-test-net.yaml
rm compose/docker/docker-compose-test-net.yaml

echo "🧹 Removing named volumes..."
vols=(docker_orderer.example.com)
if [ -d "organizations/peerOrganizations" ]; then
  for dir in organizations/peerOrganizations/*; do
    org=$(basename "$dir")
    vols+=(docker_peer0.${org})
  done
fi
${CONTAINER_CLI} volume rm "${vols[@]}" 2>/dev/null || true
echo "✅ Volumes removed"

echo "👏 сеть успешно остановлена 👏"
