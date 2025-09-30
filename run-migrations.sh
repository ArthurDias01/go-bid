#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Change to migrations directory and run tern
cd internal/store/pgstore/migrations
tern migrate --config tern.conf
