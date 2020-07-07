from teststevedore import loader

if __name__ == "__main__":
    drivers = loader.get_available_drivers()
    print(drivers)
    for d in drivers:
        dcls = loader.get_driver_instance(d)
        dcls.do()
