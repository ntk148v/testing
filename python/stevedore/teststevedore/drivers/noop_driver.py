from teststevedore import base


class NoopDriver(base.BaseDriver):
    NAME = "noop"

    def do(self):
        print("Do nothing")
