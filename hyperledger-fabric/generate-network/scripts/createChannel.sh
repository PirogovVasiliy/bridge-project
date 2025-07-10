#!/usr/bin/env bash

. scripts/internal/envVar.sh

: ${CONTAINER_CLI:="docker"}
if command -v ${CONTAINER_CLI}-compose > /dev/null 2>&1; then
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
else
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI} compose"}
fi

if ! $CONTAINER_CLI info > /dev/null 2>&1 ; then
    echo "$CONTAINER_CLI network is required to be running to create a channel"
    exit 1
fi

: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="2"}
: ${VERBOSE:="false"}
: ${BFT:=0}

if [ ! -d "channel-artifacts" ]; then
	mkdir channel-artifacts
fi

function createChannelGenesisBlock() {
    setGlobals 1
	which configtxgen
	if [ "$?" -ne 0 ]; then
		echo "configtxgen tool not found."
        exit 1
	fi
	local bft_true=$1
	set -x

	if [ $bft_true -eq 1 ]; then
		configtxgen -profile ChannelUsingBFT -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
	else
		configtxgen -profile ChannelUsingRaft -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
	fi
	res=$?
	{ set +x; } 2>/dev/null
    verifyResult $res "Failed to generate channel configuration transaction..."
}

function createChannel() {
  local rc=1
  local COUNTER=1
  local bft_true="$1"

  echo "Adding orderers"
  while [[ $rc -ne 0 && $COUNTER -lt $MAX_RETRY ]]; do
      sleep "$DELAY"

      source scripts/internal/orderer.sh "${CHANNEL_NAME}" >> log.txt 2>&1
      res=$?

      rc=$res
      (( COUNTER++ ))
  done

  cat log.txt
  verifyResult "$rc" "Channel creation failed"
}

function joinChannel () {
  ORG=$1
  FABRIC_CFG_PATH=$PWD/../config/
  setGlobals $ORG

  local rc=1  COUNTER=1
  while [[ $rc -ne 0 && $COUNTER -lt $MAX_RETRY ]]; do
      sleep "$DELAY"

      peer channel join -b "$BLOCKFILE" >&log.txt
      rc=$?

      (( COUNTER++ ))
  done

  cat log.txt
  verifyResult "$rc" "After $MAX_RETRY attempts, peer0.${ORG_DOMAINS[ORG-1]} failed to join channel '$CHANNEL_NAME'"
}

function setAnchorPeer(){
	ORG=$1
  	. scripts/internal/setAnchorPeer.sh $ORG $CHANNEL_NAME
}

BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"
echo "Generating channel genesis block '${CHANNEL_NAME}.block'"
FABRIC_CFG_PATH=${PWD}/configtx

export FABRIC_CFG_PATH=${PWD}/configtx
createChannelGenesisBlock $BFT
echo "✅ done"

echo "Creating channel ${CHANNEL_NAME}"
createChannel $BFT
echo "Channel '$CHANNEL_NAME' created ✅"

echo "Joining peers to the channel..."
for (( i=1; i<=ORG_COUNT; i++ )); do
  joinChannel $i
done
echo "✅ done"

echo "Setting anchor peers..."
for (( i=1; i<=ORG_COUNT; i++ )); do
  setAnchorPeer $i
done
echo "💫 канал функционирует 💫"


