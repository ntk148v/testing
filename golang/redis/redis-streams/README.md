# Redis streams

```bash
$ export REDIS_ADDR="localhost:6379"
$ export REDIS_STREAM="test"
# Run producer
$ go run producer/main.go
# You can run multiple consumers
$ go run consumer/main.go
```
