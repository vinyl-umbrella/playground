import logging


def echo_msg(msg: str):
    l = logging.getLogger()
    l.info(f"Echoing message: {msg}")
