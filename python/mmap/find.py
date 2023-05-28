import mmap
import timeit

filename = input("Enter the absolute file path: ")
if not filename:
    filename = "/tmp/a"


def regular_io_find(filename):
    with open(filename, mode="r", errors="ignore") as file_obj:
        text = file_obj.read()
        return text.find(" the ")


def mmap_io_find(filename):
    with open(filename, mode="r", errors="ignore") as file_obj:
        with mmap.mmap(file_obj.fileno(), length=0, access=mmap.ACCESS_READ) as mmap_obj:
            return mmap_obj.find(b" the ")


print("regular_io_find", timeit.repeat(
    "regular_io_find(filename)",
    repeat=3,
    number=1,
    setup="from __main__ import regular_io_find, filename"))
print("mmap_io_find", timeit.repeat(
    "mmap_io_find(filename)",
    repeat=3,
    number=1,
    setup="from __main__ import mmap_io_find, filename"))
