#!/bin/bash

rm beamgotest-* # Clear out old node data

screen -XS beam-go-test-node quit
screen -XS beam-go-test-node-nomine quit
screen -XS beam-go-test-wallet-api quit
screen -XS beam-go-test-wallet-api-nomine quit