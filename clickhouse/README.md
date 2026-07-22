# ClickHouse Full-Text Search Guide

> A hands-on guide to testing ClickHouse Full-Text Search (FTS) — from setup to benchmarking.

ClickHouse Full-Text Search is GA since [version 25.12+](https://clickhouse.com/blog/full-text-search-ga-release) (the blog announce details). It implements **inverted indexes** for fast token-based search over text columns, delivering 10–20x faster queries compared to full scans.

## Prerequisites

- Docker and Docker Compose
- ~10GB free disk

## 1. Start ClickHouse

```bash
docker compose up -d
```

The [docker-compose.yml](./docker-compose.yml) starts a single ClickHouse node (version ≥26.x, FTS is GA) with 8GB RAM limit, port 8123 (HTTP) and 9000 (TCP).

Verify:

```bash
docker exec clickhouse-server clickhouse-client -q "SELECT version()"
```

Expected output: `26.6.1.1193` (or newer).

## 2. Dataset

We use a synthetic **100M-row log dataset** (~3GB compressed). It mimics real-world log entries: mixed INFO/WARN/ERROR/CRITICAL messages, service names, hostnames, and trace IDs — with varying token frequencies (some terms in 1% of rows, some in 40%).

```bash
docker exec clickhouse-server clickhouse-client "
CREATE TABLE logs_bench (
    event_time DateTime,
    level String,
    service String,
    message String,
    host String,
    trace_id String
) ENGINE = MergeTree
ORDER BY (event_time, service, host);
"

# Insert 100M rows in batches
for batch in 1 2 3 4 5; do
  docker exec clickhouse-server clickhouse-client -q "
    INSERT INTO logs_bench
    SELECT
      toDateTime('2024-06-01 00:00:00') + rand() % (180*24*3600) AS event_time,
      multiIf(rand() % 100 < 5, 'CRITICAL', rand() % 100 < 15, 'ERROR',
              rand() % 100 < 30, 'WARN', rand() % 100 < 50, 'DEBUG', 'INFO') AS level,
      arrayRandomSample(['api-gw','auth','payment','notify','user','search',
                         'storage','worker','frontend','scheduler'], 1)[1] AS service,
      multiIf(
        rand() % 100 < 40,
          concat('Request OK status=200 duration=', toString(rand()%500+1), 'ms trace=', hex(rand()%1000000)),
        rand() % 100 < 60,
          concat('Slow query detected: query=', toString(rand()%1000), ' duration=', toString(rand()%5000+500), 'ms'),
        rand() % 100 < 75,
          concat('Auth failed user=', toString(rand()%10000), ' ip=', toString(rand()%255), '.',
                 toString(rand()%255), '.', toString(rand()%255), '.', toString(rand()%255)),
        rand() % 100 < 85,
          concat('Connection pool exhausted: service=', service, ' size=100 waiting=', toString(rand()%50)),
        rand() % 100 < 92,
          concat('OOM risk: heap=', toString(rand()%64+16), 'GB used=', toString(rand()%80+10), '%'),
        concat('buffer_overflow module=', toString(rand()%10), ' size=', toString(rand()%65535))
      ) AS message,
      concat('host-', toString(rand() % 5000)) AS host,
      concat(hex(rand() % 100000000), '-', hex(rand() % 100000)) AS trace_id
    FROM numbers(1, 20000000)
  "
done
```

Then create an **identical table with a Full-Text Search index**:

```bash
docker exec clickhouse-server clickhouse-client "
CREATE TABLE logs_bench_fts (
    event_time DateTime,
    level String,
    service String,
    message String,
    host String,
    trace_id String,
    INDEX fts_message message TYPE text(
        tokenizer = splitByNonAlpha,
        preprocessor = lowerUTF8(message)
    )
) ENGINE = MergeTree
ORDER BY (event_time, service, host);
"

# Copy data and materialize the index
docker exec clickhouse-server clickhouse-client "
  INSERT INTO logs_bench_fts SELECT * FROM logs_bench;
"
docker exec clickhouse-server clickhouse-client "
  ALTER TABLE logs_bench_fts MATERIALIZE INDEX fts_message;
"
```

### Index details

| Metric                | Value                                         |
| --------------------- | --------------------------------------------- |
| Rows                  | 100,000,000                                   |
| Table size (no index) | 3.09 GiB                                      |
| Table size (with FTS) | 3.65 GiB                                      |
| FTS index size        | 703 MiB (compressed) / 710 MiB (uncompressed) |
| Storage overhead      | ~18%                                          |
| Tokenizer             | `splitByNonAlpha`                             |
| Preprocessor          | `lowerUTF8(message)`                          |

## 3. Benchmarks

### 3.1. Single token: `'exhausted'` (appears in ~10% of rows)

```sql
-- WITH FTS
SELECT count() FROM logs_bench_fts WHERE hasToken(message, 'exhausted');

-- WITHOUT FTS
SELECT count() FROM logs_bench WHERE hasToken(message, 'exhausted');
```

| Metric     | Without FTS | With FTS | Improvement    |
| ---------- | ----------- | -------- | -------------- |
| Latency    | ~300 ms     | ~12 ms   | **25× faster** |
| Bytes read | 4.49 GiB    | 95 MiB   | **48× less**   |

### 3.2. Single token: `'Connection'` (appears in ~10% of rows)

```sql
SELECT count() FROM logs_bench_fts WHERE hasToken(message, 'Connection');
```

| Metric     | Without FTS | With FTS | Improvement    |
| ---------- | ----------- | -------- | -------------- |
| Latency    | ~300 ms     | ~13 ms   | **23× faster** |
| Bytes read | 4.49 GiB    | 95 MiB   | **48× less**   |

### 3.3. Multi-token aggregation query

```sql
SELECT service, count() AS cnt FROM logs_bench_fts
WHERE hasToken(message, 'exhausted')
GROUP BY service ORDER BY cnt DESC LIMIT 5;
```

| Metric     | Without FTS | With FTS | Improvement   |
| ---------- | ----------- | -------- | ------------- |
| Latency    | ~1200 ms    | ~200 ms  | **6× faster** |
| Bytes read | 5.85 GiB    | 1.03 GiB | **5.7× less** |

### Summary

| Query pattern                       | No FTS  | With FTS | Speedup  |
| ----------------------------------- | ------- | -------- | -------- |
| Single token scan (10% selectivity) | ~300ms  | ~12ms    | **~25×** |
| Single token scan (1% selectivity)  | ~300ms  | ~10ms    | **~30×** |
| Multi-token + GROUP BY              | ~1200ms | ~200ms   | **~6×**  |

**Key insight:** The more data you have, the bigger the gap. On our smaller 5M-row dataset the difference was negligible — FTS shines when scans take seconds, not milliseconds.

## 4. How FTS works

ClickHouse FTS creates an **inverted index**: a mapping from tokens (words) to the row numbers containing them. At query time, `hasToken()` uses this index to identify matching rows directly, skipping irrelevant granules entirely.

### Index configuration

```sql
INDEX name column TYPE text(
    tokenizer = splitByNonAlpha,   -- how to split text into tokens
    preprocessor = lowerUTF8(column)  -- optional transform before tokenization
)
```

Available tokenizers:

- `splitByNonAlpha` — splits on any non-alphabetic character
- `splitByWhitespace` — splits on whitespace only

Preprocessors enable case-insensitive search (via `lowerUTF8`) or other transforms before tokenization.

### Query functions

| Function                               | Behavior                                                  |
| -------------------------------------- | --------------------------------------------------------- |
| `hasToken(col, 'term')`                | Exact token match (no separators/whitespace in needle)    |
| `hasTokenCaseInsensitive(col, 'term')` | Case-insensitive token match (without index preprocessor) |
| `like` / `match`                       | Pattern / regex — won't use FTS index directly            |

The FTS index accelerates `hasToken()` and related token functions. For `LIKE '%term%'` patterns you still get a full scan (or use ngram tokenizer).

### When to use

- **Great fit:** Log search, observability, text-heavy analytics, filtering message blobs at scale
- **Not a fit:** Full-text relevance ranking (BM25/TF-IDF), phrase search with positional matching, substring search
- **vs Bloom filter:** FTS is deterministic (no false positives), row-level precise, dramatically faster at scale — but larger and slightly slower on writes

## 5. Clean up

```bash
docker compose down -v
```

## References

- [ClickHouse Full-Text Search GA Announcement](https://clickhouse.com/blog/full-text-search-ga-release)
- [Text Index Documentation](https://clickhouse.com/docs/engines/table-engines/mergetree-family/textindexes)
- [hasToken Function](https://clickhouse.com/docs/sql-reference/functions/string-search-functions#hasToken)
