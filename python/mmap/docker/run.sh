#!/bin/bash

echo "Read regular"
docker run --rm -m 100m -v $PWD:/tmp/test python python /tmp/test/read_regular.py

echo "MMAP"
docker run --rm -m 100m -v $PWD:/tmp/test python python /tmp/test/read_mmap.py
