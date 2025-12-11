import os, mmap

# Create file
with open("trace_test.bin", "wb") as f:
    f.truncate(1024)

# 5 Operations using Memory Assignment
with open("trace_test.bin", "r+b") as f:
    with mmap.mmap(f.fileno(), 0) as mm:
        for i in range(5):
            mm[i * 10] = 88  # NO SYSCALL
