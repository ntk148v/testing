import time
from concurrent import futures

import grpc

import my_service_pb2
import my_service_pb2_grpc

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class gRPCServer(my_service_pb2_grpc.MyServiceServicer):
    def __init__(self):
        print('Initilization')

    def MyMethod1(self, request, context):
        print(request.name)
        print(request.code)
        return my_service_pb2.MyResponse(name=request.name,
                                         sex='M', code=123)

    def MyMethod2(self, request, context):
        print(request.name)
        print(request.code * 12)
        return my_service_pb2.MyResponse(name=request.name,
                                         sex='F', code=1234)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    my_service_pb2_grpc.add_MyServiceServicer_to_server(
        gRPCServer(), server
    )
    server.add_insecure_port('[::]:50051')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
