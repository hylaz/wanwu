import settings
import re
import hashlib
from enum import Enum

class IndexType(Enum):
    """索引类型枚举"""
    MAIN = "main"
    SNIPPET = "snippet"
    CONTENT_CONTROL = "content_control"
    FILE_CONTROL = "file_control"

#获取主索引名
def get_main_index_name(user_id:str) -> str:
    return settings.INDEX_NAME_PREFIX + user_id

def get_snippet_index_name(user_id:str) -> str:
    return settings.SNIPPET_INDEX_NAME_PREFIX + user_id

def get_content_control_index_name(user_id:str) -> str:
    return 'content_control_' + get_main_index_name(user_id)

def get_file_index_name(user_id:str) -> str:
    return 'file_control_' + get_main_index_name(user_id)

def get_qa_index_name(user_id:str) -> str:
    return 'qa_' + settings.INDEX_NAME_PREFIX + user_id


def validate_index_name(index_name):
    # Check length
    if len(index_name) > 255:
        return False, "Index name cannot exceed 255 characters"

    # Check for illegal characters
    # if not re.match(r'^[a-z0-9_-]+$', index_name):
    if not re.match(r'^[a-z0-9_\u4e00-\u9fa5-]+$', index_name, re.UNICODE):
        return False, "Index name can only contain lowercase letters, numbers, hyphens, and underscores"

    # Check if starts with a hyphen or underscore
    if index_name.startswith('-') or index_name.startswith('_'):
        return False, "Index name cannot start with a hyphen or underscore"

    # Check for commas
    if ',' in index_name:
        return False, "Index name cannot contain commas"

    # Check if name is "." or ".."
    if index_name in ['.', '..']:
        return False, "Index name cannot be \".\" or \"..\""

    # Check if starts with "." or ".."
    if index_name.startswith('.') or index_name.startswith('..'):
        return False, "Index name cannot start with \".\" or \"..\""

    return True, "Index name is valid"


def generate_md5(content_str):
    # 创建一个md5 hash对象
    md5_obj = hashlib.md5()

    # 对字符串进行编码，因为md5需要bytes类型的数据
    md5_obj.update(content_str.encode('utf-8'))

    # 获取十六进制的MD5值
    md5_value = md5_obj.hexdigest()

    return md5_value