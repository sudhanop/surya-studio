#!/bin/bash
cd "$(dirname "$0")/backend"
echo "Starting Surya Photography API on http://localhost:8080"
go run ./cmd/server
read -p "Press [Enter] key to continue..."
