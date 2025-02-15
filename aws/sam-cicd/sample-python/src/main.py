import logging


def lambda_handler(event, context):
    logger = logging.getLogger(__name__)
    logger.info("Hello World")
    logger.warning("This is a warning")
    logger.error("This is an error")


if __name__ == "__main__":
    lambda_handler({}, None)
