from datetime import datetime
import os
import time

import futurist
from taskflow import engines
from taskflow.patterns import linear_flow as lf
from taskflow.patterns import unordered_flow as uf
from taskflow import task
import telegram
from dotenv import load_dotenv
from prettytable import PrettyTable

load_dotenv()

counter = 0


class Task1(task.Task):
    def execute(self, host_name, notification_uuid):
        global counter
        while counter <= 10:
            print("Task1: ", host_name, notification_uuid)
            time.sleep(2)
            counter += 1


class Task2(task.Task):
    def execute(self, host_name, notification_uuid):
        global counter
        bot = telegram.Bot(os.getenv("TELEGRAM_BOT_TOKEN"))
        while True:
            if counter > 10:
                break
            # Send message to telegram
            table = PrettyTable()
            table.field_names = ["City name", "Area",
                                 "Population", "Annual Rainfall"]
            table.add_row(["Adelaide", 1295, 1158259, 600.5])
            table.add_row(["Brisbane", 5905, 1857594, 1146.4])
            table.add_row(["Darwin", 112, 120900, 1714.7])
            table.add_row(["Hobart", 1357, 205556, 619.5])
            table.add_row(["Sydney", 2058, 4336374, 1214.8])
            table.add_row(["Melbourne", 1566, 3806092, 646.9])
            table.add_row(["Perth", 5386, 1554769, 869.4])

            msg = f"""
<b>🔥[Test][{host_name}][{notification_uuid}] Done {counter} times at {datetime.now()}</b>

<pre>
{table.get_string()}
</pre>
"""
            bot.send_message(chat_id=os.getenv(
                "TELEGRAM_CHAT_ID"), text=msg, parse_mode="HTML")

            time.sleep(5)  # This task supposes to take longer time


class Task3(task.Task):
    def execute(self, host_name, notification_uuid):
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
                             store={'host_name': host_name,
                                    'notification_uuid': notification_uuid})
            for st in e.run_iter():
                print(st)
        finally:
            executor.shutdown()


class Task4(task.Task):
    def execute(self, host_name, notification_uuid):
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
                             store={'host_name': host_name,
                                    'notification_uuid': notification_uuid},
                             max_workers=5)
        else:
            print("-- Running in parallel using eventlet --")
            e = engines.load(parent, engine='parallel', executor='greenthreaded',
                             store={'host_name': host_name,
                                    'notification_uuid': notification_uuid},
                             max_workers=5)
        e.run()


main = lf.Flow('main')
main.add(Task4())
engines.run(main, store={'host_name': 'computehost1',
                         'notification_uuid': 'yourphonenumber'})
