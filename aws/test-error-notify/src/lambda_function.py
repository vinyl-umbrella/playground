import logging


def lambda_handler(event, context):
    l = logging.getLogger(__name__)
    l.info(event)
    l.info("hello world")

    raise Exception("This is an exception")
