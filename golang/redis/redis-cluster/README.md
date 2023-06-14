# Redis cluster

- Test Redis cluster
- Create Redis cluster with 9 nodes (but replicas 1 -> able to recover from fail)
- Use Redisinsight to discover master nodes ip address.
- Map ip address (internal) with container name.

```bash
docker inspect --format '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}} {{ .Name  }}' $(docker ps -aq)
```

- Stop redis master node.
- See how it works.
