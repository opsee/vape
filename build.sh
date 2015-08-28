#!/bin/bash
set -e

echo "loading schema for tests..."
echo "drop database if exists vape_test; create database vape_test" | psql -U postgres -h postgresql
migrate -url postgres://postgres@postgresql/vape_test?sslmode=disable -path ./migrations up
