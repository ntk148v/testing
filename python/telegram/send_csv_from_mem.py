from datetime import datetime
import os
import io
import csv
import telegram
from dotenv import load_dotenv
from prettytable import PrettyTable

load_dotenv()

data = [
    {
        "index": 1,
        "project_id": "projectid",
        "project_name": "projectname",
        "domain_id": "domainid",
        "domain_name": "domainname",
        "instance_uuid": "instancelonguuid",
        "instance_name": "vm1",
        "instance_ipaddresses": "network - 192.168.1.253",
        "source_host": "host1",
        "dest_host": "host2",
        "down_time": 10,
        "status": "PENDING",
    },
    {
        "index": 2,
        "project_id": "projectid",
        "project_name": "projectname",
        "domain_id": "domainid",
        "domain_name": "domainname",
        "instance_uuid": "instancelonguuid",
        "instance_name": "vm2",
        "instance_ipaddresses": "network - 192.168.1.119",
        "source_host": "host1",
        "dest_host": "host4",
        "down_time": 10,
        "status": "PENDING",
    },
    {
        "index": 3,
        "project_id": "projectidprojectidprojectid",
        "project_name": "projectnameprojectname",
        "domain_id": "domainiddomainiddomainid",
        "domain_name": "domainname",
        "instance_uuid": "instancelonguuidinstancelonguuidinstancelonguuidinstancelonguuidinstancelonguuid",
        "instance_name": "vm3",
        "instance_ipaddresses": "network - 192.168.1.10",
        "source_host": "host1",
        "dest_host": "host5",
        "down_time": 10,
        "status": "PENDING",
    }
]
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

bot = telegram.Bot(os.getenv("TELEGRAM_BOT_TOKEN"))
bot.send_document(chat_id=os.getenv("TELEGRAM_CHAT_ID"), text="Long Text",
                  parse_mode="HTML", document=bytesBuf,
                  filename="test.csv")
# Send as prettytable
table = PrettyTable()
table.field_names = data[0].keys()
for row in data:
    table.add_row(row.values())
msg = f"""
<b>🔥[Test]</b>

<pre>
{table.get_string()}
</pre>
"""
bot.send_message(chat_id=os.getenv(
    "TELEGRAM_CHAT_ID"), text=msg, parse_mode="HTML")
