import grpc

import my_service_pb2
import my_service_pb2_grpc


class gRPCClient():
    def __init__(self):
        channel = grpc.insecure_channel('127.0.0.1:50051')
        self.stub = my_service_pb2_grpc.MyServiceStub(channel)

    def method1(self, name, code):
        print('method 1')
        return self.stub.MyMethod1(my_service_pb2.MyRequest(name=name, code=code))

    def method2(self, name, code):
        print('method 2')
        return self.stub.MyMethod2(my_service_pb2.MyRequest(name=name, code=code))


def main():
    print('main')
    client = gRPCClient()
    print(client.method1('Alexandre', 123))
    print(client.method2('Maria', 123))


if __name__ == '__main__':
    main()
