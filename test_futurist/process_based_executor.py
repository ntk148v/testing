import time

import futurist

def delayed_func():
    print('started')
    time.sleep(3)
    print('done')
    # return "hello"

e = futurist.ProcessPoolExecutor()
fut = e.submit(delayed_func)
#print(fut.result())
print('Hello')
e.shutdown()
