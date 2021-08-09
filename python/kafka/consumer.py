import json
import time
import os

from kafka import KafkaConsumer

kafka_host = os.getenv('KAFKA_HOST')
consumer = KafkaConsumer(
    'numtest',
     bootstrap_servers=[kafka_host],
     auto_offset_reset='earliest',
     enable_auto_commit=True,
     group_id='my-group',
     value_deserializer=lambda x: json.loads(x.decode('utf-8')))
for message in consumer:
    message = message.value
    print('Receive {}'.format(message))
    time.sleep(5)
