#!/bin/bash

set -exu
set -o pipefail

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is not installed. Please install jq first."
    exit 1
fi
# Check if curl is installed
if ! command -v curl &> /dev/null; then
    echo "Error: curl is not installed. Please install curl first."
    exit 1
fi

# NETWORK_DIR is where all files for the testnet will be stored,
# including logs and storage
NETWORK_DIR=./network

# Change this number for your desired number of nodes
NUM_NODES=4

# Port information. All ports will be incremented upon
# with more validators to prevent port conflicts on a single machine
GETH_BOOTNODE_PORT=30301

GETH_HTTP_PORT=8000
GETH_WS_PORT=8100
GETH_AUTH_RPC_PORT=8200
GETH_METRICS_PORT=8300
GETH_NETWORK_PORT=8400

PRYSM_BEACON_RPC_PORT=4000
PRYSM_BEACON_GRPC_GATEWAY_PORT=4100
PRYSM_BEACON_P2P_TCP_PORT=4200
PRYSM_BEACON_P2P_UDP_PORT=4300
PRYSM_BEACON_MONITORING_PORT=4400

PRYSM_VALIDATOR_RPC_PORT=7000
PRYSM_VALIDATOR_GRPC_GATEWAY_PORT=7100
PRYSM_VALIDATOR_MONITORING_PORT=7200


trap 'echo "Error on line $LINENO"; exit 1' ERR
# Function to handle the cleanup
cleanup() {
    echo "Caught Ctrl+C. Killing active background processes and exiting."
    kill $(jobs -p)  # Kills all background processes started in this script
    exit
}
# Trap the SIGINT signal and call the cleanup function when it's caught
trap 'cleanup' SIGINT

# Reset the data from any previous runs and kill any hanging runtimes
rm -rf "$NETWORK_DIR" || echo "no network directory"
mkdir -p $NETWORK_DIR
pkill geth || echo "No existing geth processes"
pkill beacon-chain || echo "No existing beacon-chain processes"
pkill validator || echo "No existing validator processes"
pkill bootnode || echo "No existing bootnode processes"

# Set Paths for your binaries. Configure as you wish, particularly
# if you're developing on a local fork of geth/prysm
GETH_BINARY=./dependencies/agora-el/build/bin/geth
GETH_BOOTNODE_BINARY=./dependencies/agora-el/build/bin/bootnode

PRYSM_CTL_BINARY=./dependencies/agora-cl/bazel-bin/cmd/prysmctl/prysmctl_/prysmctl
PRYSM_BEACON_BINARY=./dependencies/agora-cl/bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain
PRYSM_VALIDATOR_BINARY=./dependencies/agora-cl/bazel-bin/cmd/validator/validator_/validator

# Create the bootnode for execution client peer discovery. 
# Not a production grade bootnode. Does not do peer discovery for consensus client
mkdir -p $NETWORK_DIR/bootnode

$GETH_BOOTNODE_BINARY -genkey $NETWORK_DIR/bootnode/nodekey

$GETH_BOOTNODE_BINARY \
    -nodekey $NETWORK_DIR/bootnode/nodekey \
    -addr=:$GETH_BOOTNODE_PORT \
    -verbosity=5 > "$NETWORK_DIR/bootnode/bootnode.log" 2>&1 &

sleep 2
# Get the ENODE from the first line of the logs for the bootnode
bootnode_enode=$(head -n 1 $NETWORK_DIR/bootnode/bootnode.log)
# Check if the line begins with "enode"
if [[ "$bootnode_enode" == enode* ]]; then
    echo "bootnode enode is: $bootnode_enode"
else
    echo "The bootnode enode was not found. Exiting."
    exit 1
fi


# Generate the genesis. This will generate validators based
# on https://github.com/ethereum/eth2.0-pm/blob/a085c9870f3956d6228ed2a40cd37f0c6580ecd7/interop/mocked_start/README.md
$PRYSM_CTL_BINARY testnet generate-genesis \
--fork=capella \
--num-validators=$NUM_NODES \
--chain-config-file=./config.yml \
--geth-genesis-json-in=./genesis.json \
--output-ssz=$NETWORK_DIR/genesis.ssz \
--geth-genesis-json-out=$NETWORK_DIR/genesis.json

