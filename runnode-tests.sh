#!/bin/bash
# This bash file sets up a beam node & wallet API endpoint(s) with FakePoW, prepared for testing.
# Run this in an EMPTY directory containing suitable beam-node, wallet-api, and beam-wallet executables with a valid treasury.bin.
# Treasury.bin can be obtained from the master branch of beam on Github.
# Renaming masternet or testnet executables is OK for testing new features.

BEAM_CONSENSUS_ARGS=("--FakePoW" "1" "--Maturity.Coinbase" "1" "--Maturity.Std" "1")
BEAM_NODE_ARGS=("--treasury_path" "treasury.bin" "--pow_solve_time" "100" "--pass" "a")
WALLET_ARGS=("--pass" "a")
WALLET_ARGS_MINER=("--wallet_path" "beamgotest-wallet.db")
WALLET_ARGS_SECONDARY=("--wallet_path" "beamgotest-wallet-nomine.db")
WALLET_API_ARGS=("--use_http" "1" "${WALLET_ARGS[@]}")

trap '' 2 # Why does beam send signal 2 EVERY TIME IT EXITS?

# Initialize the wallet.
./stopnode-tests.sh

# shellcheck disable=SC2068
./beam-wallet init "${WALLET_ARGS[@]}" "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_ARGS_MINER[@]}" || true

# shellcheck disable=SC2068
./beam-wallet init "${WALLET_ARGS[@]}" "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_ARGS_SECONDARY[@]}" || true

# Get the owner key and miner key.
OUTPUT=$(./beam-wallet export_owner_key "${WALLET_ARGS[@]}" "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_ARGS_MINER[@]}")
BEAM_OWNER_KEY=$(echo $OUTPUT | sed 's/Owner Viewer key: *//')
OUTPUT=$(./beam-wallet export_owner_key "${WALLET_ARGS[@]}" "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_ARGS_SECONDARY[@]}")
BEAM_OWNER_KEY_NOMINE=$(echo $OUTPUT | sed 's/Owner Viewer key: *//') # Yes I know this is bad bash, no I don't care.

OUTPUT=$(./beam-wallet export_miner_key "${WALLET_ARGS[@]}" "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_ARGS_MINER[@]}" --subkey 1)
BEAM_MINER_KEY=$(echo $OUTPUT | sed 's/Secret Subkey 1: *//') # For some reason, sed does NOT like picking up Beam wallet's output directly. Grabbing the output as $OUTPUT and piping into sed works though.

# Copy the args for the no mining node
BEAM_NODE_ARGS_NOMINE=("${BEAM_NODE_ARGS[@]}")
# Prepare our final node args
BEAM_NODE_ARGS+=("--miner_key" $BEAM_MINER_KEY)
# Add owner keys
BEAM_NODE_ARGS+=("--owner_key" $BEAM_OWNER_KEY)
BEAM_NODE_ARGS_NOMINE+=("--owner_key" $BEAM_OWNER_KEY_NOMINE)

# Fire up the node
screen -S beam-go-test-node -d -m ./beam-node "${BEAM_CONSENSUS_ARGS[@]}" "${BEAM_NODE_ARGS[@]}" -p 10500 --storage beamgotest-node.db --mining_threads 1
screen -S beam-go-test-node-nomine -d -m ./beam-node "${BEAM_CONSENSUS_ARGS[@]}" "${BEAM_NODE_ARGS_NOMINE[@]}" -p 10501 --peer 127.0.0.1:10500 --storage beamgotest-nomine-node.db
screen -S beam-go-test-wallet-api -d -m ./wallet-api "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_API_ARGS[@]}" "${WALLET_ARGS_MINER[@]}" -p 5000 -n 127.0.0.1:10500
screen -S beam-go-test-wallet-api-nomine -d -m ./wallet-api "${BEAM_CONSENSUS_ARGS[@]}" "${WALLET_API_ARGS[@]}" "${WALLET_ARGS_SECONDARY[@]}" -p 6000 -n 127.0.0.1:10501

trap 2

echo "Node and APIs have been spun up (node at port 10500, wallets at 5000 and 6000). Locate these sessions in screen and kill them before re-running."
