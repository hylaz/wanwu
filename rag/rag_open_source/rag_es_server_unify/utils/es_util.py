# -*- coding: utf-8 -*-

import requests
import uuid
import re
import json
import time
import hashlib
import warnings
from openai import OpenAI
from model.model_manager import get_model_configure
import random
import numpy as np
import utils.mapping_util as es_mapping
from settings import GET_KB_ID_URL
from settings import KBNAME_MAPPING_INDEX, DELETE_BACTH_SIZE
from log.logger import logger
from utils.config_util import es
from elasticsearch import helpers

warnings.filterwarnings("ignore")

def get_maas_kb_id(user_id, kb_name):
    """获取maas的kb_id"""
    try:
        url = GET_KB_ID_URL + f"?userId={user_id}&categoryName={kb_name}"
        r = requests.get(url)
        result_data = json.loads(r.text)
        if result_data["code"] == 0:
            kb_id = result_data["data"].get('categoryId')
            return kb_id
        else:
            raise RuntimeError(f"{kb_name},get_maas_kb_id Error, result: {result_data}, url:{GET_KB_ID_URL}")
    except Exception as e:
        raise RuntimeError(kb_name + ",get_maas_kb_id Error: " + str(e) + "url:" + GET_KB_ID_URL) from e


def check_status():
    # 先等待随机0-10秒，模拟网络延迟
    time.sleep(random.uniform(0, 10))
    """ 检查 uk 映射表是否支持改名，如果不支持，修改结构"""
    uk_mappings = {
        "properties": {
            "index_name": {"type": "keyword"},  # 指定为 keyword，方便用于排序和聚合
            "userId": {"type": "keyword"},  # 指定为 keyword，方便用于排序和聚合
            "kb_name": {"type": "keyword"},  # 指定为 keyword，方便用于排序和聚合
            "kb_id": {"type": "keyword"},  # 指定为 keyword，方便用于排序和聚合
        }
    }
    try:
        create_index_if_not_exists(KBNAME_MAPPING_INDEX, mappings=uk_mappings)  # 确保 KBNAME_MAPPING_INDEX 已创建
        # 检查 KBNAME_MAPPING_INDEX 的 mapping
        index = es.indices.get(index=KBNAME_MAPPING_INDEX)
        if "kb_id" in index[KBNAME_MAPPING_INDEX]['mappings']['properties']:
            logger.info(f"KBNAME_MAPPING_INDEX mapping is new")
            uk_docs = fetch_all_documents(KBNAME_MAPPING_INDEX)
            actions = []
            for item in uk_docs:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
                doc_id = item["_id"]
                try:
                    kb_id = get_maas_kb_id(item["_source"]["userId"], item["_source"]["kb_name"])
                    logger.info(f"{item['_source']['userId']},{item['_source']['kb_name']},kb_id:{kb_id}")
                except Exception as e:
                    logger.error(f"{item['_source']['userId']},{item['_source']['kb_name']},get_maas_kb_id Error:{e}")
                    continue
                action = {
                    "_op_type": "update",
                    "_index": KBNAME_MAPPING_INDEX,
                    "_id": doc_id,
                    "doc": {"kb_id": kb_id}
                }
                actions.append(action)
            # 执行更新操作,并返回
            helpers.bulk(es, actions)
            es.indices.refresh(index=KBNAME_MAPPING_INDEX)
            logger.info(f"KBNAME_MAPPING_INDEX mapping updated")
        else:  # 如果不存在，则更新 mapping 及导入数据
            add_field_to_mapping(KBNAME_MAPPING_INDEX, "kb_id", "keyword")
            # 刷新导入kb_id数据，离石的和kb_name 一致
            uk_docs = fetch_all_documents(KBNAME_MAPPING_INDEX)
            actions = []
            for item in uk_docs:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
                doc_id = item["_id"]
                try:
                    kb_id = get_maas_kb_id(item["_source"]["userId"], item["_source"]["kb_name"])
                    logger.info(f"{item['_source']['userId']},{item['_source']['kb_name']},kb_id:{kb_id}")
                except Exception as e:
                    logger.error(f"{item['_source']['userId']},{item['_source']['kb_name']},get_maas_kb_id Error:{e}")
                    continue
                action = {
                    "_op_type": "update",
                    "_index": KBNAME_MAPPING_INDEX,
                    "_id": doc_id,
                    "doc": {"kb_id": kb_id}
                }
                actions.append(action)
            # 执行更新操作,并返回
            helpers.bulk(es, actions)
            es.indices.refresh(index=KBNAME_MAPPING_INDEX)
            logger.info(f"KBNAME_MAPPING_INDEX mapping updated")
    except Exception as e:
        logger.error(f"Error checking status: {e}")


def get_embs(texts: list, embedding_model_id=""):
    """ 先使用 openai embedding协议获取 文本向量"""
    emb_info = get_model_configure(embedding_model_id)
    logger.info(f"Starting embedding request for {len(texts)} texts, model: {emb_info.model_name}")

    api_key = emb_info.api_key or "fake api key"
    # 安全记录API Key（仅显示部分）
    masked_key = api_key[:4] + "****" + api_key[-4:] if len(api_key) > 8 else "****"

    client = OpenAI(
        api_key=api_key,
        base_url=emb_info.endpoint_url,
    )

    # 安全的请求日志
    request_details = {
        "url": emb_info.endpoint_url,
        "model": emb_info.model_name,
        "api_key": masked_key,  # 使用脱敏后的key
        "text_count": len(texts),
        "input": texts
    }
    logger.info(f"Sending embedding request: {json.dumps(request_details, ensure_ascii=False)}")

    # 退避间隔
    rate_limit_backoff  = [10, 20, 40, 60] # 限流退避
    other_error_max_retries = 2        # 其他错误最多重试2次
    other_error_wait = 0.5             # 每次0.5s

    attempt = 0
    while attempt < max(len(rate_limit_backoff), other_error_max_retries) + 1:
        try:
            # 记录请求开始时间
            start_time = time.time()
            completion = client.embeddings.create(
                model=emb_info.model_name,
                input=texts,
                encoding_format="float"
            )

            response_json = json.loads(completion.model_dump_json())
            dense_vec_data = response_json["data"]
            
            # 计算响应时间
            latency = time.time() - start_time
            logger.info(f"Received response in {latency:.2f}s")
            
            # 安全的响应日志（只记录元数据）
            response_metadata = {
                "object": response_json.get("object"),
                "model": response_json.get("model"),
                "usage": response_json.get("usage"),
                "data_count": len(dense_vec_data)
            }
            logger.info(f"Response metadata: {json.dumps(response_metadata)}")
            
            # 调试日志：记录前3个向量的维度
            if dense_vec_data:
                sample_info = [
                    {"index": i, "vec_len": len(item["embedding"])} 
                    for i, item in enumerate(dense_vec_data[:3])
                ]
                logger.debug(f"Sample vector dimensions: {sample_info}")
            
            # 构建结果
            result_list = [
                {"dense_vec": emb_vec["embedding"]}
                for emb_vec in dense_vec_data
            ]
            return {"result": result_list}
            
        except Exception as e:
            # 增强错误日志
            error_details = f"Error: {type(e).__name__} - {str(e)}"

            # 尝试获取OpenAI错误详情
            if hasattr(e, 'response'):
                try:
                    status_code = getattr(e.response, "status_code", "N/A")
                    error_body = e.response.text if hasattr(e.response, "text") else "N/A"
                    error_details += f" | HTTP {status_code}: {error_body[:200]}"
                except Exception as parse_err:
                    error_details += f" | Failed to parse error: {parse_err}"

            logger.error(f"Embedding request failed (attempt {attempt + 1}): {error_details}")

            # 判断是否限流
            is_rate_limited = error_details and "429" in error_details
            if is_rate_limited:
                if attempt < len(rate_limit_backoff):
                    wait_time = rate_limit_backoff[attempt]
                    logger.warning(f"Rate limited (429). Retrying after {wait_time}s...")
                    time.sleep(wait_time)
                    attempt += 1
                    continue
                else:
                    logger.error("Exceeded max retries due to rate limiting.")
                    break
            else:
                if attempt < other_error_max_retries:
                    logger.warning(f"Non-429 error. Retrying after {other_error_wait}s...")
                    time.sleep(other_error_wait)
                    attempt += 1
                    continue
                else:
                    logger.error("Exceeded max retries for non-429 errors.")
                    break
    
    # 最终错误处理
    raise RuntimeError(f"Failed to get embeddings after retries. Model config: {emb_info}")

def calculate_cosine(query, search_list, embedding_model_id="") -> list[float]:
    query_vector_scores = []
    query_vector = get_embs([query], embedding_model_id=embedding_model_id)["result"][0]["dense_vec"]
    contents = []
    for item in search_list:
        contents.append(item["snippet"])
    contents_vector = get_embs(contents, embedding_model_id=embedding_model_id)["result"]
    for item in contents_vector:
        vec1 = np.array(query_vector)
        vec2 = np.array(item["dense_vec"])

        # calculate dot product
        dot_product = np.dot(vec1, vec2)

        # calculate norm
        norm_vec1 = np.linalg.norm(vec1)
        norm_vec2 = np.linalg.norm(vec2)

        # calculate cosine similarity
        cosine_sim = dot_product / (norm_vec1 * norm_vec2)
        query_vector_scores.append(cosine_sim)

    return query_vector_scores

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


def generate_document_id(title, snippet):
    """根据文档的标题和摘要生成一个唯一的ID"""
    # 生成标题和摘要的哈希值
    hash_title = hashlib.md5(title.encode('utf-8')).hexdigest()
    hash_snippet = hashlib.md5(snippet.encode('utf-8')).hexdigest()
    # 组合哈希值生成最终的文档ID
    return f"{hash_title}-{hash_snippet}"


def generate_md5(content_str):
    # 创建一个md5 hash对象
    md5_obj = hashlib.md5()

    # 对字符串进行编码，因为md5需要bytes类型的数据
    md5_obj.update(content_str.encode('utf-8'))

    # 获取十六进制的MD5值
    md5_value = md5_obj.hexdigest()

    return md5_value


def add_field_to_mapping(index_name, field_name, field_type):
    """
    向指定索引的映射中添加一个新字段。

    参数:
    index_name -- 索引名称
    field_name -- 要添加的字段名称
    field_type -- 字段类型
    """
    try:
        # 获取当前索引的映射
        current_mapping = es.indices.get_mapping(index=index_name)[index_name]['mappings']['properties']
        print(f"当前索引 '{index_name}' 的映射: {current_mapping}")
        logger.info(f"当前索引 '{index_name}' 的映射: {current_mapping}")

        # 添加新字段到映射
        new_mapping = {
            field_name: {
                "type": field_type,
            }
        }
        current_mapping.update(new_mapping)

        # 更新索引的映射
        es.indices.put_mapping(index=index_name, body={"properties": current_mapping})
        print(f"字段 '{field_name}' 已添加到索引 '{index_name}' 的映射中。")
        logger.info(f"字段 '{field_name}' 已添加到索引 '{index_name}' 的映射中。")
    except Exception as e:
        print(f"更新映射时发生错误: {e}")
        logger.error(f"更新映射时发生错误: {e}")

def check_index_exists(index_name):
    return es.indices.exists(index=index_name)

def create_index(index_name, index_settings=None, mappings=None):
    if index_settings is None:
        index_settings = {
            "index": {
                "number_of_shards": 1,
                "number_of_replicas": 0
            }
        }
    if mappings is None:
        mappings = {}

    body = {
        "settings": index_settings,
        "mappings": mappings
    }
    es.indices.create(index=index_name, body=body)

def create_index_if_not_exists(index_name, index_settings=None, mappings=None):
    """
    创建Elasticsearch索引，如果索引不存在的话。

    参数:
    index_name -- 索引名称
    index_settings -- 索引设置，字典格式
    mappings -- 映射定义，字典格式
    """

    # 尝试获取索引信息，如果索引不存在，将会抛出NotFound异常
    exists = check_index_exists(index_name)
    if exists:
        logger.info(f"索引 '{index_name}' 已存在。")
        return True
    else:
        # 索引不存在，创建索引
        logger.info(f"创建索引 '{index_name}'...")
        create_index(index_name, index_settings=index_settings, mappings=mappings)
        logger.info(f"索引 '{index_name}' 创建成功。")
        return False


def bulk_add_index_data(index_name, kb_name, data):
    """使用 helpers.bulk() 批量上传数据到指定的 Elasticsearch 索引，并返回操作状态"""
    actions = []
    # 首先校验index命名是否合法
    is_index_valid, reason = validate_index_name(index_name)  # 创建普通文本类型索引
    if not is_index_valid:
        print("index invalid")
        return {"success": False, "uploaded": len(data), "error": reason}
    # ============ 若索引不存在则新建 ============
    create_index_if_not_exists(index_name)
    # ============ 若索引不存在则新建 ============
    # 提前设置doc_meta字段mapping，避免自动mapping
    es_mapping.update_doc_meta_mapping(index_name)
    for item in data:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
        # doc_id = generate_document_id(item['title'], item['snippet']) # 弃用，
        # 生成一个随机的UUID，相当于不校验重复
        cont_str = kb_name + item["content"] + item["file_name"] + str(item["meta_data"]["chunk_current_num"])
        content_id = generate_md5(cont_str)
        doc_id = uuid.uuid4()
        item['content_id'] = content_id
        # logger.info(f"content_id:{content_id}")
        item['chunk_id'] = doc_id
        item['kb_name'] = kb_name
        action = {
            "_op_type": "index",
            "_index": index_name,
            "_id": doc_id,
            "_source": item
        }
        actions.append(action)
    # 执行批量操作
    try:
        helpers.bulk(es, actions)
        # es.indices.refresh(index=index_name)
        logger.info(
            f"bulk_add_index_data, index_name:'{index_name}', kb_name:'{kb_name}' 添加成功。文档数量: {len(actions)}")
        return {"success": True, "uploaded": len(actions), "error": None}
    except Exception as e:
        # 专门处理批量索引错误
        error_count = len(e.errors)
        logger.error(f"批量索引失败！共 {error_count}/{len(actions)} 个文档索引失败")
        # 打印每个失败文档的详细原因
        for i, error in enumerate(e.errors[:5]):  # 最多打印前5个错误
            doc_id = error['index'].get('_id', '未指定ID')
            reason = error['index']['error']['reason']
            error_type = error['index']['error']['type']
            logger.error(f"失败文档 #{i+1} - ID: {doc_id}")
            logger.error(f"    → 错误类型: {error_type}")
            logger.error(f"    → 原因: {reason}")
        if error_count > 5:
            logger.error(f"...... 另有 {error_count-5} 个错误未显示 ......")
    
        # 如果批量操作失败，返回失败状态和错误信息
        logger.info(f"bulk_add_index_data have err, index_name:'{index_name}',kb_name:{kb_name}, item:{item}")
        import traceback
        logger.error(traceback.format_exc())
        return {"success": False, "uploaded": len(actions), "error": str(e)}


def bulk_add_cc_index_data(index_name, kb_name, data):
    """(用于content 主控索引添加数据) 使用 helpers.bulk() 批量上传数据到指定的 Elasticsearch 索引，并返回操作状态"""
    actions = []
    # ============== 直接往里添加，固定 id  ==============
    try:
        for item in data:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
            cont_str = kb_name + item["content"] + item["file_name"] + str(item["meta_data"]["chunk_current_num"])
            doc_id = generate_md5(cont_str)
            # print(doc_id)
            item['content_id'] = doc_id
            item['chunk_id'] = doc_id
            item['kb_name'] = kb_name
            action = {
                "_op_type": "index",  # 使用index,已存在就覆盖
                "_index": index_name,
                "_id": doc_id,
                "_source": item
            }
            actions.append(action)

        # 提前设置doc_meta字段mapping，避免自动mapping
        es_mapping.update_doc_meta_mapping(index_name)
        # 执行批量操作
        helpers.bulk(es, actions)
        # es.indices.refresh(index=index_name)
        logger.info(f"bulk_add_cc_index_data, index_name:'{index_name}',kb_name:{kb_name} 添加成功。{len(actions)}")
        return {"success": True, "uploaded": len(actions), "error": None}
    except Exception as e:
        # 如果批量操作失败，返回失败状态和错误信息
        logger.info(f"bulk_add_cc_index_data have err, index_name:'{index_name}',kb_name:{kb_name}, item:{item}")
        import traceback
        logger.error(traceback.format_exc())
        return {"success": False, "uploaded": len(actions), "error": str(e)}

def get_index_update_content_actions(index_name, kb_name, content_id, chunk_current_num, child_chunk_current_num, update_data):
    """
    主控表更新

    参数:
    index_name: 索引名称
    kb_name: 知识库名称
    file_name: 文件名
    update_data: 更新的数据, list
    chunk_id: 可选，指定特定chunk_id时使用
    """
    # 构建查询条件
    must_conditions = [
        {"term": {"kb_name": kb_name}},
        {"term": {"content_id": content_id}},
        {"term": {"meta_data.chunk_current_num": chunk_current_num}},
        {"term": {"meta_data.child_chunk_current_num": child_chunk_current_num}}
    ]

    query = {
        "query": {
            "bool": {
                "must": must_conditions
            }
        }
    }

    scan_kwargs = {
        "index": index_name,
        "query": query,
        "scroll": "1m",
        "size": 1
    }

    upsert_data = []
    for doc in helpers.scan(es, **scan_kwargs):
        data = {
            "chunk_id": doc["_source"]["chunk_id"],
        }
        data.update(update_data)
        upsert_data.append(data)

    actions = []
    for item in upsert_data:
        doc_id = item["chunk_id"]
        action = {
            "_op_type": "update",
            "_index": index_name,
            "_id": doc_id,
            "doc": item
        }
        actions.append(action)

    return actions


def bulk_add_uk_index_data(index_name, data):
    """(用于userid 的所有 kb_name映射表索引添加数据) 使用 helpers.bulk() 批量上传数据到指定的 Elasticsearch 索引，并返回操作状态"""
    actions = []
    # ============== 直接往里添加，固定 id  ==============
    try:
        for item in data:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
            if not item["kb_id"]:  # 如果不传递，则生成一个
                # cont_str = item["index_name"] + item["userId"] + item["kb_name"]
                # doc_id = generate_md5(cont_str)
                doc_id = uuid.uuid4()  # 不关注重复
                # print(doc_id)
                item['item_id'] = doc_id
                item['kb_id'] = doc_id
            else:
                doc_id = item["kb_id"]
                item['item_id'] = doc_id
            action = {
                "_op_type": "index",  # 使用index,已存在就覆盖
                "_index": index_name,
                "_id": doc_id,
                "_source": item
            }
            actions.append(action)

        # 执行批量操作
        helpers.bulk(es, actions)
        res = es.indices.refresh(index=index_name)

        # logger.info(f"{res}： bulk_add_uk_index_data  ----- {data}")
        return {"success": True, "uploaded": len(actions), "error": None}
    except Exception as e:
        # 如果批量操作失败，返回失败状态和错误信息
        return {"success": False, "uploaded": len(actions), "error": str(e)}


def cc_bulk_upsert_index_data(index_name, data):
    """使用 helpers.bulk() 批量上传数据到指定的 Elasticsearch 主控表索引，并返回操作状态"""
    actions = []
    for item in data:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
        doc_id = item["content_id"]
        # print(doc_id)
        action = {
            "_op_type": "update",
            "_index": index_name,
            "_id": doc_id,
            "doc": item
        }
        actions.append(action)

    # 执行批量操作
    try:
        helpers.bulk(es, actions)
        es.indices.refresh(index=index_name)
        return {"success": True, "upserted": len(actions), "error": None}
    except Exception as e:
        # 如果批量操作失败，返回失败状态和错误信息
        return {"success": False, "upserted": len(actions), "error": str(e)}


def snippet_bulk_add_index_data(index_name, kb_name, data):
    """使用 helpers.bulk() 批量上传数据到指定的 Elasticsearch 索引，并返回操作状态"""
    actions = []
    # #首先校验index命名是否合法
    # is_index_valid,reason = validate_index_name(index_name)
    # if not is_index_valid:
    #     print("index invalid")
    #     return {"success": False, "uploaded": len(data), "error": reason}
    for item in data:  # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
        try:
            if 'chunk_current_num' in item.get("meta_data"):  # 添加 md5 的 content_id
                if "parent_snippet" in item:
                    content_str = kb_name + item["parent_snippet"] + item["meta_data"]["file_name"] + str(
                        item["meta_data"]["chunk_current_num"])
                    item.pop("parent_snippet")
                else:
                    content_str = kb_name + item["snippet"] + item["meta_data"]["file_name"] + str(
                        item["meta_data"]["chunk_current_num"])

                content_id = generate_md5(content_str)
                item['content_id'] = content_id
            # 生成一个随机的UUID，相当于不校验重复
            doc_id = uuid.uuid4()
            item['chunk_id'] = doc_id
            item['kb_name'] = kb_name
            # print(doc_id)
            action = {
                "_op_type": "index",
                "_index": index_name,
                "_id": doc_id,
                "_source": item
            }
            actions.append(action)
        except Exception as e:
            # 如果在处理单个文档时出现异常，记录错误但继续处理其他文档
            return {"success": False, "uploaded": len(actions), "error": str(e)}

    # 执行批量操作
    try:
        # 提前设置doc_meta字段mapping，避免自动mapping
        es_mapping.update_doc_meta_mapping(index_name)
        helpers.bulk(es, actions)
        # es.indices.refresh(index=index_name)
        return {"success": True, "uploaded": len(actions), "error": None}
    except Exception as e:
        # 如果批量操作失败，返回失败状态和错误信息
        return {"success": False, "uploaded": len(actions), "error": str(e)}


def rescore_bm25_score(index_name, query, search_by="snippet", search_list = []):
    content_ids = []
    for item in search_list:
        content_ids.append(item['content_id'])
    """根据content id 进行过滤，重计算bm 25得分，并按分数从高到低排序"""
    search_body = {
        "query": {
            "bool": {
                "filter": [
                    {
                        "bool": {
                            "should": [
                                {"terms": {"content_id": content_ids}},
                                {"terms": {"content_id.keyword": content_ids}}
                            ],
                            "minimum_should_match": 1
                        }
                    }
                ],
                "must": [
                    {
                        "match": {
                            search_by: query
                        }
                    }
                ]
            }
        },
        "size": len(search_list),  # 指定返回的文档数量
        "sort": [
            {"_score": {"order": "desc"}}  # 按分数降序排序
        ]
    }

    response = es.search(index=index_name, body=search_body)

    search_list = []
    scores = []
    # 遍历搜索结果，填充列表
    for hit in response['hits']['hits']:
        hit_data = hit['_source']
        hit_data["score"] = hit['_score']
        search_list.append(hit_data)
        scores.append(hit['_score'])

    # 构建结果字典
    result_dict = {
        "search_list": search_list,
        "scores": scores
    }

    return result_dict

def search_data_keyword_recall(index_name, kb_name, keywords, top_k, min_score, search_by="labels",
                            filter_file_name_list=[]):
    """根据查询检索数据，仅返回分数高于 min_score 的文档，并按分数从高到低排序"""
    labels_list = keywords.keys()
    # 构建查询体，每个匹配项都有相同的权重
    should_clauses = []
    for label in labels_list:
        should_clauses.append({
            "term": {
                search_by: {
                    "value": label,
                    "boost": keywords[label]  # 每个匹配项的权重
                }
            }
        })
    must_clauses = [
        {"term": {"kb_name": kb_name}}
    ]
    # 如果提供了文件名过滤列表，则添加文件名过滤条件
    if filter_file_name_list:
        must_clauses.append({"terms": {"title.keyword": filter_file_name_list}})


    search_body = {
        "query": {
            "bool": {
                "must": must_clauses,
                "should": should_clauses,
                "minimum_should_match": 1  # 至少匹配一个
            }
        },
        "min_score": min_score,
        "size": top_k,
        "sort": [
            {"_score": {"order": "desc"}}  # 按分数降序排序
        ]
    }
    logger.info(f"search_data_keyword_recall, es index: {index_name}, search body: {search_body}")

    response = es.search(index=index_name, body=search_body)

    search_list = []
    scores = []

    # 遍历搜索结果
    for hit in response['hits']['hits']:
        hit_data = hit['_source']
        hit_data["score"] = hit['_score']
        search_list.append(hit_data)
        scores.append(hit['_score'])

    # 构建结果字典
    result_dict = {
        "search_list": search_list,
        "scores": scores
    }

    return result_dict


def search_data_text_recall(index_name, kb_name, query, top_k, min_score, search_by="snippet",
                            filter_file_name_list=[]):
    """根据查询检索数据，仅返回分数高于 min_score 的文档，并按分数从高到低排序"""
    if filter_file_name_list:
        search_body = {
            "query": {
                "bool": {
                    "must": [
                        # 假设 'search_by' 是你要查询的字段名称，query 是具体的查询值
                        {"match": {search_by: query}},
                        {"term": {"kb_name": kb_name}},
                        {"terms": {"title.keyword": filter_file_name_list}},
                    ],
                }
            },
            "min_score": min_score,
            "size": top_k,  # 指定返回的文档数量
            "sort": [
                {"_score": {"order": "desc"}}  # 按分数降序排序
            ]
        }
    else:
        # ============== TFIDF 通道召回数据 ==========
        search_body = {
            "query": {
                "bool": {
                    "must": [
                        {
                            "match": {
                                search_by: query  # 假设 'search_by' 是你要查询的字段名称，query 是具体的查询值
                            }
                        },
                        {
                            "term": {
                                "kb_name": kb_name
                            }
                        },
                    ]
                }
            },
            "min_score": min_score,
            "size": top_k,  # 指定返回的文档数量
            "sort": [
                {"_score": {"order": "desc"}}  # 按分数降序排序
            ]
        }

    response = es.search(index=index_name, body=search_body)

    search_list = []
    scores = []
    # 遍历搜索结果，填充列表
    for hit in response['hits']['hits']:
        hit_data = hit['_source']
        hit_data["score"] = hit['_score']
        search_list.append(hit_data)
        scores.append(hit['_score'])

    # 构建结果字典
    result_dict = {
        "search_list": search_list,
        "scores": scores
    }

    return result_dict


def search_text_title_list(index_name, kb_name, query, top_k, min_score=0):
    """根据查询检索数据，仅返回分数高于 min_score 的文档，并按分数从高到低排序"""
    search_body = {
        "query": {
            "bool": {
                "must": [
                    {"match": {"title": query}},
                    {"term": {"kb_name": kb_name}},
                ],
            }
        },
        "min_score": min_score,
        "size": top_k,  # 指定返回的文档数量
        "sort": [
            {"_score": {"order": "desc"}}  # 按分数降序排序
        ],
        "collapse": {
            "field": "title.keyword"  # 根据 title 字段去重
        }
    }
    response = es.search(index=index_name, body=search_body)
    search_list = []
    scores = []
    # 遍历搜索结果，填充列表
    for hit in response['hits']['hits']:
        search_list.append(hit['_source']["title"])
        scores.append(hit['_score'])
    # 构建结果字典
    result_dict = {
        "filename_list": search_list,
        "scores": scores
    }
    return result_dict

def is_field_exist(index_name:str, field_name:str)-> (bool, dict):
    mapping = es.indices.get_mapping(index=index_name)
    properties = mapping[index_name].get('mappings', {}).get('properties', {})

    if field_name not in properties:
        return False, properties

    return True, properties

def search_data_knn_recall(index_name, kb_names, query, top_k, min_score, filter_file_name_list=[], embedding_model_id=""):
    """根据查询检索数据，仅返回分数高于 min_score 的文档，并按分数从高到低排序，支持多知识库"""

    query_vector = get_embs([query], embedding_model_id=embedding_model_id)["result"][0]["dense_vec"]
    field_name = f"q_{len(query_vector)}_content_vector"
    # 检查索引映射以确定使用哪个字段
    field_exist, properties = is_field_exist(index_name, field_name)

    # 如果指定维度的字段不存在
    if not field_exist:
        # 只有1024维度可以回退到默认字段
        if len(query_vector) == 1024 and "content_vector" in properties:
            logger.info(f"es 索引 {index_name} 字段 {field_name} 不存在，回退到默认字段 content_vector")
            field_name = "content_vector"
        else:
            # 其他维度不存在对应字段时抛出错误
            available_fields = [k for k in properties.keys() if 'content_vector' in k]
            error_msg = f"向量维度不支持: {field_name} 字段在索引映射中不存在，可用的向量字段: {available_fields}"
            logger.error(error_msg)
            raise ValueError(error_msg)
    else:
        logger.info(f"es 索引 {index_name} 使用向量字段: {field_name} 执行向量检索")

    # ============== KNN 通道召回数据 ==========
    if filter_file_name_list:
        search_body = {
            "knn": {
                "field": field_name,
                "query_vector": query_vector,
                "filter": [
                    {"terms": {"kb_name": kb_names}},
                    {"terms": {"file_name": filter_file_name_list}},
                ],
                "k": 10,
                "num_candidates": max(50, top_k),
            },
            "min_score": min_score,
            "size": top_k,  # 指定返回的文档数量
            "sort": [
                {"_score": {"order": "desc"}}  # 按分数降序排序
            ],
            "_source": ["content", "embedding_content", "file_name", "kb_name", "chunk_id", "meta_data", "content_id", "is_parent"],
            # 指定您希望返回的字段
        }
    else:
        search_body = {
            "knn": {
                "field": field_name,
                "query_vector": query_vector,
                "filter": [
                    {"terms": {"kb_name": kb_names}}
                ],
                "k": 10,
                "num_candidates": max(50, top_k),
            },
            "min_score": min_score,
            "size": top_k,  # 指定返回的文档数量
            "sort": [
                {"_score": {"order": "desc"}}  # 按分数降序排序
            ],
            "_source": ["content", "embedding_content", "file_name", "kb_name", "chunk_id", "meta_data", "content_id", "is_parent"],
            # 指定您希望返回的字段
        }

    response = es.search(index=index_name, body=search_body)

    search_list = []
    scores = []
    # 遍历搜索结果，填充列表
    for hit in response['hits']['hits']:
        hit_data = hit['_source']
        hit_data["score"] = hit['_score']
        # 父子分段模式
        if "is_parent" in hit_data and not hit_data["is_parent"]:
            hit_data["content"] = hit_data["embedding_content"]
        search_list.append(hit_data)
        scores.append(hit['_score'])

    # 构建结果字典
    result_dict = {
        "search_list": search_list,
        "scores": scores
    }

    return result_dict


def get_kb_name_list(index_name):
    """ 获取 index_name 下 所有的知识库名称的集合"""
    body = {
        # "query": {
        #     "match_all": {}  # 使用 match_all 查询来获取所有文档
        # },
        "aggs": {
            "unique_res": {
                "terms": {
                    "field": "kb_name",
                    "size": 100000,  # 根据需要设置大小
                }
            }
        },
        "size": 0  # 不需要原始文档，只用于聚合
    }

    response = es.search(index=index_name, body=body)
    unique_res = [bucket['key'] for bucket in response['aggregations']['unique_res']['buckets']]
    return unique_res


def get_uk_kb_name_list(index_name, user_id):
    """ 获取 userid 的所有 kb_name 映射表下 某个 user_id 所有的知识库名称的集合"""
    body = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"userId": user_id}},
                    # {"match": {"userId": user_id}},
                ]
            }
        },
        "aggs": {
            "unique_res": {
                "terms": {
                    "field": "kb_name",
                    "size": 100000,  # 根据需要设置大小
                }
            }
        },
        "size": 0  # 不需要原始文档，只用于聚合
    }

    response = es.search(index=index_name, body=body)
    unique_res = [bucket['key'] for bucket in response['aggregations']['unique_res']['buckets']]
    return unique_res


def get_uk_kb_id_list(index_name, user_id):
    """ 获取 userid 的所有 kb_name 映射表下 某个 user_id 所有的知识库名称的集合"""
    kb_id_list = []
    kb_name_list = get_uk_kb_name_list(index_name, user_id)
    for kb_name in kb_name_list:
        kb_id_list.append(get_uk_kb_id(user_id, kb_name))
    return kb_id_list


def get_file_name_list(index_name, kb_name):
    """ 获取 index_name 下某个知识库下 file_name 的集合"""
    body = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"kb_name": kb_name}},
                ]
            }
        },
        "aggs": {
            "unique_res": {
                "terms": {
                    "field": "file_name",
                    "size": 100000,  # 根据需要设置大小
                }
            }
        },
        "size": 0  # 不需要原始文档，只用于聚合
    }

    response = es.search(index=index_name, body=body)
    unique_res = [bucket['key'] for bucket in response['aggregations']['unique_res']['buckets']]
    return unique_res


def get_file_download_link_list(index_name, kb_name):
    """ 获取 index_name 下某个知识库下 file_download_link 的集合"""
    body = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"kb_name": kb_name}},
                ]
            }
        },
        "aggs": {
            "unique_res": {
                "terms": {
                    "field": "meta_data.file_name.keyword",
                    "size": 100000,  # 根据需要设置大小
                }
            }
        },
        "size": 0  # 不需要原始文档，只用于聚合
    }

    response = es.search(index=index_name, body=body)
    unique_res = [bucket['key'] for bucket in response['aggregations']['unique_res']['buckets']]
    return unique_res


def fetch_all_documents(index_name):
    """ 从指定索引中获取所有文档 """
    query = {
        "query": {
            "match_all": {}  # 使用 match_all 查询来获取所有文档
        }
    }
    results = helpers.scan(
        es,
        query=query,
        index=index_name,
        scroll='5m',  # 每次滚动窗口持续时间
        size=1000  # 每个批次返回的文档数量
    )

    for doc in results:
        yield doc


def fetch_all_kb_documents(index_name, kb_name):
    """ 从指定索引中获取 某个知识库下所有文档 """
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"kb_name": kb_name}},
                ]
            }
        }
    }
    results = helpers.scan(
        es,
        query=query,
        index=index_name,
        scroll='5m',  # 每次滚动窗口持续时间
        size=1000  # 每个批次返回的文档数量
    )

    for doc in results:
        yield doc


def delete_data_by_kbname_file_names(index_name: str, kb_name: str, file_names: list):
    """根据索引名和 kb_name字段和 file_name 字段 精确匹配删除文档，并返回删除操作的状态"""
    # # === term查询默认是进行精确匹配的，它不会进行分词处理，而是会匹配整个字段的值。但是，term查询是区分大小写的 ===
    # query = {
    #     "query": {
    #         "term": {
    #             "title.keyword": title,
    #             "kb_name.keyword": kb_name,
    #         }
    #     }
    # }
    # === 想要确保file_name和kb_name两个字段完全等于某个字符串，你可以使用bool查询来组合这两个条件 ===
    # 构建查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"kb_name": kb_name}},
                    {"terms": {"file_name": file_names}},
                ]
            }
        }
    }
    # # =============== 一次性删除所有匹配的文档 ===============
    # try:
    #     response = es.delete_by_query(index=index_name, body=query)
    #     delete_status = {
    #         "success": True,
    #         "deleted_count": response['deleted'],
    #         "failures": response.get('failures', [])
    #     }
    #     es.indices.refresh(index=index_name)
    # except Exception as e:
    #     delete_status = {
    #         "success": False,
    #         "error": str(e),
    #         "deleted_count": 0
    #     }
    #
    # return delete_status
    # # =============== 一次性删除所有匹配的文档 ===============
    # =============== 分batch 删除所有匹配的文档 ===============
    try:
        deleted_num = 0
        # 使用 scan API 获取匹配的文档 ID
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }
        delete_actions = []
        for doc in helpers.scan(es, **scan_kwargs):
            delete_actions.append({
                "_op_type": "delete",
                "_index": index_name,
                "_id": doc['_id']
            })
            if len(delete_actions) >= DELETE_BACTH_SIZE:
                logger.info(
                    f"索引 '{index_name}' kb_name:{kb_name} ,file_names:{file_names} 删除文档数量: {deleted_num}")
                # 使用 bulk API 批量删除
                res = helpers.bulk(es, delete_actions)
                deleted_num += res[0]
                delete_actions = []  # 清空 delete_actions
        if len(delete_actions) > 0:
            logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_names:{file_names} 删除文档数量: {deleted_num}")
            # 最后的残留 bulk API 也批量删除
            res = helpers.bulk(es, delete_actions)
            deleted_num += res[0]
        delete_status = {
            "success": True,
            "deleted": deleted_num
        }
        es.indices.refresh(index=index_name)
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
        }

    return delete_status


def delete_data_by_kbname_file_name(index_name: str, kb_name: str, file_name: str):
    """根据索引名和 kb_name字段和 file_name 字段 精确匹配删除文档，并返回删除操作的状态"""
    # === 想要确保file_name和kb_name两个字段完全等于某个字符串，你可以使用bool查询来组合这两个条件 ===
    # 构建查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    {"term": {"file_name": file_name}},
                ]
            }
        }
    }
    # # =============== 一次性删除所有匹配的文档 ===============
    # try:
    #     response = es.delete_by_query(index=index_name, body=query)
    #     delete_status = {
    #         "success": True,
    #         "deleted_count": response['deleted'],
    #         "failures": response.get('failures', [])
    #     }
    #     es.indices.refresh(index=index_name)
    # except Exception as e:
    #     delete_status = {
    #         "success": False,
    #         "error": str(e),
    #         "deleted_count": 0
    #     }
    #
    # return delete_status
    # # =============== 一次性删除所有匹配的文档 ===============
    # =============== 分batch 删除所有匹配的文档 ===============
    try:
        deleted_num = 0
        # 使用 scan API 获取匹配的文档 ID
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }
        if check_index_exists(index_name): #兼容老索引没有file_content_xxx索引
            delete_actions = []
            logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_name:{file_name} 删除文档数量: {deleted_num}")
            for doc in helpers.scan(es, **scan_kwargs):
                delete_actions.append({
                    "_op_type": "delete",
                    "_index": index_name,
                    "_id": doc['_id']
                })
                if len(delete_actions) >= DELETE_BACTH_SIZE:
                    # 使用 bulk API 批量删除
                    res = helpers.bulk(es, delete_actions)
                    deleted_num += res[0]
                    delete_actions = []  # 清空 delete_actions
                    logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_name:{file_name} 删除文档数量: {deleted_num}")
            if len(delete_actions) > 0:
                # 最后的残留 bulk API 也批量删除
                res = helpers.bulk(es, delete_actions)
                deleted_num += res[0]
                logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_name:{file_name} 删除文档数量: {deleted_num}")
            es.indices.refresh(index=index_name)
        delete_status = {
            "success": True,
            "deleted": deleted_num
        }
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
        }

    return delete_status


def delete_data_by_kbname_title(index_name, kb_name, title):
    """根据索引名和 kb_name字段和 title字段 精确匹配删除文档，并返回删除操作的状态"""
    # # === term查询默认是进行精确匹配的，它不会进行分词处理，而是会匹配整个字段的值。但是，term查询是区分大小写的 ===
    # query = {
    #     "query": {
    #         "term": {
    #             "title.keyword": title,
    #             "kb_name.keyword": kb_name,
    #         }
    #     }
    # }
    # === 想要确保title和kb_name两个字段完全等于某个字符串，你可以使用bool查询来组合这两个条件 ===
    # 构建查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"title.keyword": title}},
                    {"term": {"kb_name": kb_name}},
                    # ===== 此 title 字段类型为 text，会进行分词处理，所以不能使用 match =====
                    # {"match": {"title": title}},
                    # {"match": {"kb_name": kb_name}}
                ]
            }
        }
    }
    # # =============== 一次性删除所有匹配的文档 ===============
    # try:
    #     response = es.delete_by_query(index=index_name, body=query)
    #     delete_status = {
    #         "success": True,
    #         "deleted_count": response['deleted'],
    #         "failures": response.get('failures', [])
    #     }
    #     es.indices.refresh(index=index_name)
    # except Exception as e:
    #     delete_status = {
    #         "success": False,
    #         "error": str(e),
    #         "deleted_count": 0
    #     }
    #
    # return delete_status
    # # =============== 一次性删除所有匹配的文档 ===============
    # =============== 分batch 删除所有匹配的文档 ===============
    try:
        deleted_num = 0
        # 使用 scan API 获取匹配的文档 ID
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }
        delete_actions = []
        logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_name:{title} 删除文档数量: {deleted_num}")
        for doc in helpers.scan(es, **scan_kwargs):
            delete_actions.append({
                "_op_type": "delete",
                "_index": index_name,
                "_id": doc['_id']
            })
            if len(delete_actions) >= DELETE_BACTH_SIZE:
                # 使用 bulk API 批量删除
                res = helpers.bulk(es, delete_actions)
                deleted_num += res[0]
                delete_actions = []  # 清空 delete_actions
                logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_name:{title} 删除文档数量: {deleted_num}")
        if len(delete_actions) > 0:
            # 最后的残留 bulk API 也批量删除
            res = helpers.bulk(es, delete_actions)
            deleted_num += res[0]
            logger.info(f"索引 '{index_name}' kb_name:{kb_name} ,file_name:{title} 删除文档数量: {deleted_num}")
        delete_status = {
            "success": True,
            "deleted": deleted_num
        }
        es.indices.refresh(index=index_name)
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
        }

    return delete_status


def delete_data_by_kbname(index_name: str, kb_name: str):
    """根据索引名和 kb_name字段 精确匹配删除文档，并返回删除操作的状态"""
    # === 想要确保file_name和kb_name两个字段完全等于某个字符串，你可以使用bool查询来组合这两个条件 ===
    # 构建查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                ]
            }
        }
    }
    try:
        deleted_num = 0
        # 使用 scan API 获取匹配的文档 ID
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }
        if check_index_exists(index_name): #兼容老知识库没有file_control_xxx索引
            delete_actions = []
            for doc in helpers.scan(es, **scan_kwargs):
                delete_actions.append({
                    "_op_type": "delete",
                    "_index": index_name,
                    "_id": doc['_id']
                })
                if len(delete_actions) >= DELETE_BACTH_SIZE:
                    logger.info(f"索引 '{index_name}' kb_name:{kb_name} , 删除文档数量: {deleted_num}")
                    # 使用 bulk API 批量删除
                    res = helpers.bulk(es, delete_actions)
                    deleted_num += res[0]
                    delete_actions = []  # 清空 delete_actions
            if len(delete_actions) > 0:
                logger.info(f"索引 '{index_name}' kb_name:{kb_name} , 删除文档数量: {deleted_num}")
                # 最后的残留 bulk API 也批量删除
                res = helpers.bulk(es, delete_actions)
                deleted_num += res[0]
            es.indices.refresh(index=index_name)
        delete_status = {
            "success": True,
            "deleted": deleted_num
        }
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
        }

    return delete_status


def delete_uk_data_by_kbname(userId: str, kb_name: str):
    """根据索引名和 kb_name字段和 file_name 字段 精确匹配删除文档，并返回删除操作的状态"""
    # === 想要确保file_name和kb_name两个字段完全等于某个字符串，你可以使用bool查询来组合这两个条件 ===
    # 构建查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    {"term": {"userId": userId}},
                    # {"match": {"kb_name": kb_name}},
                    # {"match": {"userId": userId}},
                ]
            }
        }
    }
    try:
        response = es.delete_by_query(index=KBNAME_MAPPING_INDEX, body=query)
        delete_status = {
            "success": True,
            "deleted_count": response['deleted'],
            "failures": response.get('failures', [])
        }
        es.indices.refresh(index=KBNAME_MAPPING_INDEX)
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
            "deleted_count": 0
        }

    return delete_status


def delete_index(index_name):
    """根据索引名删除整个索引，并返回操作的状态"""
    try:
        response = es.indices.delete(index=index_name)
        # 如果索引成功删除，通常响应中会包含 acknowledged = True
        delete_status = {
            "success": response.get('acknowledged', False),
            "error": None
        }
    except Exception as e:
        # 捕获异常，如索引不存在或其他Elasticsearch错误
        delete_status = {
            "success": False,
            "error": str(e)
        }

    return delete_status


def get_cc_file_content_list(index_name: str, kb_name: str, file_name: str, page_size: int, search_after: int):
    """ 获取 主控表中 知识片段的分页展示 """
    # ======== 分页查询参数 =============
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    {"term": {"file_name": file_name}},
                ]
            }
        },
        #"search_after": [search_after],  # 初始化search_after参数
        "from": search_after,
        "size": page_size,
        "sort": {"meta_data.chunk_current_num": {"order": "asc"}},  # 确保按照文档ID升序排序
    }
    # 执行查询
    response = es.search(
        index=index_name,
        body=query
        #size=page_size
    )

    # 获取当前页的文档列表
    page_hits = response['hits']['hits']
    res_content_list = []
    for doc in page_hits:
        res_content_list.append(doc['_source'])

    # 获取匹配总数
    total_hits = response['hits']['total']['value']

    return {
        "content_list": res_content_list,
        "chunk_total_num": int(total_hits)
    }

def update_chunk_labels_mapping(index_name: str):
    field_exist, _ = is_field_exist(index_name, "labels")
    if not field_exist:
        # 如果 labels 字段不存在，添加它
        es.indices.put_mapping(
            index=index_name,
            body={
                "properties": {
                    "labels": {
                        "type": "keyword",
                    }
                }
            }
        )
        logger.info(f"已为索引 '{index_name}' 添加 labels 字段映射")

def get_cc_index_update_label_actions(index_name, kb_name, file_name, update_data, chunk_id=None):
    """
    主控表更新

    参数:
    index_name: 索引名称
    kb_name: 知识库名称
    file_name: 文件名
    update_data: 更新的数据, list
    chunk_id: 可选，指定特定chunk_id时使用
    """
    # 构建查询条件
    must_conditions = [
        {"term": {"kb_name": kb_name}},
        {"term": {"file_name": file_name}}
    ]

    # 如果指定了chunk_id，则添加到查询条件中
    if chunk_id:
        # use content_id, content_id always is chunk_id, and chunk_id may be not keyword type
        must_conditions.append({"term": {"content_id": chunk_id}})

    query = {
        "query": {
            "bool": {
                "must": must_conditions
            }
        }
    }

    scan_kwargs = {
        "index": index_name,
        "query": query,
        "scroll": "1m",
        "size": 1 if chunk_id else 100
    }

    upsert_data = []
    for doc in helpers.scan(es, **scan_kwargs):
        data = {
            "content_id": doc["_source"]["content_id"],
            "labels": update_data
        }

        upsert_data.append(data)

    actions = []
    for item in upsert_data:
        doc_id = item["content_id"]
        action = {
            "_op_type": "update",
            "_index": index_name,
            "_id": doc_id,
            "doc": item
        }
        actions.append(action)

    return actions


def update_index_data(index_actions: dict, mapping_update_func=None):
    """
    索引数据更新函数

    参数:
    index_actions: 索引名到操作列表的映射
    mapping_update_func: 可选的mapping更新函数
    """
    actions = []
    index_names = []

    try:
        for index_name, action in index_actions.items():
            actions.extend(action)
            index_names.append(index_name)
            # 如果提供了mapping更新函数，则执行
            if mapping_update_func:
                mapping_update_func(index_name)

        # 执行批量操作
        helpers.bulk(es, actions)
        for index_name in index_names:
            es.indices.refresh(index=index_name)
        result = {
            "code": 0,
            "message": "success",
            "updated_count": len(actions)
        }
        return result
    except Exception as e:
        # 如果批量操作失败，返回失败状态和错误信息
        logger.info(f"update index data error: {e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        return result

def update_chunk_labels(index_actions: dict):
    """更新chunk标签"""
    return update_index_data(index_actions, update_chunk_labels_mapping)


def delete_chunks_by_content_ids(index_name, kb_name, content_ids):
    """
    根据知识库名称和content_id列表删除分段

    参数:
    index_name: 索引名称
    kb_name: 知识库名称,
    content_ids: content_id列表
    """
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    {
                        "bool": {
                            "should": [
                                {"terms": {"content_id": content_ids}},
                                {"terms": {"content_id.keyword": content_ids}}   #兼容老索引
                            ],
                            "minimum_should_match": 1
                        }
                    }
                ]
            }
        }
    }

    try:
        deleted_num = 0
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }

        delete_actions = []
        for doc in helpers.scan(es, **scan_kwargs):
            delete_actions.append({
                "_op_type": "delete",
                "_index": index_name,
                "_id": doc['_id']
            })
            if len(delete_actions) >= DELETE_BACTH_SIZE:
                logger.info(f"索引 '{index_name}' kb_name:{kb_name} , 删除文档数量: {deleted_num}")
                # 使用 bulk API 批量删除
                res = helpers.bulk(es, delete_actions)
                deleted_num += res[0]
                delete_actions = []  # 清空 delete_actions
        if len(delete_actions) > 0:
            logger.info(f"索引 '{index_name}' kb_name:{kb_name} , 删除文档数量: {deleted_num}")
            # 最后的残留 bulk API 也批量删除
            res = helpers.bulk(es, delete_actions)
            deleted_num += res[0]
        es.indices.refresh(index=index_name)
        delete_status = {
            "success": True,
            "deleted": deleted_num
        }
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
        }

    return delete_status



def delete_child_chunks(index_name, kb_name, content_id, chunk_current_num, child_chunk_current_nums):
    """
    根据知识库名称和child_chunk_current_nums列表删除子分段

    参数:
    index_name: 索引名称
    kb_name: 知识库名称,
    child_chunk_current_nums: 子段列表
    """
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    {
                        "bool": {
                            "should": [
                                {"term": {"content_id": content_id}},
                                {"term": {"content_id.keyword": content_id}}   #兼容老索引
                            ],
                            "minimum_should_match": 1
                        }
                    },
                    {"term": {"meta_data.chunk_current_num": chunk_current_num}},
                    {"terms": {"meta_data.child_chunk_current_num": child_chunk_current_nums}}
                ]
            }
        }
    }

    try:
        deleted_num = 0
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }

        delete_actions = []
        for doc in helpers.scan(es, **scan_kwargs):
            delete_actions.append({
                "_op_type": "delete",
                "_index": index_name,
                "_id": doc['_id']
            })
            if len(delete_actions) >= DELETE_BACTH_SIZE:
                logger.info(f"索引 '{index_name}' kb_name:{kb_name} , 删除文档数量: {deleted_num}")
                # 使用 bulk API 批量删除
                res = helpers.bulk(es, delete_actions)
                deleted_num += res[0]
                delete_actions = []  # 清空 delete_actions
        if len(delete_actions) > 0:
            logger.info(f"索引 '{index_name}' kb_name:{kb_name} , 删除文档数量: {deleted_num}")
            # 最后的残留 bulk API 也批量删除
            res = helpers.bulk(es, delete_actions)
            deleted_num += res[0]
        es.indices.refresh(index=index_name)
        delete_status = {
            "success": True,
            "deleted": deleted_num
        }
    except Exception as e:
        delete_status = {
            "success": False,
            "error": str(e),
        }

    return delete_status

def add_file(file_index_name, kb_name, file_name, file_meta):
    try:
        file_id = generate_md5(kb_name + file_name)
        file_doc = {
            "file_id": file_id,
            "kb_name": kb_name,
            "file_name": file_name,
            "meta_data": file_meta
        }

        # 写入file_index_name索引
        es.index(
            index=file_index_name,
            id=file_id,
            body=file_doc
        )

        # 刷新索引
        es.indices.refresh(index=file_index_name)
        logger.info(f"新增文件记录: file_name={file_name}, kb_name={kb_name}")
        result = {
            "code": 0,
            "message": "success",
        }
        return result
    except Exception as e:
        logger.info(f"新增文件 error: {e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        return result

def sync_file_record(file_index_name, cc_index_name, kb_name, file_name):
    """
    如果file_id在file_index_name索引中不存在，则从cc_index_name中获取符合条件的文档，
    并将file相关信息写入file_index_name索引中

    参数:
    file_index_name: 文件索引名称
    cc_index_name: 主控索引名称
    file_name: 文件名称
    kb_name: 知识库名称
    """
    # 从cc_index_name中查找符合条件的文档
    cc_query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"file_name": file_name}},
                    {"term": {"kb_name": kb_name}}
                ]
            }
        },
        "size": 1  # 只获取一个符合条件的文档
    }

    cc_response = es.search(index=cc_index_name, body=cc_query)

    # 如果在cc_index_name中找到了文档
    if cc_response['hits']['hits']:
        cc_doc = cc_response['hits']['hits'][0]['_source']

        file_id = generate_md5(kb_name + file_name)
        # 根据cc_index_name的结构为file_doc赋值
        file_doc = {
            "file_id": file_id,
            "kb_name": kb_name,
            "file_name": file_name,
            "meta_data": {
                # 从cc_doc的meta_data中提取相关信息
                "bucket_name": cc_doc.get("meta_data", {}).get("bucket_name", ""),
                "chunk_total_num": cc_doc.get("meta_data", {}).get("chunk_total_num", 0),
                "doc_meta": cc_doc.get("meta_data", {}).get("doc_meta", []),
                "download_link": cc_doc.get("meta_data", {}).get("download_link", ""),
                "object_name": cc_doc.get("meta_data", {}).get("object_name", "")
            }
        }

        # 写入file_index_name索引
        es.index(
            index=file_index_name,
            id=file_id,
            body=file_doc
        )

        # 刷新索引
        es.indices.refresh(index=file_index_name)
        logger.info(f"成功同步文件记录: file_name={file_name}, kb_name={kb_name}")


def allocate_chunk_nums(file_index_name: str, cc_index_name: str, kb_name: str, file_name: str, count: int):
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"file_name": file_name}},
                    {"term": {"kb_name": kb_name}}
                ]
            }
        }
    }

    response = es.search(index=file_index_name, body=query)

    # 如果在file_index_name中没有找到对应的记录，兼容老知识库
    if response['hits']['total']['value'] == 0:
        sync_file_record(file_index_name, cc_index_name, kb_name, file_name)
        # 再次查询
        response = es.search(index=file_index_name, body=query)

    # 如果找到了记录，则使用script原子性地增加chunk_total_num
    if response['hits']['hits']:
        doc = response['hits']['hits'][0]
        file_id = doc['_id']

        # 使用script原子性地增加chunk_total_num
        script_body = {
            "script": {
                "source": """
                    if (ctx._source.meta_data.chunk_total_num == null) {
                        ctx._source.meta_data.chunk_total_num = params.count;
                    } else {
                        ctx._source.meta_data.chunk_total_num += params.count;
                    }
                """,
                "lang": "painless",
                "params": {
                    "count": count
                }
            }
        }

        try:
            # 执行更新操作，注意 refresh 和 _source 作为参数传递
            update_response = es.update(
                index=file_index_name,
                id=file_id,
                body=script_body,
                refresh=True,
                source=True
            )

            # 获取更新后的值
            meta_data = (update_response.get("get", {})
                .get("_source", {})
                .get("meta_data", {})
            )

            # 计算起始编号
            chunk_total_num = meta_data.get("chunk_total_num", 0)

            return {
                "code": 0,
                "message": "",
                "data": {
                    "chunk_total_num": chunk_total_num,
                    "allocated_count": count,
                    "file_name": file_name,
                    "kb_name": kb_name,
                    "meta_data": meta_data
                }
            }
        except Exception as e:
            logger.error(f"分配chunk编号时出错: {e}")
            return {
                "code": 1,
                "message": f"分配chunk编号失败: {str(e)}"
            }
    else:
        return {
            "code": 1,
            "message": "未找到文件记录"
        }



def allocate_child_chunk_nums(cc_index_name: str, kb_name: str, file_name: str, chunk_id: str, count: int):
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"chunk_id": chunk_id}},
                    {"term": {"file_name": file_name}},
                    {"term": {"kb_name": kb_name}}
                ]
            }
        }
    }

    response = es.search(index=cc_index_name, body=query)

    # 如果找到了记录，则使用script原子性地增加chunk_total_num
    if response['hits']['hits']:
        doc = response['hits']['hits'][0]
        doc_id = doc['_id']

        # 使用script原子性地增加child_chunk_total_num
        script_body = {
            "script": {
                "source": """
                    if (ctx._source.child_chunk_total_num == null) {
                        ctx._source.child_chunk_total_num = params.count;
                    } else {
                        ctx._source.child_chunk_total_num += params.count;
                    }
                """,
                "lang": "painless",
                "params": {
                    "count": count
                }
            }
        }

        try:
            # 执行更新操作，注意 refresh 和 _source 作为参数传递
            update_response = es.update(
                index=cc_index_name,
                id=doc_id,
                body=script_body,
                refresh=True,
                source=True
            )

            # 获取更新后的值
            doc_data = (update_response.get("get", {}).get("_source", {}))
            meta_data = doc_data.get("meta_data", {})

            # 计算起始编号
            child_chunk_total_num = doc_data.get("child_chunk_total_num", 0)
            content = doc_data.get("content")
            response_data = {
                    "content": content,
                    "child_chunk_total_num": child_chunk_total_num,
                    "allocated_count": count,
                    "file_name": file_name,
                    "kb_name": kb_name,
                    "meta_data": meta_data
                }
            if "labels" in doc_data:
                response_data["labels"] = doc_data["labels"]

            return {
                "code": 0,
                "message": "",
                "data": response_data
            }
        except Exception as e:
            logger.error(f"分配child_chunk编号时出错: {e}")
            return {
                "code": 1,
                "message": f"分配child_chunk编号失败: {str(e)}"
            }
    else:
        return {
            "code": 1,
            "message": "未找到文件记录"
        }


def update_cc_content_status(index_name, kb_name, file_name, content_id, status, on_off_switch):
    """根据 on_off_switch 或 content_id更新知识库文件片段状态 """
    if on_off_switch in [True, False]:  # 一键启停
        query = {
            "query": {
                "bool": {
                    "must": [
                        {"term": {"kb_name": kb_name}},
                        {"term": {"file_name": file_name}},
                        # {"match": {"kb_name": kb_name}},
                        # {"match": {"file_name": file_name}},
                    ]
                }
            }
        }
        scan_kwargs = {
            "index": index_name,
            "query": query,
            "scroll": "1m",
            "size": 100  # 每次返回的文档数量
        }
        upsert_data = []
        for doc in helpers.scan(es, **scan_kwargs):
            upsert_data.append({
                "content_id": doc["_source"]["content_id"],
                "status": on_off_switch,
            })
        upsert_res = cc_bulk_upsert_index_data(index_name, upsert_data)
        if upsert_res["success"]:
            result = {
                "code": 0,
                "message": "success",
            }
            return result
        else:  # 一键启停失败
            result = {
                "code": 1,
                "message": upsert_res["error"],
            }
            return result
    else:  # 单个 content_id 启停
        upsert_data = [{"content_id": content_id, "status": status}]
        upsert_res = cc_bulk_upsert_index_data(index_name, upsert_data)
        if upsert_res["success"]:
            result = {
                "code": 0,
                "message": "success",
            }
            return result
        else:  # 单个 content_id 启停失败
            result = {
                "code": 1,
                "message": upsert_res["error"],
            }
            return result

def get_child_contents(index_name, kb_name, content_id):
    """ 获取子分段"""
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    {
                        "bool": {
                            "should": [
                                {"term": {"content_id": content_id}},
                                {"term": {"content_id.keyword": content_id}}
                            ],
                            "minimum_should_match": 1
                        }
                    }
                ]
            }
        },
        "size": 500  # 增加返回的条数
    }
    response = es.search(index=index_name, body=query)
    # 遍历搜索结果，填充列表
    result = []
    for hit in response["hits"]["hits"]:
        cleaned_hit = {k: v for k, v in hit['_source'].items()
                   if (not k.startswith('vector') and k != 'content_vector')}
        embedding_content = cleaned_hit["embedding_content"]
        cleaned_hit["content"] = embedding_content
        cleaned_hit.pop("embedding_content")
        result.append(cleaned_hit)

    # 获取匹配总数
    total_hits = response['hits']['total']['value']

    return {
        "parent_chunk_id": content_id,
        "child_content_list": result,
        "child_chunk_total_num": int(total_hits)
    }


def get_cc_contents(index_name, kb_name, content_id_list):
    """ 获取文本分块状态用于进行检索后过滤。"""
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"kb_name": kb_name}},
                    {
                        "bool": {
                            "should": [
                                {"terms": {"content_id": content_id_list}},
                                {"terms": {"content_id.keyword": content_id_list}}
                            ],
                            "minimum_should_match": 1
                        }
                    }
                ]
            }
        },
        "size": 500  # 增加返回的条数
    }
    response = es.search(index=index_name, body=query)
    # 遍历搜索结果，填充列表
    result = []
    for hit in response["hits"]["hits"]:
        result.append(hit['_source'])
    # ========= 返回 =========
    return result

def get_cc_content_status(index_name, kb_name, content_id_list):
    """ 获取文本分块状态用于进行检索后过滤。"""
    response = get_cc_contents(index_name, kb_name, content_id_list)
    useful_content_id_list = []
    # 遍历搜索结果，填充列表
    for hit in response:
        if hit["status"]:
            useful_content_id_list.append(hit["content_id"])
    # ========= 返回 =========
    return useful_content_id_list


def get_uk_kb_id(userId, kb_name):
    """ 获取知识库映射的 kb_id """
    kb_id = ""
    logger.info(f"userId:{userId},kb_name:{kb_name} ====== get_uk_kb_id")
    # 查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"userId": userId}},
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"userId": userId}},
                    # {"match": {"kb_name": kb_name}},
                ]
            }
        }
    }
    response = es.search(index=KBNAME_MAPPING_INDEX, body=query)
    # 遍历搜索结果，获取 kb_id
    for hit in response["hits"]["hits"]:
        kb_id = hit['_source']["kb_id"]
    # ========= 返回 =========
    if not kb_id:
        kb_id = get_maas_kb_id(userId, kb_name)  # 如果没有找到，则从 maas 知识库中获取
    logger.info(f"userId:{userId},kb_name:{kb_name} 对应的 kb_id 为:{kb_id}")
    return kb_id


def get_uk_kb_emb_model_id(userId, kb_name):
    """ 获取知识库映射的 embedding_model_id  """
    embedding_model_id = ""
    logger.info(f"userId:{userId},kb_name:{kb_name} ====== get_uk_kb_emb_model_id")
    # 查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"userId": userId}},
                    {"term": {"kb_name": kb_name}},
                    # {"match": {"userId": userId}},
                    # {"match": {"kb_name": kb_name}},
                ]
            }
        }
    }
    response = es.search(index=KBNAME_MAPPING_INDEX, body=query)
    # 遍历搜索结果，获取 kb_id
    for hit in response["hits"]["hits"]:
        embedding_model_id = hit['_source']["embedding_model_id"]
    logger.info(f"userId:{userId},kb_name:{kb_name} 对应的 embedding_model_id 为:{embedding_model_id}")
    return embedding_model_id


def update_uk_kb_name(userId, old_kb_name, new_kb_name):
    """ 更新 uk映射表 知识库名 """
    # 查询条件
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"userId": userId}},
                    {"term": {"kb_name": old_kb_name}},
                    # {"match": {"userId": userId}},
                    # {"match": {"kb_name": old_kb_name}},
                ]
            }
        }
    }
    response = es.search(index=KBNAME_MAPPING_INDEX, body=query)
    # 遍历搜索结果，将actions
    actions = []
    for hit in response["hits"]["hits"]:
        # 往索引里插数据，以index的方式，若_id已存在则先删除再添加
        doc_id = hit['_id']
        # print(doc_id)
        action = {
            "_op_type": "update",
            "_index": KBNAME_MAPPING_INDEX,
            "_id": doc_id,
            "doc": {"kb_name": new_kb_name}
        }
        actions.append(action)
    if len(actions) < 1:
        return {'code': 1, 'message': f'没有找到对应的知识库:{old_kb_name}'}
    # 执行更新操作,并返回
    try:
        helpers.bulk(es, actions)
        es.indices.refresh(index=KBNAME_MAPPING_INDEX)
        return {'code': 0, 'message': 'success'}
    except Exception as e:
        # 如果批量操作失败，返回失败状态和错误信息
        return {'code': 1, 'message': f'{e}'}


if __name__ == "__main__":
    # 示例使用
    pass
