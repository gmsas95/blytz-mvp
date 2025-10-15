#!/bin/bash
# Quick test script for individual services

SERVICE=${1:-all}

if [ "$SERVICE" == "all" ]; then
    ./scripts/local-ci.sh test-all
else
    ./scripts/local-ci.sh test SERVICE=$SERVICE
fi