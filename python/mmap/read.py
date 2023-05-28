import timeit
import mmap


def regular_io(filename):
    with open(filename, mode="r", errors='ignore') as file_obj:
        text = file_obj.read()
        return text


def mmap_io(filename):
    with open(filename, mode="r", errors='ignore') as file_obj:
        with mmap.mmap(file_obj.fileno(), length=0, access=mmap.ACCESS_READ) as mmap_obj:
            text = mmap_obj.read()
            return text


filename = input("Enter the absolute file path: ")
if not filename:
    filename = "/tmp/a"

# Benchmark time
print("regular_io", timeit.repeat(
    "regular_io(filename)",
    repeat=3,
    number=1,
    setup="from __main__ import regular_io, filename"))

print("mmap_io", timeit.repeat(
    "mmap_io(filename)",
    repeat=3,
    number=1,
    setup="from __main__ import mmap_io, filename"))

# regular_io [6.941209831999913, 7.164341439000054, 6.920215922000352]
# mmap_io [0.2982539229997201, 0.2962546029993973, 0.2988888530007898]
