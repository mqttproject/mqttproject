#!/bin/bash

sleep 5

cd "$(dirname "$0")" || exit 1

echo "Compiling the Go application..."
if ! go build; then
    echo "Compilation failed."
    exit 1
fi

echo "Compilation successful."

echo "Running the application..."
./laite