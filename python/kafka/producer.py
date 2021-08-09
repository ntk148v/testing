import time
import json
import os

from kafka import KafkaProducer

kafka_host = os.getenv('KAFKA_HOST')
producer = KafkaProducer(bootstrap_servers=[kafka_host],
                         value_serializer=lambda x: json.dumps(x).encode('utf-8'))

for e in range(1000):
    data = {'number' : e}
    print('Send {}'. format(data))
    producer.send('numtest', value=data)
    time.sleep(5)
