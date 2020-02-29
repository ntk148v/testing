import futurist
from futurist import periodics

import time
import threading


@periodics.periodic(1)
def run_only_once(started_at):
    print("1: %s" % (time.time() - started_at))
    raise periodics.NeverAgain("No need to run again after first run !!")


@periodics.periodic(1)
def keep_running(started_at):
    print("2: %s" % (time.time() - started_at))


started_at = time.time()
callables = [
    # The function to run + any automatically provided positional and
    # keyword arguments to provide to it everytime it is activated.
    (run_only_once, (started_at,), {}),
    (keep_running, (started_at,), {}),
]

executor_factory = lambda: futurist.ThreadPoolExecutor(max_workers=2)
w = periodics.PeriodicWorker(callables, executor_factory=executor_factory)

t = threading.Thread(target=w.start)
t.daemon = True
t.start()

# Run for 10 seconds and then stop.
while (time.time() - started_at) <= 10:
    time.sleep(0.1)
w.stop()
w.wait()
t.join()
