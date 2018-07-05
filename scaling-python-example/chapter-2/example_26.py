"""Example 2.6: Worker using multipleprocessing"""
import multiprocessing
import random


def compute(n):
    return sum(
        [random.randint(1, 100) for i in range(100000)])


# Start 8 worker
pool = multiprocessing.Pool(processes=8)
print("Results: %s" % pool.map(compute, range(8)))
