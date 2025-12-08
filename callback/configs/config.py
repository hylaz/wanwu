import configparser
import os
from configparser import ConfigParser

from flask import Flask


class Config:
    callback_cfg = None


config = Config()  # 全局单例对象


def load_config():
    path = "configs/config.ini"
    ini_config = configparser.ConfigParser()
    ini_config.read(path, encoding="utf-8")
    update_config_with_env(ini_config)

    config.callback_cfg = ini_config


def flatten_config_to_app(app: Flask, config: ConfigParser):
    """
    将 ConfigParser 配置平铺到 Flask app.config 中
    例如：
        [redis]
        host=127.0.0.1
    会被转为 app.config["REDIS_HOST"] = "127.0.0.1"
    """
    for section in config.sections():
        for key, value in config.items(section):
            config_key = f"{section.upper()}_{key.upper()}"
            app.config[config_key] = value


def update_config_with_env(config: ConfigParser):
    # 1. 对 ini 中的键做环境变量覆盖
    for section in config.sections():
        for key in config[section]:
            env_key = f"{section}_{key}".upper().replace(".", "_").replace("-", "_")
            if env_key in os.environ:
                # 环境变量有 → 覆盖
                config[section][key] = os.environ[env_key]
