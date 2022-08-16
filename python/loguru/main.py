import sys

from loguru import logger

logger.debug("This is a ready-to-use debug log")
# Use add() to setup handler, formatter, filter
logger.add(sys.stderr, format="{time} {level} {message}",
           filter="my_module", level="INFO")
# File logging with rotation/retention/compression
logger.add("/tmp/file_{time}.log")
logger.add("/tmp/file_1.log", rotation="500 MB")    # Automatically rotate too big file
# logger.add("/tmp/file_2.log", rotation="12:00")     # New file is created each day at noon
# logger.add("/tmp/file_3.log", rotation="1 week")    # Once the file is too old, it's rotated
# logger.add("/tmp/file_X.log", retention="10 days")  # Cleanup after some time
# logger.add("/tmp/file_Y.log", compression="zip")    # Save some loved space
logger.info("If you're using Python {}, prefer {feature} of course!", 3.6,
            feature="f-strings")
# Enable backtrace
logger.add("out.log", backtrace=True, diagnose=True)  # Caution, may leak sensitive data in prod

def func(a, b):
    return a / b

def nested(c):
    try:
        func(5, c)
    except ZeroDivisionError:
        logger.exception("What?!")

nested(0)
# Check more here: https://github.com/Delgan/loguru
