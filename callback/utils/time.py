from datetime import datetime, timedelta, timezone


def get_current_cst_time(time_format="%Y-%m-%d %H:%M:%S"):
    """
    获取当前中国标准时间（CST）
    :param time_format: 时间格式字符串，默认为"%Y-%m-%d %H:%M:%S"
    :return: 格式化后的当前CST时间字符串
    """
    cst_tz = timezone(timedelta(hours=8))
    current_cst_time = datetime.now(cst_tz)
    return current_cst_time.strftime(time_format)
