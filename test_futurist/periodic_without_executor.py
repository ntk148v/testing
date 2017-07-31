import futurist
from futurist import periodics

import time
import threading


@periodics.periodic(1)
def every_one(started_at):
    print("1: %s" % (time.time() - started_at))


@periodics.periodic(2)
def every_two(started_at):
    print("2: %s" % (time.time() - started_at))


@periodics.periodic(4)
def every_four(started_at):
    print("4: %s" % (time.time() - started_at))


@periodics.periodic(6)
def every_six(started_at):
    print("6: %s" % (time.time() - started_at))


started_at = time.time()
callables = [
    # The function to run + any automatically provided positional and
    # keyword arguments to provide to it everytime it is activated.
    (every_one, (started_at,), {}),
    (every_two, (started_at,), {}),
    (every_four, (started_at,), {}),
    (every_six, (started_at,), {}),
]
w = periodics.PeriodicWorker(callables)

# In this example we will run the periodic functions using a thread, it
# is also possible to just call the w.start() method directly if you do
# not mind blocking up the current program.
t = threading.Thread(target=w.start)
t.daemon = True
t.start()

# Run for 10 seconds and then stop.
while (time.time() - started_at) <= 10:
    time.sleep(0.1)
w.stop()
w.wait()
t.join()
