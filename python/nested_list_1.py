import collections
import time


def flatten1(l):
    result = []
    for i in l:
        if isinstance(i, collections.Iterable):
            result.extend(flatten1(i))
        else:
            result.append(i)
    return result


def flatten2(l):
    if isinstance(l, collections.Iterable):
        return [a for i in l for a in flatten2(i)]
    else:
        return [l]


def flatten3(l):
    if not isinstance(l, collections.Iterable):
        yield l
    else:
        for i in l:
            yield from flatten3(i)


def execution_time(fn, args):
    t = time.process_time()
    print(fn(args))
    print(time.process_time() - t)


if __name__ == "__main__":
    xs = [1, [2, 3, [4, 5]]]
    execution_time(flatten1, xs)
    execution_time(flatten2, xs)
    t = time.process_time()
    print(list(flatten3(xs)))
    print(time.process_time() - t)
