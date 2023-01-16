import asyncio
import random


class Worker:
    def __init__(self, num, target_q):
        # Worker number
        self.num = num
        # The job queue
        self.target_q = target_q
        # Create asyncio task
        self.task = asyncio.create_task(self.run())

    async def run(self):
        # Work on jobs until task is cancelled
        while True:
            print(f"Worker {self.num}: Waiting for new target")

            # Receive a new job from the queue
            target = await self.target_q.get()
            print(f"Worker {self.num}: Processing target {target}")

            try:
                # Simulating some work (e.g. a POST request)
                await asyncio.sleep(1.0)

                # Depending on the outcome, some new work results
                # 0-2 new targets are generated
                new_target_count = random.randint(0, 2)

                if new_target_count > 0:
                    print(f"Worker {self.num}: Target {target} generating "
                          f"{new_target_count} more targets")
                    for _ in range(new_target_count):
                        # Create a new random target
                        new_target = random.randint(1, 10000)

                        # Put new targets into queue. This will wait if queue is
                        # currently full.
                        await self.target_q.put(new_target)

            finally:
                print(f"Worker {self.num}: Target {target} done")

                # Decrease queue count by one
                self.target_q.task_done()


class Manager:
    def __init__(self):
        # A common queue holding the jobs for the workers. It just stores
        # integers here but could hold any data.
        self.target_q = asyncio.Queue(10)
        # The list of workers
        self.workers = None

    async def run(self):
        # Create some initial work
        await self.target_q.put(1)

        # Create 3 workers
        self.workers = [Worker(num, self.target_q) for num in range(3)]

        # Wait until queue of unfinished tasks is empty
        await self.target_q.join()

        # Cancel other workers
        for worker in self.workers:
            worker.task.cancel()


if __name__ == "__main__":
    manager = Manager()
    asyncio.run(manager.run())
