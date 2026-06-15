#!/bin/bash
cd "$(dirname "$0")/frontend"
echo "Starting Surya Photography website on http://localhost:3000"

if [ ! -d "node_modules" ]; then
  echo "Installing frontend dependencies..."
  npm install
fi

npx next dev -p 3000
read -p "Press [Enter] key to continue..."
