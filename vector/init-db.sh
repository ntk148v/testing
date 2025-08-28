#!/bin/bash
set +e

clickhouse client -n <<-EOSQL
CREATE DATABASE IF NOT EXISTS logs;
CREATE TABLE IF NOT EXISTS logs.syslog (
    timestamp DateTime,
    host String,
    job String,
    message String
)
ENGINE = MergeTree()
ORDER BY tuple();
EOSQL
