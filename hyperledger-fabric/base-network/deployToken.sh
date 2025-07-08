#!/bin/bash

# Путь к бинарникам Fabric и конфигам
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/../config/

# Название и путь до чейнкода
CHAINCODE_NAME=basic
CHAINCODE_LABEL=basic_1.0
CHAINCODE_PATH=../../gotoken/
CHAINCODE_LANG=golang
CHAINCODE_VERSION=1.0
CHANNEL_NAME=mychannel
SEQUENCE=1

# Удалим предыдущий пакет, если есть
rm -f ${CHAINCODE_LABEL}.tar.gz

# Упаковка чейнкода
peer lifecycle chaincode package ${CHAINCODE_LABEL}.tar.gz \
  --path ${CHAINCODE_PATH} --lang ${CHAINCODE_LANG} --label ${CHAINCODE_LABEL}

echo "📦 Чейнкод упакован"

# ==========================
# Установка на Org1
# ==========================
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install ${CHAINCODE_LABEL}.tar.gz
echo "✅ Установлено на Org1"

# ==========================
# Установка на Org2
# ==========================
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

peer lifecycle chaincode install ${CHAINCODE_LABEL}.tar.gz
echo "✅ Установлено на Org2"

# ==========================
# Получаем package ID
# ==========================
PACKAGE_ID=$(peer lifecycle chaincode queryinstalled | grep ${CHAINCODE_LABEL} | awk -F "[, ]+" '{print $3}')
echo "📦 Package ID: $PACKAGE_ID"

# ==========================
# Одобрение от Org2
# ==========================
peer lifecycle chaincode approveformyorg -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} \
--version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} \
--sequence ${SEQUENCE} --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

echo "✅ Одобрено Org2"

# ==========================
# Одобрение от Org1
# ==========================
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode approveformyorg -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com \
--channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} \
--version ${CHAINCODE_VERSION} --package-id ${PACKAGE_ID} \
--sequence ${SEQUENCE} --tls \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

echo "✅ Одобрено Org1"

# ==========================
# Проверка готовности
# ==========================
peer lifecycle chaincode checkcommitreadiness \
--channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} \
--version ${CHAINCODE_VERSION} --sequence ${SEQUENCE} \
--tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
--output json

# ==========================
# Коммит
# ==========================
peer lifecycle chaincode commit -o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com --channelID ${CHANNEL_NAME} \
--name ${CHAINCODE_NAME} --version ${CHAINCODE_VERSION} --sequence ${SEQUENCE} \
--tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
--peerAddresses localhost:7051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
--peerAddresses localhost:9051 \
--tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

echo "✅ Коммит выполнен"

# ==========================
# Проверка коммита
# ==========================
peer lifecycle chaincode querycommitted \
--channelID ${CHANNEL_NAME} --name ${CHAINCODE_NAME} \
--cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

echo "🎉 Чейнкод '${CHAINCODE_NAME}' успешно задеплоен на канал '${CHANNEL_NAME}'"