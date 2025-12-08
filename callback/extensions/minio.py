import io
import json
import os
from re import S
from typing import Generator, Optional

from minio import Minio
from minio.commonconfig import ENABLED, Filter
from minio.error import S3Error
from minio.lifecycleconfig import Expiration, LifecycleConfig, Rule

from configs.config import config
from utils.log import logger


class Client:
    def __init__(
        self,
    ):
        # 1. 获取原始配置值
        secure_config = config.callback_cfg["MINIO"]["SECURE"]

        # 2. 安全地转换为布尔值
        # 如果 secure_config 是字符串 "False"、"false" 或 "0"，这个逻辑能确保它变为 False
        if isinstance(secure_config, str):
            # 只要字符串是 'true' (不区分大小写) 才为真，否则为假
            is_secure = secure_config.lower() == "true"
        else:
            # 如果本身就是布尔类型或整数，直接转换
            is_secure = bool(secure_config)
        self.cli = Minio(
            endpoint=config.callback_cfg["MINIO"]["ENDPOINT"],
            access_key=config.callback_cfg["MINIO"]["USER"],
            secret_key=config.callback_cfg["MINIO"]["PASSWORD"],
            secure=is_secure,
        )

    def get_cli(self) -> Minio:

        return self.cli

    def put_object(self, bucket_name: str, object_name: str, data: bytes) -> object:
        """
        接收 bytes，内部封装为 io.BytesIO
        """
        # Python 的 put_object 需要类似文件的对象 (read 方法) 和 长度
        data_stream = io.BytesIO(data)
        length = len(data)

        return self.cli.put_object(
            bucket_name=bucket_name,
            object_name=object_name,
            data=data_stream,
            length=length,
            content_type="application/octet-stream",  # 默认类型，可视情况修改
        )

    def put_object_from_stream(
        self, bucket_name: str, object_name: str, file_data
    ) -> object:
        """
        新方法：支持 file_buffer (BytesIO 或 打开的文件对象)
        会自动处理 seek(0) 和 计算长度
        """
        try:
            # 1. 自动重置指针到开头 (防止“磁带”跑完)
            if hasattr(file_data, "seek"):
                file_data.seek(0)

            # 2. 自动获取长度
            length = 0
            # 情况 A: 如果是 BytesIO，直接用 getbuffer
            if hasattr(file_data, "getbuffer"):
                length = file_data.getbuffer().nbytes
            # 情况 B: 如果是普通文件对象，用 stat
            elif hasattr(file_data, "fileno"):
                length = os.fstat(file_data.fileno()).st_size
            # 情况 C: 通用兜底 (移到最后看位置，再移回来)
            else:
                file_data.seek(0, 2)  # 移到末尾
                length = file_data.tell()
                file_data.seek(0)  # 移回开头

            # 3. 执行上传
            return self.cli.put_object(
                bucket_name=bucket_name,
                object_name=object_name,
                data=file_data,
                length=length,
                content_type="application/octet-stream",
            )
        except Exception as e:
            logger.info(f"Stream upload failed: {e}")
            return None  # 或者 raise e

    def get_object(self, bucket_name: str, object_name: str) -> bytes:
        """
        返回 bytes
        """
        response = None
        try:
            response = self.cli.get_object(bucket_name, object_name)
            return response.read()
        finally:
            if response:
                response.close()
                response.release_conn()

    def delete_object(self, bucket_name: str, object_name: str):
        self.cli.remove_object(bucket_name, object_name)

    def list_objects(self, bucket_name: str, prefix: str) -> Generator:
        return self.cli.list_objects(bucket_name, prefix=prefix, recursive=True)

    def create_bucket_if_not_exist(self, bucket_name: str):
        try:
            found = self.cli.bucket_exists(bucket_name)
            if not found:
                self.cli.make_bucket(bucket_name)
                logger.info(f"Bucket {bucket_name} created.")
        except S3Error as err:
            # 如果错误码是 BucketAlreadyOwnedByYou，说明被其他进程抢先创建了，可以忽略
            if err.code == "BucketAlreadyOwnedByYou":
                logger.info(
                    f"Bucket {bucket_name} already owned by you (concurrency handled)."
                )
            else:
                # 其他真正的错误则抛出
                raise err

    def create_public_bucket_if_not_exist(self, bucket_name: str):
        found = self.cli.bucket_exists(bucket_name)
        if not found:
            self.cli.make_bucket(bucket_name)
            self.set_bucket_public(bucket_name)

    def set_path_expire_by_day(self, bucket_name: str, path: str, days: int):
        """
        配置生命周期规则
        """
        config = LifecycleConfig(
            [
                Rule(
                    status=ENABLED,
                    rule_id="delete-after-days",
                    rule_filter=Filter(prefix=path),
                    expiration=Expiration(days=days),
                ),
            ]
        )
        self.cli.set_bucket_lifecycle(bucket_name, config)

    def set_bucket_public(self, bucket_name: str):
        """
        设置桶策略为公开只读
        """
        policy = {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": ["s3:GetObject"],
                    "Resource": [f"arn:aws:s3:::{bucket_name}/*"],
                }
            ],
        }
        # 将字典转换为 JSON 字符串
        policy_json = json.dumps(policy)
        self.cli.set_bucket_policy(bucket_name, policy_json)


# 定义全局变量
minio_client: Optional["Client"] = None


def init_minio() -> "Client":
    """
    初始化并返回全局单例 MinIO 客户端。
    如果已经初始化过，则直接返回现有实例，忽略新的 cfg。
    """
    global minio_client

    if minio_client is None:
        # 实例化上一轮定义的 Client 类
        minio_client = Client()
        logger.info(f"MinIO Client initialized with endpoint")
        bucket_name = config.callback_cfg["MINIO"]["BUCKET_NAME"]
        minio_client.create_bucket_if_not_exist(bucket_name)
        minio_client.set_bucket_public(bucket_name)
