import futurist
from futurist import periodics

import time
import threading


@periodics.periodic(1)
def run_only_once(started_at):
    print("1: %s" % (time.time() - started_at))
    raise periodics.NeverAgain("No need to run again after first run !!")


@periodics.periodic(2)
def run_for_some_time(started_at):
    print("2: %s" % (time.time() - started_at))
    if (time.time() - started_at) > 5:
        raise periodics.NeverAgain("No need to run again !!")


started_at = time.time()
callables = [
    # The function to run + any automatically provided positional and
    # keyword arguments to provide to it everytime it is activated.
    (run_only_once, (started_at,), {}),
    (run_for_some_time, (started_at,), {}),
]
w = periodics.PeriodicWorker(callables)

# In this example we will run the periodic functions using a thread, it
# is also possible to just call the w.start() method directly if you do
# not mind blocking up the current program.
t = threading.Thread(target=w.start, kwargs={'auto_stop_when_empty': True})
t.daemon = True
t.start()

# Run for 10 seconds and then check to find out that it had
# already stooped.
while (time.time() - started_at) <= 10:
    time.sleep(0.1)
print(w.pformat())
t.join()
