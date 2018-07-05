"""Example 2.3: Workers using multithreading"""
import random
import threading


results = []

def compute():
    results.append(sum(
        [random.randint(1, 100) for i in range(100000)]))


workers = [threading.Thread(target=compute) for x in range(8)]
for worker in workers:
    worker.start()

for worker in workers:
    worker.join()

print("Results: %s" % results)
# Run time python example_23.py --> CPU usage around 150%, not as expect 400%
