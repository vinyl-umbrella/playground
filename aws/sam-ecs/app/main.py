import logging
import time

from mylib import sample
from utils.log_config import JsonlFormatter


def main():
    l = logging.getLogger(__name__)
    l.info("Starting ECS app")
    l.info("hello", extra={"key1": "value2", "key2": 1234})
    l.warning(
        "This is a warning message", extra={"key1": {"nested_key": "nested_value", "key2": 5678}}
    )
    time.sleep(2)
    l.info("Waiting for 2 \nseconds...")
    sample.echo_msg("Hello from ECS app!")

    try:
        1 / 0
    except ZeroDivisionError as e:
        l.error(f"An error occurred: {e}", exc_info=True)
    l.info("ECS app finished")


if __name__ == "__main__":
    jsonl_formatter = JsonlFormatter()

    logging.basicConfig(
        level=logging.INFO,
        format="%(message)s",  # NOTE: 実際のフォーマットはカスタムフォーマッタが処理
        handlers=[logging.StreamHandler()],
    )

    for handler in logging.getLogger().handlers:
        handler.setFormatter(jsonl_formatter)

    main()
