## 项目目录结构
```text
callback/
├── callback/
│   ├── models/                      # 数据模型层，如 ORM 模型
│   ├── routes/                      # 路由层，定义 API 接口
│   ├── services/                    # 业务逻辑层，处理具体业务逻辑
│   ├── static/                      # 静态文件（如 js、css、图片等）
│   └── __init__.py                  # 包初始化，注册蓝图、初始化应用
│
├── configs/
│   ├── config.ini                   # 配置文件，存储各种环境和服务配置
│   └── config.py                    # 配置加载脚本，将 config.ini注册到Flask App中
├── extensions/                      # 初始化第三方服务
│   ├── minio.py                     # MinIO 客户端初始化
│   └── redis.py                     # Redis 客户端初始化
│
├── utils/
│   ├── response.py                  # 响应封装工具，例如统一返回格式
│   └── router_trace.py              # 路由追踪工具，可能用于日志或调试
├── static/                          # 全局静态文件
├── .gitignore                       # Git 忽略文件配置
├── README.md                        # 项目说明文档
├── requirements.txt                 # Python 依赖包清单
└── run.py                           # 项目入口文件，启动 Flask 服务
```

## 运行callback服务
```bash
# 启动容器
make -f Makefile.develop run-callback
# 关闭容器
make -f Makefile.develop stop-callback
```

## 构建镜像
```bash
# 构建base镜像
make docker-image-callback-base
# 构建发版镜像
make docker-image-callback
```

## Swagger文档
运行容器之后访问
```bash
http://localhost:8669/apidocs
```

## 配置文件
配置文件放在`configs/config.ini`文件中，docker compose中的环境变量必须在`config.ini`中定义。
比如docker compose中的redis配置
```yaml
REDIS_HOST: ${WANWU_REDIS_HOST}
REDIS_PORT: ${WANWU_REDIS_PORT}
REDIS_PASSWORD: ${WANWU_REDIS_PASSWORD}
```
在`config.ini`中定义为：
```ini
[REDIS]
HOST=
PORT=
PASSWORD=
```
docker compose环境变量的优先级更高，会覆盖config.ini中的配置。如果需要私有化配置，则只需在`config.ini`中定义。

项目配置通过全局单例对象读取，例如在`extensions/redis.py`中读取配置。
```python
from configs.config import config
redis_client = redis.Redis(
        host=config.callback_cfg["REDIS"]["HOST"],
        port=config.callback_cfg["REDIS"]["PORT"],
        password=config.callback_cfg["REDIS"]["PASSWORD"],
        decode_responses=True 
    )
```

## 代码规范
[![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/psf/black)
[![Imports: isort](https://img.shields.io/badge/%20imports-isort-%231674b1?style=flat&labelColor=ef8336)](https://pycqa.github.io/isort/)