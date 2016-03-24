#!/bin/bash
set -e

echo "loading schema for tests..."
echo "drop database if exists vape_test; create database vape_test" | psql -U postgres -h postgres
migrate -url $POSTGRES_CONN -path ./migrations up
