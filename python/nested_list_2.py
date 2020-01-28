import collections
import copy
from functools import reduce
from itertools import chain
import operator
import time


def flatten1(l):
    return list(chain.from_iterable(l))


def flatten2(l):
    return list(chain(*l))


def flatten3(l):
    return reduce(lambda x, y: x+y, l)


def flatten4(l):
    return reduce(operator.concat, l) if isinstance(l, list) else [l]


def flatten5(l):
    return sum(l, [])


def flatten6(l):
    return reduce(operator.concat, l)


def execution_time(fn, args):
    t = time.process_time()
    print(fn(args))
    print(time.process_time() - t)


if __name__ == "__main__":
    xs = [[1], [2, 3], [4, 5]]
    execution_time(flatten1, xs)
    execution_time(flatten2, xs)
    execution_time(flatten3, xs)
    execution_time(flatten4, xs)
    execution_time(flatten5, xs)
    execution_time(flatten6, xs)
