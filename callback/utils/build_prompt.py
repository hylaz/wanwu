import logging
from pathlib import Path
from re import sub

from configs.config import config
from utils.llm_tools import format_prompt_template, prompts_base_path
from utils.log import logger
from utils.tokenizers import CustomTokenizer

# 常量定义
MAX_INPUT_TOKENS = int(
    config.callback_cfg["MODEL"]["MAX_INPUT_TOKENS"]
)  # 从配置读取的最大模型输入 token 数

# 拼接文本时的截断阈值（单位：token）
# 搜索内容拼接时如果剩余token数量低于该值，则舍弃该chunk，否则将chunk截断后再拼接
TRUNCATION_THRESHOLD = 20


def assemble_search_context(
    search_list: list,
    citation_label: str,
    tokenizer: CustomTokenizer,
    max_token_length: int,
) -> str:
    """
    根据搜索结果构建受限长度的上下文内容字符串，拼接内容不超过指定 token 上限。

    参数：
        search_list (list): 搜索结果组成的列表，每项为字典，包含以下字段：
            - file_name (str): 文件名或标题。
            - sub_query (str): 子查询内容。
            - snippet (str): 匹配片段。
        citation_label (str): 每段拼接内容前缀标签，例如 "citation" 或 "参考信息"。
        tokenizer (CustomTokenizer): 用于统计和截断 token 的 tokenizer 实例。
        max_token_length (int): 限定的搜索内容的最大 token 数，超过后停止拼接。

    返回：
        str: 构建完成的上下文字符串，用于填充提示词模板。
    """
    search_context = ""
    current_tokens = 0

    for idx, item in enumerate(search_list):
        file_name = (
            f"参考文件：{item.get('file_name', '')}\n" if item.get("file_name") else ""
        )
        sub_query = f"{item.get('sub_query', '')}\n" if item.get("sub_query") else ""
        snippet = item.get("snippet", "")

        segment = (
            f"{citation_label}-{idx + 1:02d}:\n{sub_query}{file_name}{snippet}"
            if citation_label
            else f"{idx + 1:02d}:\n{sub_query}{file_name}{snippet}"
        )

        segment_tokens = tokenizer.count_tokens(segment)
        new_total = current_tokens + segment_tokens

        if new_total > max_token_length:
            leftover = max_token_length - current_tokens
            if leftover > TRUNCATION_THRESHOLD:
                truncated = tokenizer.truncate_text(segment, leftover)
                search_context += "\n\n" + truncated
                logger.warn(
                    f"超出token限制，终止拼接，使用chunk数量: {idx+1}（最后一项已截断）"
                )
            else:
                logger.warn(f"超出token限制，终止拼接，使用chunk数量: {idx}")
            break

        search_context += "\n\n" + segment
        current_tokens = new_total

    return search_context


def build_prompt_from_search_list(
    query: str,
    search_list: list,
    template_prefix: str,
    auto_citation: bool = False,
    max_length: int = MAX_INPUT_TOKENS,
    **kwargs,
) -> str:
    """
    构建完整提示词（Prompt），将模板与搜索上下文拼接后填充参数。

    参数：
        query (str): 用户原始提问内容。
        search_list (list): 搜索结果组成的列表，见 build_limited_search_contents 中字段要求。
        template_prefix (str): 所用模板的前缀（例如 "bing"、"docqa"、"rag"）。
        auto_citation (bool): 是否启用 citation 样式引用（影响模板选择及编号格式）。
        model (str): 当前使用的语言模型名（用于 tokenizer 选择）。
        max_length (int): 总 token 限制上限，包括模板 + 搜索内容。
        **kwargs: 其他模板所需参数，例如 cur_date 等。

    返回：
        str: 构造好的提示词，可直接送入 LLM 使用。
    """
    tokenizer = CustomTokenizer()
    citation_label = "citation" if auto_citation else "参考信息"

    template_name = (
        f"{template_prefix}_prompt{'_citation' if auto_citation else ''}.txt"
    )
    template_path = Path(
        prompts_base_path, template_name
    )  # 确保 prompts_base_path 为 Path 类型
    template_content = template_path.read_text(encoding="utf-8")
    template_token_num = tokenizer.count_tokens(template_content)

    search_content_max_length = max_length - template_token_num
    search_contents = assemble_search_context(
        search_list, citation_label, tokenizer, search_content_max_length
    )

    prompt = format_prompt_template(
        template_name, question=query, context=search_contents, **kwargs
    )
    return prompt


def build_docqa_prompt_from_search_list(
    query: str,
    search_list: list,
    auto_citation: bool = False,
    max_length: int = MAX_INPUT_TOKENS,
) -> str:
    """
    构建适用于 DocQA（基于文档问答）场景的提示词。
    与 Bing 相比不包含日期。
    """
    return build_prompt_from_search_list(
        query=query,
        search_list=search_list,
        template_prefix="docqa",
        auto_citation=auto_citation,
        max_length=max_length,
    )
