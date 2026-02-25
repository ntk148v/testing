#!/bin/bash
clickhouse client -n <<-EOSQL
CREATE DATABASE IF NOT EXISTS default;
EOSQL
