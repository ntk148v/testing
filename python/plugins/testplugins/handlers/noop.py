from testplugins.handlers import base


class NoopHandler(base.BaseHandler):
    NAME = "handler"

    def do(self):
        print("Do Noop thing")

def handler():
    return NoopHandler()
