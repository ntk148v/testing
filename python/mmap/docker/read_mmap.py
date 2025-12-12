# Create file
with open("large.txt", "wb") as f:
    f.truncate(1024 * 1024 * 1024 * 1024)

import mmap


def mmap_io(filename):
    with open(filename, mode="r", encoding="utf8") as file_obj:
        with mmap.mmap(
            file_obj.fileno(), length=0, access=mmap.ACCESS_READ
        ) as mmap_obj:
            print(mmap_obj.find(b"abc"))


mmap_io("large.txt")
