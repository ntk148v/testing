"""Example 2.7 Worker using concurrent.futures.ThreadPoolExecutor"""
from concurrent import futures
import random


def compute():
    return sum(
        [random.randint(1, 100) for i in range(100000)])


with futures.ThreadPoolExecutor(max_workers=8) as executor:
    futures = [executor.submit(compute) for _ in range(8)]

results = [f.result() for f in futures]

print("Results: %s" % results)
