#!/usr/bin/env bash

CFG_DIR="organizations/cryptogen"

echo "🧹  Удаляем peer-конфиги в $CFG_DIR …"

find "$CFG_DIR" -maxdepth 1 -type f -name 'crypto-config-*.yaml' \
     ! -name 'crypto-config-orderer.yaml' \
     -exec rm -v {} \;

echo "✅  Остался только crypto-config-orderer.yaml"

if [ -d "organizations/peerOrganizations" ]; then
    rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
fi

echo "👏 сеть успешно остановлена 👏"
