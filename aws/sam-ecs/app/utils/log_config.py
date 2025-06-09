import json
import logging


class JsonlFormatter(logging.Formatter):
    """JSONL で出力するカスタムフォーマッタ"""

    def format(self, record: logging.LogRecord) -> str:
        log_data: dict[str, object] = {
            "timestamp": self.formatTime(record, "%Y-%m-%d %H:%M:%S"),
            "level": record.levelname,
            "func": f"{record.filename}.{record.funcName}:{record.lineno}",
            "message": self.trim_message(record.getMessage()),
        }

        if record.exc_info:
            import traceback

            exc_type, exc_value, exc_tb = record.exc_info
            if exc_type is not None:
                log_data["exception"] = {
                    "type": exc_type.__name__ if exc_type else None,
                    "message": str(exc_value),
                    "traceback": traceback.format_exception(exc_type, exc_value, exc_tb),
                }

        for key, value in record.__dict__.items():
            if key not in {
                "args",
                "asctime",
                "created",
                "exc_info",
                "exc_text",
                "filename",
                "funcName",
                "id",
                "levelname",
                "levelno",
                "lineno",
                "module",
                "msecs",
                "message",
                "msg",
                "name",
                "pathname",
                "process",
                "processName",
                "relativeCreated",
                "stack_info",
                "thread",
                "threadName",
                "taskName",
            }:
                log_data[key] = value

        return json.dumps(log_data, cls=JsonEncoder)

    def trim_message(self, message: object) -> str:
        if message and isinstance(message, str):
            return message.replace("\n", " ").replace("\r", " ")
        return str(message).strip()


class JsonEncoder(json.JSONEncoder):
    def default(self, o: object) -> object:
        try:
            return super().default(o)
        except TypeError:
            return str(o)
