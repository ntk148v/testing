# ClickHouse Full-Text Search Guide

Hands-on test plan for ClickHouse Full-Text Search (FTS) on a single-node Docker setup.

## Prerequisites

- Docker and Docker Compose
- ~10GB free disk
- Network access to ClickHouse's public dataset bucket

## 1. Start ClickHouse

```bash
docker compose up -d
```

Verify the server version:

```bash
docker exec clickhouse-server clickhouse-client -q "SELECT version()"
```

Use ClickHouse 25.12+ for GA Full-Text Search. The current image tested here was `26.6.1.1193`.

## 2. Import a public dataset

Use ClickHouse's public **Hacker News** dataset:

- Source: `https://datasets-documentation.s3.eu-west-3.amazonaws.com/hackernews/hacknernews.parquet`
- Size: ~28M rows
- Why this dataset: public, text-heavy (`text`, `title`), and small enough for this host's 8GB ClickHouse memory limit.

Create the table without an index first:

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

## 3. Benchmark before FTS

Run each query 3 times and record the best result from `system.query_log`.

```bash
# Query 1: single token
docker exec clickhouse-server clickhouse-client -q "
SELECT count()
FROM hackernews
WHERE hasToken(lowerUTF8(text), 'clickhouse')
SETTINGS enable_full_text_index = 0;
"

# Query 2: multi-token filter
docker exec clickhouse-server clickhouse-client -q "
SELECT count()
FROM hackernews
WHERE hasToken(lowerUTF8(text), 'olap')
  AND hasToken(lowerUTF8(text), 'oltp')
SETTINGS enable_full_text_index = 0;
"

# Query 3: token filter + aggregation
docker exec clickhouse-server clickhouse-client -q "
SELECT toYYYYMM(time) AS month, count() AS mentions
FROM hackernews
WHERE hasToken(lowerUTF8(text), 'clickhouse')
GROUP BY month
ORDER BY month
SETTINGS enable_full_text_index = 0;
"
```

Collect timings:

```bash
docker exec clickhouse-server clickhouse-client -q "SYSTEM FLUSH LOGS"

docker exec clickhouse-server clickhouse-client -q "
SELECT
    query_duration_ms,
    read_rows,
    formatReadableSize(read_bytes) AS read_bytes,
    query
FROM system.query_log
WHERE type = 'QueryFinish'
  AND current_database = 'default'
  AND query LIKE '%hackernews%'
  AND query LIKE '%hasToken%'
ORDER BY event_time DESC
LIMIT 10
FORMAT Vertical;
"
```

## 4. Add Full-Text Search

Add the FTS index to the **same table** so the before/after benchmark only changes one thing: index usage.

```bash
docker exec clickhouse-server clickhouse-client --multiquery "
ALTER TABLE hackernews ADD INDEX fts_text text TYPE text(
    tokenizer = splitByNonAlpha,
    preprocessor = lowerUTF8(text)
);

ALTER TABLE hackernews MATERIALIZE INDEX fts_text;
"
```

Wait for materialization:

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT command, is_done, latest_fail_reason
FROM system.mutations
WHERE database = 'default' AND table = 'hackernews'
ORDER BY create_time DESC;
"
```

Check index size:

```bash
docker exec clickhouse-server clickhouse-client -q "
SELECT
    name,
    type,
    formatReadableSize(data_compressed_bytes) AS compressed,
    formatReadableSize(data_uncompressed_bytes) AS uncompressed
FROM system.data_skipping_indices
WHERE database = 'default' AND table = 'hackernews';
"
```

Confirm ClickHouse plans to use the index:

```bash
docker exec clickhouse-server clickhouse-client -q "
EXPLAIN indexes = 1
SELECT count()
FROM hackernews
WHERE hasToken(text, 'clickhouse')
SETTINGS enable_full_text_index = 1;
"
```

## 5. Benchmark after FTS

Use the same predicates as the baseline, but enable FTS. Because the index already lowercases `text`, query the original column.

```bash
# Query 1: single token
docker exec clickhouse-server clickhouse-client -q "
SELECT count()
FROM hackernews
WHERE hasToken(text, 'clickhouse')
SETTINGS enable_full_text_index = 1;
"

# Query 2: multi-token filter
docker exec clickhouse-server clickhouse-client -q "
SELECT count()
FROM hackernews
WHERE hasToken(text, 'olap')
  AND hasToken(text, 'oltp')
SETTINGS enable_full_text_index = 1;
"

# Query 3: token filter + aggregation
docker exec clickhouse-server clickhouse-client -q "
SELECT toYYYYMM(time) AS month, count() AS mentions
FROM hackernews
WHERE hasToken(text, 'clickhouse')
GROUP BY month
ORDER BY month
SETTINGS enable_full_text_index = 1;
"
```

Collect timings again from `system.query_log` using the same query from section 3.

## 6. Results table

Fill this in from `system.query_log`, not shell `time`, so client startup overhead is excluded.

| Query | FTS | query_duration_ms | read_rows | read_bytes |
| --- | --- | ---: | ---: | ---: |
| single token: `clickhouse` | off |  |  |  |
| single token: `clickhouse` | on |  |  |  |
| multi-token: `olap AND oltp` | off |  |  |  |
| multi-token: `olap AND oltp` | on |  |  |  |
| token + monthly aggregation | off |  |  |  |
| token + monthly aggregation | on |  |  |  |

Expected outcome: FTS should reduce `read_bytes` and latency most for selective terms. Common terms may improve less because ClickHouse still has many matching rows to read.

## 7. Notes

- FTS is for token filtering, not relevance ranking like BM25/TF-IDF.
- `hasToken()` needles must be one token; no spaces or separator characters.
- `LIKE '%term%'` and regex queries do not directly use this FTS index.
- For fair no-index comparisons, use `SETTINGS enable_full_text_index = 0` on the same table.

## Clean up

```bash
docker compose down -v
```

## References

- [ClickHouse Full-Text Search GA Announcement](https://clickhouse.com/blog/full-text-search-ga-release)
- [Text Index Documentation](https://clickhouse.com/docs/engines/table-engines/mergetree-family/textindexes)
- [Hacker News public dataset](https://clickhouse.com/docs/getting-started/example-datasets/hacker-news)
