# -*- coding: utf-8 -*-

from utils import es_utils
from utils import milvus_utils
from typing import List, Dict, Any


# ---------------------- 1. 问答库生命周期 ----------------------
def init_qa_base(user_id: str, qa_base: str, qa_id: str, embedding_model_id: str) -> Dict[str, Any]:
    """
    创建问答库
    """
    if not user_id or not qa_base or not qa_id:
        return {"code": 1, "message": "缺失必填参数"}
    # 伪逻辑：检查是否已存在
    # if exist(user_id, qa_id):
    #     return {"code": 1, "message": "qa_id 已存在"}
    return {"code": 0, "message": "success"}


def delete_qa_base(user_id: str, qa_base: str, qa_id: str) -> Dict[str, Any]:
    """
    删除问答库
    """
    if not user_id or not qa_id:
        return {"code": 1, "message": "缺失必填参数"}
    # 伪逻辑：检查是否存在
    # if not exist(user_id, qa_id):
    #     return {"code": 1, "message": "qa_id 不存在"}
    return {"code": 0, "message": "success"}


# ---------------------- 2. 问答对 CRUD ----------------------

def batch_add_qas(user_id: str, qa_base: str, qa_id: str, qa_pairs: List[Dict[str, str]]) -> Dict[str, Any]:
    """
    批量新增问答对
    qa_pairs: [{"qa_pair_id":"123","question":"q","answer":"a"}, ...]
    """
    if not qa_pairs:
        return {"code": 1, "message": "QAPairs 为空"}

    return {"code": 0, "message": "success"}


def get_qa_list(user_id: str, qa_base: str, qa_id: str, page_size: int, search_after: int) -> Dict[str, Any]:
    """
    分页获取问答对（冗余列表）
    """
    # 伪数据
    mock_list = []
    return {"code": 0, "message": "success", "data": {"qa_list": mock_list, "qa_pair_total_num": 3}}


def update_qa(user_id: str, qa_base: str, qa_id: str, qa_pairs: List[Dict[str, str]]) -> Dict[str, Any]:
    """
    批量更新问答对（全量覆盖）
    """
    if not qa_pairs:
        return {"code": 1, "message": "QAPairs 为空"}
    return {"code": 0, "message": "success"}


def batch_delete_qas(user_id: str, qa_base: str, qa_id: str, qa_pair_ids: List[str]) -> Dict[str, Any]:
    """
    批量删除问答对
    """
    if not qa_pair_ids:
        return {"code": 1, "message": "QAPairIds 为空"}
    return {"code": 0, "message": "success"}


def update_qa_status(user_id: str, qa_base: str, qa_id: str, qa_pair_id: str, status: bool) -> Dict[str, Any]:
    """
    启停单个问答对
    """
    return {"code": 0, "message": "success"}


# ---------------------- 3. 元数据管理 ----------------------

def update_qa_metas(user_id: str, qa_base: str, qa_id: str, metas: List[Dict[str, Any]]) -> Dict[str, Any]:
    """
    全量覆盖更新元数据
    metas: [{"qa_pair_id":"123","metadata_list":[...]}, ...]
    """
    if not metas:
        return {"code": 1, "message": "metas 为空"}
    return {"code": 0, "message": "success"}


def delete_meta_by_keys(user_id: str, qa_base: str, qa_id: str, keys: List[str]) -> Dict[str, Any]:
    """
    批量删除指定 key 的元数据（跨所有问答对）
    """
    if not keys:
        return {"code": 1, "message": "keys 为空"}
    return {"code": 0, "message": "success"}


def rename_meta_keys(user_id: str, qa_base: str, qa_id: str, mappings: List[Dict[str, str]]) -> Dict[str, Any]:
    """
    批量重命名元数据 key
    mappings: [{"old_key":"company","new_key":"organization"}, ...]
    """
    if not mappings:
        return {"code": 1, "message": "mappings 为空"}
    return {"code": 0, "message": "success"}
