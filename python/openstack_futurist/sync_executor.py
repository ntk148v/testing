import datetime
import sys
import futurist

import eventlet

def delayed_func():
    print("started")
    eventlet.sleep(3)
    print("done")

print(datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
e = futurist.SynchronousExecutor()
fut = e.submit(delayed_func)
eventlet.sleep(1)
print("Hello")
eventlet.sleep(1)
e.shutdown()
print(datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
