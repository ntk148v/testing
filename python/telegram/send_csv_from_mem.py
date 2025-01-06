from datetime import datetime
import os
import io
import csv
import json
import random
import uuid
import telegram
from dotenv import load_dotenv
from prettytable import PrettyTable

load_dotenv()

# data = [
#     {
#         "index": 1,
#         "project_id": "projectid",
#         "project_name": "projectname",
#         "domain_id": "domainid",
#         "domain_name": "domainname",
#         "instance_uuid": "instancelonguuid",
#         "instance_name": "vm1",
#         "instance_ipaddresses": "network - 192.168.1.253",
#         "source_host": "host1",
#         "dest_host": "host2",
#         "down_time": 10,
#         "status": "PENDING",
#     },
#     {
#         "index": 2,
#         "project_id": "projectid",
#         "project_name": "projectname",
#         "domain_id": "domainid",
#         "domain_name": "domainname",
#         "instance_uuid": "instancelonguuid",
#         "instance_name": "vm2",
#         "instance_ipaddresses": "network - 192.168.1.119",
#         "source_host": "host1",
#         "dest_host": "host4",
#         "down_time": 10,
#         "status": "PENDING",
#     },
#     {
#         "index": 3,
#         "project_id": "projectidprojectidprojectid",
#         "project_name": "projectnameprojectname",
#         "domain_id": "domainiddomainiddomainid",
#         "domain_name": "domainname",
#         "instance_uuid": "instancelonguuidinstancelonguuidinstancelonguuidinstancelonguuidinstancelonguuid",
#         "instance_name": "vm3",
#         "instance_ipaddresses": "network - 192.168.1.10",
#         "source_host": "host1",
#         "dest_host": "host5",
#         "down_time": 10,
#         "status": "PENDING",
#     }
# ]

# Function to generate a random IP address


def generate_ip():
    return f"192.168.{random.randint(0, 255)}.{random.randint(1, 254)}"


# Generate 100 samples
data = []
for i in range(1, 100):
    sample = {
        "index": i,
        "project_id": f"projectid_{i}",
        "project_name": f"projectname_{i}",
        "domain_id": f"domainid_{i}",
        "domain_name": f"domainname_{i}",
        "instance_uuid": str(uuid.uuid4()),
        "instance_name": f"vm{i}",
        "instance_ipaddresses": f"network - {generate_ip()}",
        "source_host": f"host{random.randint(1, 10)}",
        "dest_host": f"host{random.randint(1, 10)}",
        "down_time": random.randint(1, 100),
        "status": random.choice(["PENDING", "COMPLETED", "FAILED"])
    }
    data.append(sample)
strBuf = io.StringIO()
csv_writer = csv.DictWriter(
    strBuf, fieldnames=data[0].keys())
csv_writer.writeheader()
csv_writer.writerows(data)
# Seek to the beginning of the BytesIO buffer
strBuf.seek(0)

# python-telegram-bot library can send files only from io.BytesIO buffer
# we need to convert StringIO to BytesIO
bytesBuf = io.BytesIO()

# extract csv-string, convert it to bytes and write to buffer
bytesBuf.write(strBuf.getvalue().encode())
bytesBuf.seek(0)


def divide_into_batches(data, batch_size):
    # Divide the data into batches of the specified size
    return [data[i:i + batch_size] for i in range(0, len(data), batch_size)]


bot = telegram.Bot(os.getenv("TELEGRAM_BOT_TOKEN"))
try:
    bot.send_document(chat_id=os.getenv("TELEGRAM_CHAT_ID"), text="Long Text",
                      parse_mode="HTML", document=bytesBuf,
                      filename="test.csv")
    # Send as prettytable
    table = PrettyTable()
    table.field_names = data[0].keys()
    batches = divide_into_batches(data, 10)
    for batch in batches:
        for row in batch:
            table.add_row(row.values())
        msg = f"""
<b>🔥[Test]</b>

<pre>
{table.get_string()}
</pre>
        """
        bot.send_message(chat_id=os.getenv(
            "TELEGRAM_CHAT_ID"), text=msg, parse_mode="HTML")
        table.clear_rows()
except Exception as e:
    raise e
