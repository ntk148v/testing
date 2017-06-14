import oslo_messaging as messaging
import conf

version = '1.0'
version_cap = '1.2'

transport = messaging.get_transport(conf.CONF)
target = messaging.Target(topic='topic', version=version)

client = messaging.RPCClient(transport, target, version_cap=version_cap)
client.call({}, 'print_message', message='this is test')