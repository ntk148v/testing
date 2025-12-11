import os

# Create file
with open("trace_test.bin", "wb") as f:
    f.truncate(1024)

# 5 Operations using Seek + Write
with open("trace_test.bin", "r+b") as f:
    for i in range(5):
        f.seek(i * 10)  # SYSCALL
        f.write(b"X")  # SYSCALL
