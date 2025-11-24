import os
import logging
import datetime
from logging.handlers import TimedRotatingFileHandler
from logging.handlers import RotatingFileHandler
import sys

###转化为北京时间
current_date = datetime.datetime.now()
current_time = current_date.strftime('%Y%m%d')

# 全局变量定义
LOG_DIRECTORY = f'./logs'
LOG_LEVEL = logging.INFO
INTERVAL = 1 
BACKUP_COUNT = 10  # 保留的日志文件数量


def get_log_directory():
    """动态获取日志目录路径"""
    # 判断是否为打包环境
    if getattr(sys, 'frozen', False):
        # 打包环境：使用可执行文件所在目录
        base_dir = os.path.dirname(sys.executable)
    else:
        # 开发环境：使用当前文件所在目录
        base_dir = os.path.dirname(os.path.abspath(__file__))

    # 在基础目录下创建 logs 子目录
    log_dir = os.path.join(base_dir, 'logs')
    os.makedirs(log_dir, exist_ok=True)
    return log_dir

def setup_logging(app_name,logger_name):
    """
    初始化日志配置。

    参数:
    app_name (str): 应用名称，用于日志文件命名
    """
    # 使用动态路径获取日志目录
    LOG_DIRECTORY = get_log_directory()

    # 定义日志文件的完整路径，日志文件命名规则 {app_name}.log
    log_file_path = os.path.join(LOG_DIRECTORY, f'{app_name}.log')

    # 创建logger
    logger = logging.getLogger(logger_name)
    logger.setLevel(LOG_LEVEL)

    # 确保日志目录存在
    os.makedirs(LOG_DIRECTORY, exist_ok=True)

    # 创建一个handler，用于写入日志文件  
    # file_handler = TimedRotatingFileHandler(log_file_path, when='D', interval=INTERVAL, backupCount=BACKUP_COUNT, encoding='utf-8')
    file_handler = RotatingFileHandler(log_file_path, maxBytes=1024*1024*5, backupCount=5, encoding='utf-8') 
    file_handler.setLevel(logging.INFO)
  
    # 再创建一个handler，用于输出到控制台  
    console_handler = logging.StreamHandler()  
    console_handler.setLevel(logging.INFO)
    
    # 定义handler的输出格式  
    formatter = logging.Formatter('%(asctime)s - %(filename)s:%(funcName)s:%(lineno)d - %(levelname)s - %(message)s',  datefmt='%Y-%m-%d %H:%M:%S')   
    
    file_handler.setFormatter(formatter)  
    console_handler.setFormatter(formatter) 

    # 清除已存在的处理器，防止重复添加
    if logger.hasHandlers():
        logger.handlers.clear()
        
    # 给logger添加handler  
    logger.addHandler(file_handler)  

    return logger
