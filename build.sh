#!/bin/bash
set -e

echo "loading schema for tests..."
apk add --update postgresql-client
psql -U postgres -h postgresql -d vape_test -q < ./schema.sql >/dev/null 2>&1
