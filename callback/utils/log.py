import logging
import os
from logging.handlers import RotatingFileHandler


def init_logger(
    log_name="callback_logger", log_file="log/app.log", log_level=logging.DEBUG
):
    formatter = logging.Formatter(
        "[%(asctime)s][%(levelname)s][%(module)s] %(message)s [%(pathname)s:%(lineno)d|%(funcName)s]"
    )

    # 确保日志目录存在
    os.makedirs(os.path.dirname(log_file), exist_ok=True)

    # rotating file
    rotate_handler = RotatingFileHandler(
        log_file, maxBytes=1024 * 1024 * 10, backupCount=10
    )
    rotate_handler.setFormatter(formatter)

    # console
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)

    logger = logging.getLogger(log_name)
    logger.setLevel(log_level)
    logger.addHandler(rotate_handler)
    logger.addHandler(console_handler)

    # 防止日志重复（如果多次调用 init_logger）
    logger.propagate = False

    return logger


logger = init_logger()
