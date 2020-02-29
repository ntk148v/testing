import oslo_messaging as messaging
import conf

version = '1.2'

transport = messaging.get_transport(conf.CONF)
target = messaging.Target(topic='topic', server='server')


class RPCEndpoint(object):
    def print_message(self, ctxt, message):
        print('Received: ' + message)
        return message


access_policy = messaging.rpc.dispatcher.DefaultRPCAccessPolicy
endpoints = [RPCEndpoint(), ]
server = messaging.get_rpc_server(transport, target, endpoints,
                                  executor='blocking',
                                  access_policy=access_policy)
server.start()
server.wait()
