#!/bin/bash
# Go Environment Setup Script for Blytz MVP
# This script sets up the proper Go environment for development

export PATH="/home/sas/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.24.9.linux-amd64/bin:$PATH"
export GOROOT="/home/sas/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.24.9.linux-amd64"
export GOPATH="/home/sas/go"

echo "✅ Go environment configured:"
echo "   PATH: $PATH"
echo "   GOROOT: $GOROOT"
echo "   GOPATH: $GOPATH"
echo
echo "Go version: $(go version)"
echo
echo "🚀 Environment ready for microservices development!"