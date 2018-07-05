"""Example 2.2: Starting a new thread in daemon mode.
"""
import threading


def print_something(something):
    print(something)


t = threading.Thread(target=print_something, args=("hello",))
t.daemon = True
t.start()
print("thread started")
