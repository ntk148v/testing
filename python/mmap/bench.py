import os
import mmap
import time
import random

FILENAME = "perf_test.bin"
FILE_SIZE_MB = 1024
FILE_SIZE = FILE_SIZE_MB * 1024 * 1024
NUM_OPERATIONS = 100_000


def setup_file():
    """Creates a sparse file of the specified size instantly."""
    print(f"Creating {FILE_SIZE_MB}MB file: {FILENAME}...")
    with open(FILENAME, "wb") as f:
        f.truncate(FILE_SIZE)

    # Generate random positions ahead of time so the RNG calculation
    # doesn't affect the I/O benchmark time.
    print(f"Generating {NUM_OPERATIONS} random indices...")
    indices = [random.randint(0, FILE_SIZE - 1) for _ in range(NUM_OPERATIONS)]
    return indices


def benchmark_old_way(indices):
    print("\n--- Starting STANDARD I/O Benchmark (seek + write) ---")
    start_time = time.perf_counter()

    # For every single byte written, halts, hands control to the Linux kernel
    # seek(), waits for the Kernel, halts again (write), hands control to
    # the kernel, and waits again
    with open(FILENAME, "r+b") as f:
        for idx in indices:
            # SYSCALL 1: Tell OS to move the disk head/pointer
            f.seek(idx)
            # SYSCALL 2: Copy byte from user to kernel
            f.write(b"X")

    end_time = time.perf_counter()
    duration = end_time - start_time
    print(f"Standard I/O: {duration:.4f} seconds")
    return duration


def benchmark_mmap_way(indices):
    print("\n--- Starting MMAP Benchmark (memory assignment) ---")
    start_time = time.perf_counter()

    with open(FILENAME, "r+b") as f:
        # Map the file into memory
        # Treats file as a RAM array. The Kernel is only involved if it
        # needs to fetch a new page (Page Fault), but since we are just writing
        # to the memory, the CPU executes it at RAM speeds
        with mmap.mmap(f.fileno(), 0, access=mmap.ACCESS_WRITE) as mm:
            for idx in indices:
                # NO SYSCALL: Just writing to a memory address
                mm[idx] = 88  # 88 is ASCII for 'X'

    end_time = time.perf_counter()
    duration = end_time - start_time
    print(f"Mmap I/O:     {duration:.4f} seconds")
    return duration


def main():
    try:
        indices = setup_file()

        # Run benchmarks
        t_old = benchmark_old_way(indices)
        t_mmap = benchmark_mmap_way(indices)

        # Results
        print("\n--- RESULTS ---")
        print(f"Standard I/O Operations per sec: {NUM_OPERATIONS / t_old:.0f}")
        print(f"Mmap Operations per sec:         {NUM_OPERATIONS / t_mmap:.0f}")
        print(f"Speedup Factor:                  {t_old / t_mmap:.2f}x FASTER")

    finally:
        # Cleanup
        if os.path.exists(FILENAME):
            os.remove(FILENAME)
            print(f"\nCleaned up {FILENAME}")


if __name__ == "__main__":
    main()
