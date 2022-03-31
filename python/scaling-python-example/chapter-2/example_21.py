"""Example 2.1. Starting a new thread.

Run this multiple times, output might be different each time.
"""
import threading


def print_something(something):
    print(something)


t = threading.Thread(target=print_something, args=("hello", ))
t.start()
print("thread started")
t.join() # Main thread waits for the 2nd thread to complete.
# If do not join all threads and wait for them to finish, it is possible
# that main thread finishs and exit before the other threads.
