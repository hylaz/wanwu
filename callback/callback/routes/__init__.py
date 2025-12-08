import importlib
import pkgutil

from flask import Blueprint

callback_bp = Blueprint("callback", __name__)

# 加载所有 blueprint 路由模块
# __path__ 是当前包所在的路径列表
# __name__ 是当前包的名称（用于相对导入）
for loader, module_name, is_pkg in pkgutil.iter_modules(__path__):
    # 排除不需要自动导入的文件（例如 __init__ ）
    if module_name == "__init__":
        continue

    # 动态导入模块: from . import module_name
    importlib.import_module(f".{module_name}", package=__name__)

# 现在，所有文件夹下的 .py 文件都会被自动执行，路由也会自动注册
