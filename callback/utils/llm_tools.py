import logging
from pathlib import Path

from langchain_core.prompts import PromptTemplate  # type: ignore

from utils.log import logger

prompts_base_path = "static/prompts"


def load_prompt_template(
    template_file: Path,
    template_file_base: str = prompts_base_path,
    encoding: str = "utf-8",
) -> PromptTemplate:
    """Load and validate a prompt template from file.

    Args:
        template_path: Path to the template file
        encoding: File encoding (default: utf-8)

    Returns:
        Loaded PromptTemplate instance

    Raises:
        FileNotFoundError: If template file doesn't exist
        ValueError: If template format is invalid
    """
    try:
        prompt_template_path = Path(template_file_base, template_file)
        return PromptTemplate.from_file(prompt_template_path, encoding=encoding)
    except FileNotFoundError as e:
        logger.error(f"Template file not found: {prompt_template_path}")
        raise
    except ValueError as e:
        logger.error(f"Invalid template format in {prompt_template_path}: {e}")
        raise
    except Exception as e:
        logger.error(f"Unexpected error loading template: {e}")
        raise


def format_prompt_template(template_name, **kwargs):
    """
    格式化提示词模板

    参数:
    template_name (str): 模板文件名称
    kwargs: 模板所需的参数字典

    返回:
    str: 格式化后的提示词
    """
    prompt_template = load_prompt_template(template_name)
    return prompt_template.format(**kwargs)
