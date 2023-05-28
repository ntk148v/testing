import timeit
import os
import shutil
import mmap

filename = input("Enter the absolute file path: ")
if not filename:
    filename = "/tmp/a"


def mmap_io_write(filename):
    with open(filename, mode="r+") as file_obj:
        #  ACCESS_WRITE specifies write-through semantics, meaning the data will be written through memory and persisted on disk.
        # ACCESS_COPY does not write the changes to disk, even if flush() is called.
        with mmap.mmap(file_obj.fileno(), length=0, access=mmap.ACCESS_WRITE) as mmap_obj:
            mmap_obj[10:16] = b"python"
            mmap_obj.flush()


def regular_io_find_and_replace(filename):
    with open(filename, "r", errors="ignore") as orig_file_obj:
        with open("/tmp/b", "w", errors="ignore") as new_file_obj:
            orig_text = orig_file_obj.read()
            new_text = orig_text.replace(" the ", " eht ")
            new_file_obj.write(new_text)

    shutil.copyfile("/tmp/b", filename)
    os.remove("/tmp/b")


def mmap_io_find_and_replace(filename):
    with open(filename, mode="r+", errors="ignore") as file_obj:
        with mmap.mmap(file_obj.fileno(), length=0, access=mmap.ACCESS_WRITE) as mmap_obj:
            orig_text = mmap_obj.read()
            new_text = orig_text.replace(b" the ", b" eht ")
            mmap_obj[:] = new_text
            mmap_obj.flush()


print("regular_io_find_and_replace", timeit.repeat(
    "regular_io_find_and_replace(filename)",
    repeat=3,
    number=1,
    setup="from __main__ import regular_io_find_and_replace, filename"))

print("mmap_io_find_and_replace", timeit.repeat(
    "mmap_io_find_and_replace(filename)",
    repeat=3,
    number=1,
    setup="from __main__ import mmap_io_find_and_replace, filename"))
