from teststevedore import base


class EchoDriver(base.BaseDriver):
    NAME = "echo"

    def do(self):
        print("Do echo thing")
