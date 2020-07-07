import importlib
import os
import pkgutil

import testplugins.handlers
from testplugins.handlers import base


def load_handlers():
    modules = {}
    iter_modules = pkgutil.iter_modules(testplugins.handlers.__path__)
    for (_, name, ispkg) in iter_modules:
        if not ispkg:
            modules[name] = getattr(importlib.import_module(
                'testplugins.handlers.' + name), 'handler')
    return modules


if __name__ == "__main__":
    print(load_handlers())
    for v in load_handlers().values():
        print(v().do())
    print(base.BaseHandler.__subclasses__())
