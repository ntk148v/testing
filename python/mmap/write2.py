import mmap
import time
import os

# Define the number of lines to write
num_lines = 10000
filename1 = "/tmp/log_file.txt"
filename2 = "/tmp/log_mmap.txt"

def regular_io_write(filename):
    with open(filename, mode="w", encoding="utf-8") as f:
        f.write(f"non empty string\n")
        for i in range(num_lines):
            f.write(f'This is line {i}\n')

        f.flush()

def mmap_io_write(filename):
    # mmap module doesn’t allow memory mapping of an empty file
    with open(filename, mode="w", encoding="utf-8") as f:
        f.write(f"non empty string\n")
        f.flush()

    with open(filename, mode="r+", encoding="utf-8")  as f:
        size = f.seek(0,2) # Seek relative to end of file
        with mmap.mmap(f.fileno(), size, access=mmap.ACCESS_WRITE) as m:
            for i in range(num_lines):
                m.write(f'This is line {i}\n'.encode('utf-8'))


# Benchmark regular file write
start = time.time()
regular_io_write(filename1)
end = time.time()
print(f'Regular file write took {end - start:.2f} seconds')


# Benchmark mmap write
start = time.time()
mmap_io_write(filename2)
end = time.time()
print(f'Mmap write took {end - start:.2f} seconds')
