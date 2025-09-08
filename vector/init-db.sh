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
echo -n '
CREATE TABLE logs.demo
(
    `bytes` UInt8,
    `event_time`  DateTime,
    `host` String,
    `method` String,
    `protocol` String,
    `referer` String,
    `request` String,
    `status` String,
    `user-identifier` String,
)
ENGINE = MergeTree
PARTITION BY toStartOfHour(event_time)
ORDER BY (event_time)
SETTINGS index_granularity = 8192
;' | clickhouse-client
