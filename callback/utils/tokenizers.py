import logging
import time
from typing import List, Union

from utils.log import logger


class CustomTokenizer:
    # 默认字符到token的转换比例 (经验值)
    # 中文环境下通常 1个token ≈ 1.5 ~ 2.0 个字符
    _default_char_to_token_ratio = 1.7

    def __init__(self, char_to_token_ratio: float = None):
        """
        初始化仅用于估算的Tokenizer

        Args:
            model_name: 仅作为标识符
            char_to_token_ratio: 自定义转换比例，如不传则使用默认值 1.7
        """
        self.ratio = char_to_token_ratio or self._default_char_to_token_ratio
        logger.info(f"Initialized estimator  with ratio {self.ratio}")

    def count_tokens(self, text: Union[str, List[str]]) -> int:
        """
        使用字符长度估算token数量
        公式: token数量 = 字符数量 / 转换比例
        """
        # start_time = time.time()
        result = 0

        if isinstance(text, str):
            result = int(len(text) / self.ratio)
        elif isinstance(text, list):
            result = sum(int(len(t) / self.ratio) for t in text)
        else:
            raise ValueError("输入必须是字符串或字符串列表")

        # logger.debug(f"计算耗时: {time.time() - start_time:.6f}秒")
        return result

    def truncate_text(self, text: str, max_tokens: int) -> str:
        """
        基于估算按最大token数量截断文本
        """
        if max_tokens < 1:
            raise ValueError("max_tokens必须大于0")

        # 反向计算最大允许的字符数
        max_chars = int(max_tokens * self.ratio)

        if len(text) <= max_chars:
            return text

        return text[:max_chars]
