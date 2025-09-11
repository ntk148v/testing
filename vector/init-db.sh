#!/bin/bash
set +e

clickhouse client -n <<-EOSQL
-- Create database
CREATE DATABASE IF NOT EXISTS logs;

-- Optimized syslog table
CREATE TABLE IF NOT EXISTS logs.syslog (
    timestamp DateTime CODEC(DoubleDelta, ZSTD(3)),
    host LowCardinality(String) CODEC(LZ4),
    job LowCardinality(String) CODEC(LZ4),
    message String CODEC(ZSTD(5))
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, host)
SETTINGS
    index_granularity = 8192,
    compress_marks = 1,
    compress_primary_key = 1,
    min_bytes_for_wide_part = 0;

-- Optimized demo table
CREATE TABLE IF NOT EXISTS logs.demo (
    event_time DateTime CODEC(DoubleDelta, ZSTD(3)),
    host LowCardinality(String) CODEC(LZ4),
    method LowCardinality(String) CODEC(LZ4),
    protocol LowCardinality(String) CODEC(LZ4),
    status LowCardinality(String) CODEC(LZ4),
    bytes UInt32 CODEC(Delta, LZ4),
    request String CODEC(ZSTD(3)),
    referer String CODEC(ZSTD(3)),
    \`user-identifier\` String CODEC(LZ4)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(event_time)
ORDER BY (event_time, host, status)
SETTINGS
    index_granularity = 8192,
    compress_marks = 1,
    compress_primary_key = 1,
    min_bytes_for_wide_part = 0;
EOSQL

echo "Database and tables created successfully with optimized compression settings"
