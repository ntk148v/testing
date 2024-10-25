import logging
import time

import futurist
from taskflow import engines
from taskflow.patterns import linear_flow as lf
from taskflow.patterns import unordered_flow as uf
from taskflow import task


class PrinterTask(task.Task):
    def __init__(self, name, show_name=True, inject=None):
        super(PrinterTask, self).__init__(name, inject=inject)
        self._show_name = show_name

    def execute(self, output):
        if self._show_name:
            print("%s: %s" % (self.name, output))
        else:
            print(output)


class Task1(task.Task):
    def execute(self, hostname, notification_uuid):
        i = 0
        while i <= 5:
            print("Task1: ", hostname, notification_uuid)
            time.sleep(1)
            i += 1


class Task2(task.Task):
    def execute(self, hostname, notification_uuid):
        i = 0
        while i <= 5:
            print("Task2: ", hostname, notification_uuid)
            time.sleep(2)  # This task supposes to take longer time
            i += 1


class Task3(task.Task):
    def execute(self, hostname, notification_uuid):
        # Calling flow from task
        parent = uf.Flow('parent')
        parent.add(Task1(), Task2())

        # Now run it (using the specified executor)...
        try:
            print("-- Running in parallel using eventlet --")
            executor = futurist.GreenThreadPoolExecutor(max_workers=5)
        except RuntimeError:
            # No eventlet currently active, use real threads instead.
            print("-- Running in parallel using real threads --")
            executor = futurist.ThreadPoolExecutor(max_workers=5)
        try:
            e = engines.load(parent, engine='parallel', executor=executor,
                             store={'hostname': hostname,
                                    'notification_uuid': notification_uuid})
            for st in e.run_iter():
                print(st)
        finally:
            executor.shutdown()


class Task4(task.Task):
    def execute(self, hostname, notification_uuid):
        # Calling flow from task
        parent = uf.Flow('parent')
        parent.add(Task1(), Task2())

        # Now run it (using the specified executor)...
        try:
            import eventlet
        except ImportError:
            # No eventlet currently active, use real threads instead.
            print("-- Running in parallel using real threads --")
            e = engines.load(parent, engine='parallel', executor='threaded',
                             store={'hostname': hostname,
                                    'notification_uuid': notification_uuid},
                             max_workers=5)
        else:
            print("-- Running in parallel using eventlet --")
            e = engines.load(parent, engine='parallel', executor='greenthreaded',
                             store={'hostname': hostname,
                                    'notification_uuid': notification_uuid},
                             max_workers=5)
        e.run()


main = lf.Flow('main')
main.add(Task4())
engines.run(main, store={'hostname': 'computehost1',
                         'notification_uuid': 'yourphonenumber'})
