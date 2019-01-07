# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import my_service_pb2 as my__service__pb2


class MyServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.MyMethod1 = channel.unary_unary(
        '/testing.MyService/MyMethod1',
        request_serializer=my__service__pb2.MyRequest.SerializeToString,
        response_deserializer=my__service__pb2.MyResponse.FromString,
        )
    self.MyMethod2 = channel.unary_unary(
        '/testing.MyService/MyMethod2',
        request_serializer=my__service__pb2.MyRequest.SerializeToString,
        response_deserializer=my__service__pb2.MyResponse.FromString,
        )


class MyServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def MyMethod1(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def MyMethod2(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_MyServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'MyMethod1': grpc.unary_unary_rpc_method_handler(
          servicer.MyMethod1,
          request_deserializer=my__service__pb2.MyRequest.FromString,
          response_serializer=my__service__pb2.MyResponse.SerializeToString,
      ),
      'MyMethod2': grpc.unary_unary_rpc_method_handler(
          servicer.MyMethod2,
          request_deserializer=my__service__pb2.MyRequest.FromString,
          response_serializer=my__service__pb2.MyResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'testing.MyService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
