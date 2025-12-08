import redis

from configs.config import config
from utils.log import logger

redis_client: redis.Redis = None


def init_redis():
    global redis_client

    redis_client = redis.Redis(
        host=config.callback_cfg["REDIS"]["HOST"],
        port=config.callback_cfg["REDIS"]["PORT"],
        password=config.callback_cfg["REDIS"]["PASSWORD"],
        decode_responses=True,
    )

    # 测试连接（可选）
    try:
        redis_client.ping()
        logger.info("Redis connected.")
    except redis.ConnectionError:
        logger.info("Redis connection failed!")

    return redis_client
