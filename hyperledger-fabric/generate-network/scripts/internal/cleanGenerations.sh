#!/usr/bin/env bash

function cleaningGenerations(){
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
  rm scripts/internal/envVar.env
  rm configtx/configtx.yaml
  rm -rf channel-artifacts
  rm log.txt
  rm basic_1.0.tar.gz
  rm deployToken.sh
}

cleaningGenerations