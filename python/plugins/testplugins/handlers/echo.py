from testplugins.handlers import base


class EchoHandler(base.BaseHandler):
    NAME = "echo"

    def do(self):
        print("Do Echo thing")


def handler():
    return EchoHandler()
