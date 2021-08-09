import json

import msgpack
from kafka import KafkaProducer
from kafka.errors import KafkaError

producer = KafkaProducer(bootstrap_servers=['localhost:9092'],
                         value_serializer=lambda x:
                         json.dumps(x).encode('utf-8'))

# Asynchronous by default
future = producer.send('common-topic', 'raw_bytes')

# Block for 'synchronous' sends
try:
    record_metadata = future.get(timeout=10)
except KafkaError as e:
    # Decide what to do if produce request failed
    raise(e)

# Successful result returns assigned partition and offset
print(record_metadata.topic)
print(record_metadata.partition)
print(record_metadata.offset)

# produce keyed messages to enable hashed partitioning
producer.send('common-topic', key='boo', value='bar')

# encode objects via msgpack
producer = KafkaProducer(value_serializer=msgpack.dumps)
producer.send('msgspack-topic', {'key': 'value'})

# produce asynchronously
for _ in range(100):
    producer.send('common-topic', 'msg')


def on_send_success(record_metadata):
    print(record_metadata.topic)
    print(record_metadata.partition)
    print(record_metadata.offset)


def on_send_error(excp):
    print('I am an errback', exc_info=excp)
    # handle exception


# produce asynchronously with callbacks
producer.send('common-topic', 'raw_bytes').add_callback(
    on_send_success).add_errback(on_send_error)

# block until all async messages are sent
producer.flush()

# configure multiple retries
producer = KafkaProducer(retries=5)
