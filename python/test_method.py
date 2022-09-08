"""This is test about class method and static method"""

class Kls(object):
    def __init__(self, data):
        self.data = data

    def printd(self):
        print self.data

    @staticmethod
    def smethod(*args):
        print 'static', args

    @classmethod
    def cmethod(*args):
        print 'class', args


ins = Kls('test')
ins.printd()
Kls.printd()
ins.smethod('test')
Kls.smethod('test')
ins.cmethod('test')
Kls.cmethod('test')

