# ClickHouse Full-Text Search Guide

> Hands-on benchmark of ClickHouse Full-Text Search (FTS) on a single-node Docker setup.

**tl;dr**: FTS is **100–1000× faster** than full scan for selective terms. At 28M Hacker News rows, queries drop from seconds to single-digit milliseconds.

## Prerequisites

- Docker and Docker Compose
- ~15GB free disk
- Network access to ClickHouse's public dataset bucket

## 1. Start ClickHouse

```bash
docker compose up -d
```

Verify the server:

```bash
docker exec clickhouse-server clickhouse-client -q "SELECT version()"
```

Requires ClickHouse 25.12+. Tested with `26.6.1.1193`.

## 2. Import the Hacker News dataset

Use ClickHouse's public dataset from S3:

```bash
docker exec clickhouse-server clickhouse-client --multiquery "
DROP TABLE IF EXISTS hackernews;

CREATE TABLE hackernews
(
    id UInt64,
    deleted UInt8,
    type String,
    author String,
    time DateTime,
    text String,
    dead UInt8,
    parent UInt64,
    poll UInt64,
    kids Array(UInt32),
    url String,
    score UInt32,
    title String,
    parts Array(UInt32),
    descendants UInt32
)
ENGINE = MergeTree
ORDER BY (type, time, id);

INSERT INTO hackernews
SELECT * FROM s3(
    'https://datasets-documentation.s3.eu-west-3.amazonaws.com/hackernews/hacknernews.parquet',
    'Parquet',
    'id UInt64,
     deleted UInt8,
     type String,
     by String,
     time DateTime,
     text String,
     dead UInt8,
     parent UInt64,
     poll UInt64,
     kids Array(UInt32),
     url String,
     score UInt32,
     title String,
     parts Array(UInt32),
     descendants UInt32'
);
"
```

Check size:

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT
    count() AS rows,
    formatReadableSize(total_bytes) AS table_size
FROM system.tables
WHERE database = 'default' AND name = 'hackernews'
GROUP BY total_bytes;
"
```

Expected: **28.74M rows, ~6.61 GiB**.

## 3. Baseline: without FTS

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT count() FROM hackernews
WHERE hasToken(lowerUTF8(text), 'clickhouse')
SETTINGS enable_full_text_index = 0;
"

docker exec clickhouse-server clickhouse-client -q "
SELECT count() FROM hackernews
WHERE hasToken(lowerUTF8(text), 'olap')
  AND hasToken(lowerUTF8(text), 'oltp')
SETTINGS enable_full_text_index = 0;
"

docker exec clickhouse-server clickhouse-client -q "
SELECT toYYYYMM(time) AS month, count() AS mentions
FROM hackernews
WHERE hasToken(lowerUTF8(text), 'clickhouse')
GROUP BY month ORDER BY month
SETTINGS enable_full_text_index = 0;
"
```

## 4. Add Full-Text Search index

Add the index to the **same table**, then materialize it:

```bash
docker exec clickhouse-server clickhouse-client --multiquery "
ALTER TABLE hackernews ADD INDEX fts_text text TYPE text(
    tokenizer = splitByNonAlpha,
    preprocessor = lowerUTF8(text)
);

ALTER TABLE hackernews MATERIALIZE INDEX fts_text;
"
```

Verify materialization is done:

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT command, is_done FROM system.mutations
WHERE database='default' AND table='hackernews'
ORDER BY create_time DESC;
"
```

Check index size:

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT name, type,
    formatReadableSize(data_compressed_bytes) AS compressed,
    formatReadableSize(data_uncompressed_bytes) AS uncompressed
FROM system.data_skipping_indices
WHERE database='default' AND table='hackernews';
"
```

## 5. With FTS

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT count() FROM hackernews
WHERE hasToken(text, 'clickhouse')
SETTINGS enable_full_text_index = 1;
"

docker exec clickhouse-server clickhouse-client -q "
SELECT count() FROM hackernews
WHERE hasToken(text, 'olap')
  AND hasToken(text, 'oltp')
SETTINGS enable_full_text_index = 1;
"

docker exec clickhouse-server clickhouse-client -q "
SELECT toYYYYMM(time) AS month, count() AS mentions
FROM hackernews
WHERE hasToken(text, 'clickhouse')
GROUP BY month ORDER BY month
SETTINGS enable_full_text_index = 1;
"
```

## 6. Results

Index stats:

| Metric                              | Value                   |
| ----------------------------------- | ----------------------- |
| Rows                                | 28,737,557              |
| Table size                          | 6.61 GiB                |
| FTS index size                      | 1.81 GiB (27% of table) |
| Granules without FTS                | 3,529                   |
| Granules with FTS (`clickhouse`)    | 441 (87% skipped)       |
| Granules with FTS (`olap` + `oltp`) | 416 (88% skipped)       |

All times from `system.query_log` (best of 3 hot runs).

### Q1: Single token `clickhouse` (appears in 1,145 rows)

| Metric            | Without FTS | With FTS | Improvement      |
| ----------------- | ----------- | -------- | ---------------- |
| query_duration_ms | 537 ms      | **2 ms** | **~270× faster** |
| read_bytes        | 1.30 GiB    | 3.45 MiB | **~390× less**   |
| marks scanned     | 442         | 441      | negligible       |

### Q2: Multi-token `olap AND oltp` (appears in 476 rows)

| Metric            | Without FTS | With FTS | Improvement        |
| ----------------- | ----------- | -------- | ------------------ |
| query_duration_ms | 3,980 ms    | **4 ms** | **~1,000× faster** |
| read_bytes        | 7.76 GiB    | 6.47 MiB | **~1,200× less**   |
| marks scanned     | 3,091       | 416      | **87% skipped**    |

### Q3: Token filter + GROUP BY monthly aggregation

| Metric            | Without FTS | With FTS  | Improvement      |
| ----------------- | ----------- | --------- | ---------------- |
| query_duration_ms | 611 ms      | **6 ms**  | **~100× faster** |
| read_bytes        | 1.31 GiB    | 17.23 MiB | **~78× less**    |
| marks scanned     | 442         | 441       | negligible       |

### Summary

| Query pattern    | No FTS   | With FTS | Speedup     |
| ---------------- | -------- | -------- | ----------- |
| Single token     | 537 ms   | 2 ms     | **~270×**   |
| Multi-token      | 3,980 ms | 4 ms     | **~1,000×** |
| Token + GROUP BY | 611 ms   | 6 ms     | **~100×**   |

**Why the gap is so wide here**: `clickhouse` appears in only 0.004% of rows. The FTS index identifies those 1,145 rows instantly and reads only their data. Without FTS, ClickHouse must scan the entire 1.30 GiB of text data (or 7.76 GiB for the multi-token case) even though 99.996% of rows don't match.

## 7. How FTS works

ClickHouse FTS is an **inverted index**: a mapping from tokens to the row numbers that contain them. At query time, `hasToken()` looks up the token in this index and reads only the matching granules.

```sql
INDEX name column TYPE text(
    tokenizer = splitByNonAlpha,     -- how to split text into words
    preprocessor = lowerUTF8(column) -- optional transform before tokenization
)
```

Key points:

- `hasToken()` needles must be a single token (no spaces/separators).
- `LIKE '%term%'` and `match()` don't use this FTS index.
- The index is deterministic (no Bloom filter false positives).
- Best for selective terms: the rarer the match, the bigger the win.
- Write throughput is ~50% slower with FTS (trade-off).

## 8. Clean up

```bash
docker compose down -v
```

## References

- [ClickHouse Full-Text Search GA Announcement](https://clickhouse.com/blog/full-text-search-ga-release)
- [Text Index Documentation](https://clickhouse.com/docs/engines/table-engines/mergetree-family/textindexes)
- [Hacker News public dataset](https://clickhouse.com/docs/getting-started/example-datasets/hacker-news)
