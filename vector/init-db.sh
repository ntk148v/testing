#!/bin/bash
set -e

clickhouse client -n <<-EOSQL
CREATE TABLE IF NOT EXISTS logs.logs (
    timestamp DateTime,
    host String,
    job String,
    message String
) ENGINE = MergeTree()
ORDER BY timestamp;
EOSQL
