import eventlet
from eventlet import greenpool

thread_pool = greenpool.GreenPool(size=4)

taskids = [1, 2, 3, 4]

def do(tid: str):
    print("Executing task", tid)
    eventlet.sleep(2)
    print("Done task", tid)

for tid in taskids:
    thread_pool.spawn_n(do, tid)

thread_pool.waitall()
