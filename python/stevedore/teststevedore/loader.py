NAMESPACE = "teststevedore.drivers"


def get_available_drivers(invoke_on_load=False):
    """
    Return names of the available drivers
    """
    from stevedore.extension import ExtensionManager
    manager = ExtensionManager(
        namespace=NAMESPACE, invoke_on_load=invoke_on_load)
    return manager.names()


def get_driver(name, invoke_on_load=False):
    """
    Retrieve a driver
    """
    from stevedore.driver import DriverManager
    manager = DriverManager(namespace=NAMESPACE, name=name,
                            invoke_on_load=invoke_on_load)
    return manager.driver


def get_driver_instance(name, invoke_on_load=False):
    """
    Retrieve a class instance for the driver
    """
    cls = get_driver(name, invoke_on_load=invoke_on_load)
    return cls()
