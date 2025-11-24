import hashlib

def generate_md5(content_str):
    # 创建一个md5 hash对象
    md5_obj = hashlib.md5()

    # 对字符串进行编码，因为md5需要bytes类型的数据
    md5_obj.update(content_str.encode('utf-8'))

    # 获取十六进制的MD5值
    md5_value = md5_obj.hexdigest()

    return md5_value

