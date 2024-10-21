import logging
import time

logger = logging.getLogger()


def main():
    logger.info("Starting ECS app")
    time.sleep(10)
    logger.info("ECS app finished")


if __name__ == "__main__":
    main()
