# -*- coding: utf-8 -*-
import json
import time
import copy
from copy import deepcopy

from log.logger import logger
from datetime import datetime
from flask import Flask, request, Response

from settings import EMBEDDING_BATCH_SIZE
from settings import INDEX_NAME_PREFIX, SNIPPET_INDEX_NAME_PREFIX, KBNAME_MAPPING_INDEX
import utils.es_util as es_ops
import utils.meta_util as meta_ops
import utils.mapping_util as es_mapping

app = Flask(__name__)

def batch_list(lst: list, batch_size=32):
    """ 切分生成器 """
    for i in range(0, len(lst), batch_size):
        yield lst[i:i + batch_size]


@app.route('/rag/kn/init_kb', methods=['POST'])
def init_kb():
    """ ES 模拟RAG主控 初始化 init_kb 接口"""
    logger.info("--------------------------启动向量库初始化---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    content_index_name = 'content_control_' + index_name
    userId = data.get("userId")
    kb_name = data.get("kb_name")
    kb_id = data["kb_id"]  # 必须字段
    embedding_model_id = data["embedding_model_id"]  # 必须字段
    enable_knowledge_graph = data.get("enable_knowledge_graph", False)
    userId_kb_names = []
    dense_dim = 1024
    try:
        judge_time1 = time.time()

        es_ops.create_index_if_not_exists(content_index_name, mappings=es_mapping.cc_mappings)  # 确保 主控表 已创建
        es_ops.create_index_if_not_exists(KBNAME_MAPPING_INDEX, mappings=es_mapping.uk_mappings)  # 确保 KBNAME_MAPPING_INDEX 已创建
        is_exists = es_ops.create_index_if_not_exists(index_name, mappings=es_mapping.mappings)
        if is_exists:  # 如果之前创建过了，则查询是否有 kb_name
            # kb_names = es_ops.get_kb_name_list(index_name) # 不使用
            kb_names = es_ops.get_uk_kb_name_list(KBNAME_MAPPING_INDEX, userId)  # 从映射表中获取
            logger.info(f"当前用户:{userId},共有知识库：{len(kb_names)}个，分别为{kb_names}")
            judge_time2 = time.time()
            judge_time = judge_time2 - judge_time1
            logger.info(f"--------------------------查询kb_map时间:{judge_time}---------------------------\n")
            s_time1 = time.time()
            if kb_name in kb_names:
                result = {
                    "code": 1,
                    "message": f"已存在同名知识库{kb_name}"
                }
                jsonarr = json.dumps(result, ensure_ascii=False)
                logger.info(f"当前用户:{userId},知识库:{kb_name},ini知识库的接口返回结果为：{jsonarr}")
                s_time2 = time.time()
                s_time = s_time2 - s_time1
                logger.info(f"--------------------------同名知识库判断时间:{s_time}---------------------------\n")
                return jsonarr
            else:
                # ES 需提前 init_kb 添加到 KBNAME_MAPPING_INDEX
                # 获取当前UTC时间
                utc_now = datetime.utcnow()
                # # 上海时区是UTC+8
                # shanghai_tz = pytz.timezone('Asia/Shanghai')
                # # 将UTC时间转换为上海时间
                # shanghai_now = utc_now.replace(tzinfo=pytz.utc).astimezone(shanghai_tz)
                # # 格式化上海时间为年月日时分秒
                # formatted_time = shanghai_now.strftime('%Y-%m-%d %H:%M:%S')
                formatted_time = utc_now.strftime('%Y-%m-%d %H:%M:%S')
                uk_data = [
                    {"index_name": index_name, "userId": userId, "kb_name": kb_name,
                     "creat_time": formatted_time, "kb_id": kb_id, "embedding_model_id": embedding_model_id,
                     "enable_graph": enable_knowledge_graph}
                ]
                es_ops.bulk_add_uk_index_data(KBNAME_MAPPING_INDEX, uk_data)
                # ====== 新建完成，需要获取一下 kb_id,看看是否新建成功 ======
                save_kb_id = es_ops.get_uk_kb_id(userId, kb_name)
                if save_kb_id != kb_id:  # 新建失败，返回错误
                    result = {
                        "code": 1,
                        "message": "ini知识库失败，ES写入失败",
                    }
                    jsonarr = json.dumps(result, ensure_ascii=False)
                    logger.info(f"当前用户:{userId},知识库:{kb_name},ini知识库的接口返回结果为：{jsonarr}")
                    return jsonarr
                # 新建成功，返回
                logger.info(f"当前用户:{userId},知识库:{kb_name},save_kb_id:{save_kb_id}")
                result = {
                    "code": 0,
                    "message": "success"
                }
                jsonarr = json.dumps(result, ensure_ascii=False)
                logger.info(f"当前用户:{userId},知识库:{kb_name},ini知识库的接口返回结果为：{jsonarr}")
                return jsonarr
        else:
            # ES 需提前 init_kb 添加到 KBNAME_MAPPING_INDEX
            # 获取当前UTC时间
            utc_now = datetime.utcnow()
            # # 上海时区是UTC+8
            # shanghai_tz = pytz.timezone('Asia/Shanghai')
            # # 将UTC时间转换为上海时间
            # shanghai_now = utc_now.replace(tzinfo=pytz.utc).astimezone(shanghai_tz)
            # # 格式化上海时间为年月日时分秒
            # formatted_time = shanghai_now.strftime('%Y-%m-%d %H:%M:%S')
            formatted_time = utc_now.strftime('%Y-%m-%d %H:%M:%S')
            # ES 无需提前 init_kb kb_name,直接添加一条 kb_name 记录
            uk_data = [
                {"index_name": index_name, "userId": userId, "kb_name": kb_name,
                 "creat_time": formatted_time, "kb_id": kb_id, "embedding_model_id": embedding_model_id}
            ]
            es_ops.bulk_add_uk_index_data(KBNAME_MAPPING_INDEX, uk_data)
            # ====== 新建完成，需要获取一下 kb_id,看看是否新建成功 ======
            save_kb_id = es_ops.get_uk_kb_id(userId, kb_name)
            if save_kb_id != kb_id:  # 新建失败，返回错误
                result = {
                    "code": 1,
                    "message": "ini知识库失败，ES写入失败",
                }
                jsonarr = json.dumps(result, ensure_ascii=False)
                logger.info(f"当前用户:{userId},知识库:{kb_name},ini知识库的接口返回结果为：{jsonarr}")
                return jsonarr
            # 新建成功，返回
            logger.info(f"当前用户:{userId},知识库:{kb_name},save_kb_id:{save_kb_id}")
            result = {
                "code": 0,
                "message": "success"
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"当前用户:{userId},知识库:{kb_name},ini知识库的接口返回结果为：{jsonarr}")
            return jsonarr

    except Exception as e:
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},ini知识库的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/add', methods=['POST'])
def add_vector_data():
    """ 往 ES 中建向量索引数据，当前方法要校验索引名 kb_name 是否已存在"""
    logger.info("--------------------------启动数据添加---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    content_index_name = 'content_control_' + index_name
    userId = data.get("userId")
    kb_name = data.get("kb_name")
    kb_id = data.get("kb_id")
    embedding_model_id = es_ops.get_uk_kb_emb_model_id(userId, kb_name)
    doc_list = data.get("data")
    userId_kb_names = []
    cc_doc_list = []  # content主控表的数据
    cc_duplicate_list = []
    if not kb_id:  # 如果没有传入 kb_id,则从映射表中获取
        kb_id = es_ops.get_uk_kb_id(userId, kb_name)  # 从映射表中获取 kb_id ,添加往里传 kb_id
        if not kb_id:  # 如果映射表中没有，则返回错误
            result = {
                "code": 1,
                "message": f"{kb_name}知识库不存在"
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            return jsonarr
    # # **************** 校验 kb_name 是否已经初始化过 ****************
    # userId_kb_ids = es_ops.get_uk_kb_id_list(KBNAME_MAPPING_INDEX, userId)  # 从映射表中获取
    # if kb_id not in userId_kb_ids:
    #     result = {
    #         "code": 1,
    #         "message": f"{kb_id}知识库不存在"
    #     }
    #     jsonarr = json.dumps(result, ensure_ascii=False)
    #     logger.info(f"{userId},/rag/kn/add的接口返回结果为：{jsonarr},userId_kb_names:{userId_kb_ids}")
    #     return jsonarr
    # # **************** 校验 kb_name 是否已经初始化过 ****************
    # ========= 将 content 主控表数据过滤好 =============
    for doc in copy.deepcopy(doc_list):
        cc_str = str(doc["content"]) + doc["file_name"] + str(doc["meta_data"]["chunk_current_num"])
        if cc_str not in cc_duplicate_list:
            doc.pop("embedding_content")  # 去掉不需要的字段
            doc["status"] = True  # 初始化启停状态
            if "is_parent" in doc:
                doc["is_parent"] = True
                doc["child_chunk_total_num"] = doc["meta_data"]["child_chunk_total_num"]
                doc["meta_data"].pop("child_chunk_current_num")
                doc["meta_data"].pop("child_chunk_total_num")
            cc_doc_list.append(doc)
            cc_duplicate_list.append(cc_str)
    # ========= 将 content 主控表数据过滤好 =============
    for doc in doc_list:
        doc.pop("labels", None)  # 去掉不需要的字段, labels 只写content 主控表

    try:
        # ========= 将 embedding_content 编码好向量 =============
        content_vector_exist = False
        mapping_properties = {}
        for batch_doc in batch_list(doc_list, batch_size=EMBEDDING_BATCH_SIZE):
            res = es_ops.get_embs([x["embedding_content"] for x in batch_doc], embedding_model_id=embedding_model_id)
            dense_vector_dim = len(res["result"][0]["dense_vec"]) if res["result"] else 1024
            field_name = f"q_{dense_vector_dim}_content_vector"
            if dense_vector_dim == 1024:
                # 兼容老索引，避免创建两个1024 dim的向量字段
                if not mapping_properties:
                    content_vector_exist, mapping_properties = es_ops.is_field_exist(index_name, "content_vector")
                if content_vector_exist:
                    logger.info(f"es 索引 {index_name} 字段 {field_name} 存在，回退到默认字段 content_vector")
                    field_name = "content_vector"

            for i, x in enumerate(batch_doc):
                if len(batch_doc) != len(res["result"]):
                    raise RuntimeError(f"Error getting embeddings:{batch_doc}")
                x[field_name] = res["result"][i]["dense_vec"]
        # ========= 将 embedding_content 编码好向量 =============
        es_result = es_ops.bulk_add_index_data(index_name, kb_id, doc_list)  # 注意 存储的时候传入 kb_id
        logger.info(f"{es_result}")
        es_cc_result = es_ops.bulk_add_cc_index_data(content_index_name, kb_id, cc_doc_list)  # 注意 存储的时候传入 kb_id
        if es_result["success"] and es_cc_result["success"]:  # bulk_add_index_data 成功了则返回
            result = {
                "code": 0,
                "message": "success"
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"当前用户:{userId},知识库:{kb_name},add的接口返回结果为：{jsonarr}")
            return jsonarr
        else:  # bulk_add_index_data 报错了则返回错误信息
            result = {
                "code": 1,
                "message": es_result.get("error", "") + es_cc_result.get("error", "")
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"当前用户:{userId},知识库:{kb_name},add的接口返回结果为：{jsonarr}")
            return jsonarr
    except Exception as e:
        result = {
            "code": 2,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},add的接口返回结果为：{jsonarr}")
        return jsonarr
    finally:
        logger.info(f"{userId},{kb_name},bulk_add end")

@app.route('/rag/kn/get_kb_info', methods=['POST'])
def get_kb_info():
    """ 查询知识库是否开启知识图谱"""
    logger.info("-----------------------启动知识库info查询-------------------\n")
    data = request.get_json()
    userId = data.get("userId")
    kb_name = data.get("kb_name")
    try:
        # ******** 先检查 是否有新建 index ***********
        es_ops.create_index_if_not_exists(KBNAME_MAPPING_INDEX, mappings=es_mapping.uk_mappings)  # 确保 KBNAME_MAPPING_INDEX 已创建
        kb_info = es_ops.get_uk_kb_info(userId, kb_name)
        logger.info(f"当前用户:{userId},知识库:{kb_name}, kb_info: {kb_info}")
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "kb_info": kb_info
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库info查询接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库是否开启知识图谱查询的接口返回结果为：{jsonarr}")
        return jsonarr

@app.route('/rag/kn/list_kb_names', methods=['POST'])
def list_kb_names():
    logger.info("--------------------------启动知识库查询---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    try:
        # ******** 先检查 是否有新建 index ***********
        es_ops.create_index_if_not_exists(KBNAME_MAPPING_INDEX, mappings=es_mapping.uk_mappings)  # 确保 KBNAME_MAPPING_INDEX 已创建
        is_exists = es_ops.create_index_if_not_exists(index_name, mappings=es_mapping.mappings)
        # ******** 先检查 是否有新建 index ***********
        # userId_kb_names = es_ops.get_kb_name_list(index_name) # 不使用此方式
        userId_kb_names = es_ops.get_uk_kb_name_list(KBNAME_MAPPING_INDEX, userId)  # 从映射表中获取
        logger.info(f"/rag/kn/list_kb_names:用户{index_name}共有{len(userId_kb_names)}个知识库")
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "kb_names": userId_kb_names
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库查询的接口返回结果为：{jsonarr}")
        return jsonarr

    except Exception as e:
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库查询的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/list_file_names', methods=['POST'])
def list_file_names():
    logger.info("--------------------------启动文件查询---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    kb_id = data.get("kb_id")
    try:
        if not kb_id:  # 如果没有指定 kb_id，则从映射表中获取
            kb_id = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        # **************** 校验 kb_name 是否已经初始化过 ****************
        userId_kb_ids = es_ops.get_uk_kb_id_list(KBNAME_MAPPING_INDEX, userId)  # 从映射表中获取
        if kb_id not in userId_kb_ids:
            result = {
                "code": 1,
                "message": f"{kb_id}知识库不存在"
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"{userId},/rag/kn/list_file_names的接口返回结果为：{jsonarr},userId_kb_names:{userId_kb_ids}")
            return jsonarr
        # **************** 校验 kb_name 是否已经初始化过 ****************
        file_names = es_ops.get_file_name_list(index_name, kb_id)
        logger.info(f"用户{index_name}的知识库{kb_id}共有{len(file_names)}个文件")
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "file_names": file_names
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_id},文件查询的接口返回结果为：{jsonarr}")
        return jsonarr

    except Exception as e:
        logger.info(f"查询文件名称时发生错误：{str(e)}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_id},文件查询的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/list_file_names_after_filtering', methods=['POST'])
def list_file_names_after_filtering():
    logger.info("--------------------------启动文件过滤查询---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    kb_id = data.get("kb_id")
    filtering_conditions = data.get("filtering_conditions")
    try:
        if not kb_id:  # 如果没有指定 kb_id，则从映射表中获取
            kb_id = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},display_kb_name: {display_kb_name},请求的kb_id为:{kb_id}, filtering_conditions: {filtering_conditions}")

        final_conditions = []
        for condition in filtering_conditions:
            if condition["filtering_kb_name"] == display_kb_name:
                condition["filtering_kb_name"] = kb_id
                final_conditions.append(deepcopy(condition))
        file_names = []
        if final_conditions:
            file_names = meta_ops.search_with_doc_meta_filter(index_name, final_conditions)
        logger.info(f"用户{index_name}的知识库{display_kb_name}过滤后共有{len(file_names)}个文件")
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "file_names": file_names
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_id},文件过滤查询的接口返回结果为：{jsonarr}")
        return jsonarr

    except Exception as e:
        logger.info(f"过滤查询文件名称时发生错误：{str(e)}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_id},文件过滤查询的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/list_file_download_links', methods=['POST'])
def list_file_download_links():
    logger.info("--------------------------启动获取知识库里所有文档的下载链接查询---------------------------")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    kb_id = data.get("kb_id")
    try:
        if not kb_id:  # 如果没有指定 kb_id，则从映射表中获取
            kb_id = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        file_names = es_ops.get_file_download_link_list(index_name, kb_id)
        logger.info(f"用户{index_name}的知识库{kb_id}共有{len(file_names)}个文件的下载链接")
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "file_download_links": file_names
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_id},文件下载链接查询的接口返回结果为：{jsonarr}")
        return jsonarr

    except Exception as e:
        logger.info(f"查询文件下载链接时发生错误：{str(e)}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_id},文件下载链接查询的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/search', methods=['POST'])
def es_knn_search():
    """ 多知识库 KNN检索 """
    logger.info("--------------------------启动向量库检索---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    content_index_name = 'content_control_' + index_name
    userId = data.get("userId")
    display_kb_names = data.get("kb_names")  # list
    kb_names = []
    top_k = data.get("topk", 10)
    query = data.get("question")
    min_score = data.get("threshold", 0)
    filter_file_name_list = data.get("filter_file_name_list", [])
    metadata_filtering_conditions = data.get("metadata_filtering_conditions", [])
    kb_id_2_kb_name = {}
    embedding_model_id = es_ops.get_uk_kb_emb_model_id(userId, display_kb_names[0])
    logger.info(f"用户:{index_name},请求查询的kb_names为:{display_kb_names},embedding_model_id:{embedding_model_id}")
    logger.info(f"用户请求的query为:{query}")
    try:
        # ============= 先检查 kb_names 是不是都存在 =============
        # exists_kb_names = es_ops.get_kb_name_list(index_name) # 不使用
        exists_kb_names = es_ops.get_uk_kb_name_list(KBNAME_MAPPING_INDEX, userId)  # 从映射表中获取
        filtering_conditions = {}
        for condition in metadata_filtering_conditions:
            kb_name = condition["filtering_kb_name"]
            filtering_conditions[kb_name] = condition

        final_conditions = []
        for kb_name in display_kb_names:
            if kb_name not in exists_kb_names:
                result = {
                    "code": 1,
                    "message": f"用户:{index_name}里,{kb_name}知识库不存在"
                }
                jsonarr = json.dumps(result, ensure_ascii=False)
                logger.info(f"\n向量库检索的接口返回结果为：{jsonarr}")
                return jsonarr
            # ======== kb_name 是存在的，则往 kb_names 里添加=======
            kb_id = es_ops.get_uk_kb_id(userId, kb_name)
            kb_names.append(kb_id)  # 从映射表中获取 kb_id ，这是真正的名字
            kb_id_2_kb_name[kb_id] = kb_name
            if kb_name in filtering_conditions:
                condition = filtering_conditions[kb_name]
                condition["filtering_kb_name"] = kb_id
                final_conditions.append(deepcopy(condition))
        meta_filter_file_name_list = []
        if final_conditions:
            meta_filter_file_name_list = meta_ops.search_with_doc_meta_filter(content_index_name, final_conditions)
            logger.info(f"用户请求的query为:{query}, filter_file_name_list: {filter_file_name_list}, meta_filter_file_name_list: {meta_filter_file_name_list}")
            if len(meta_filter_file_name_list) == 0:
                result = {
                    "code": 0,
                    "message": "success",
                    "data": {
                        "search_list": [],
                        "scores": []
                    }
                }
                jsonarr = json.dumps(result, ensure_ascii=False)
                logger.info(f"当前用户:{userId},知识库:{kb_names},query:{query},向量库检索的接口返回结果为：{jsonarr}")
                return jsonarr

        if meta_filter_file_name_list:
            filter_file_name_list = filter_file_name_list + meta_filter_file_name_list
        # ============= 先检查 kb_names 是不是都存在 =============
        # ============= 开始检索召回 ===============
        result_dict = es_ops.search_data_knn_recall(index_name, kb_names, query, top_k, min_score,
                                                    filter_file_name_list=filter_file_name_list,
                                                    embedding_model_id=embedding_model_id)
        search_list = result_dict["search_list"]
        scores = result_dict["scores"]
        for item in search_list:  # 将 kb_id 转换为 kb_name
            item["kb_name"] = kb_id_2_kb_name[item["kb_name"]]
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "search_list": search_list,
                "scores": scores
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_names},query:{query},向量库检索的接口返回结果为：{jsonarr}")
        return jsonarr

    except Exception as e:
        logger.info(f"查询知识库时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_names},query:{query},向量库检索的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/del_kb', methods=['POST'])
def del_kb():
    logger.info("--------------------------启动知识库删除---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    content_index_name = 'content_control_' + index_name
    file_index_name = 'file_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        es_result = es_ops.delete_data_by_kbname(index_name, kb_name)
        es_cc_result = es_ops.delete_data_by_kbname(content_index_name, kb_name)  # 主控表也需要删除
        es_file_result = es_ops.delete_data_by_kbname(file_index_name, kb_name)
        es_uk_result = es_ops.delete_uk_data_by_kbname(userId, display_kb_name)  # uid索引映射表需要删除,传display_kb_name
        if es_result["success"] and es_cc_result["success"] and es_uk_result["success"] and es_file_result["success"]:  # delete_data_by_kbname 成功了则返回
            logger.info(f"用户{index_name},对应的{kb_name}记录删除成功")
            result = {
                "code": 0,
                "message": "success"
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(
                f"当前用户:{userId},知识库:{kb_name},知识库删除的接口返回结果为：{jsonarr},{es_result},{es_cc_result},{es_uk_result}")
            return jsonarr
        else:
            logger.info(
                f"当前用户:{userId},知识库:{kb_name},知识库删除时发生错误：{es_result},{es_cc_result},{es_uk_result},{es_file_result}")
            result = {
                "code": 1,
                "message": es_result.get("error", "") + es_cc_result.get("error", "") + es_file_result.get("error", "")
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"当前用户:{userId},知识库:{kb_name},知识库删除的接口返回结果为：{jsonarr}")
            return jsonarr

    except Exception as e:
        logger.info(f"用户{index_name},对应的{kb_name}知识库删除时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},知识库删除的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/del_files', methods=['POST'])
def del_files():
    logger.info("--------------------------启动文件删除---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_names = data.get("file_names")
    content_index_name = 'content_control_' + index_name
    file_index_name = 'file_control_' + index_name

    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字

        # # =============== 一步删除，不使用 ===================
        # es_result = es_ops.delete_data_by_kbname_file_names(index_name, kb_name, file_names)
        # # =============== 一步删除，不使用 ===================

        # ********* 单独删除，获取每一个文件状态
        success = []
        failed = []
        for file in file_names:
            es_result = es_ops.delete_data_by_kbname_file_name(index_name, kb_name, file)
            es_cc_result = es_ops.delete_data_by_kbname_file_name(content_index_name, kb_name, file)
            es_file_result = es_ops.delete_data_by_kbname_file_name(file_index_name, kb_name, file)
            if es_result["success"] and es_cc_result["success"] and es_file_result["success"]:  # delete_data_by_kbname_file_names 成功了则返回
                logger.info(f"当前用户{index_name}的知识库{kb_name}删除的文档为：{file}")
                success.append(file)
            else:
                logger.info(
                    f"当前用户:{userId},知识库:{kb_name},file_names:{file_names},文件删除时发生错误：{es_result}")
                result = {
                    "code": 1,
                    "message": es_result.get("error", "") + es_cc_result.get("error", "") + es_file_result.get("error", "")
                }
                jsonarr = json.dumps(result, ensure_ascii=False)
                logger.info(
                    f"当前用户:{userId},知识库:{kb_name},file_names:{file_names},知识库删除的接口返回结果为：{jsonarr}")
                return jsonarr

        # ======== 没有报错，则返回成功 ========
        failed = [file for file in file_names if file not in success]
        logger.info(f"----------当前用户:{userId},知识库{kb_name}完成{file_names}的delete--------------")
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "success": success,
                "failed": failed
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},file_names:{file_names},文件删除的接口返回结果为：{jsonarr}")
        return jsonarr

    except Exception as e:
        logger.info(f"知识库{kb_name},{file_names},在文件删除时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},file_names:{file_names},文件删除的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/get_content_list', methods=['POST'])
def get_content_list():
    """ 获取 主控表中 知识片段的分页展示 """
    logger.info("--------------------------获取主控表中知识片段的分页展示---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    page_size = data.get("page_size")
    search_after = data.get("search_after")
    content_index_name = 'content_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},请求的kb_name为:{kb_name},file_name:{file_name},page_size:{page_size},search_after:{search_after}")
        content_result = es_ops.get_cc_file_content_list(content_index_name, kb_name, file_name, page_size, search_after)
        content_list = content_result["content_list"]
        for content in content_list:
            if "is_parent" in content and content["is_parent"]:
                child_result = es_ops.get_child_contents(index_name, kb_name, content["content_id"])
                content["child_chunk_total_num"] = child_result["child_chunk_total_num"]
        result = {
            "code": 0,
            "message": "success",
            "data": content_result
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},page_size:{page_size},search_after:{search_after},知识片段分页查询的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"获取主控表中知识片段的分页展示时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},获取主控表中知识片段的分页展示的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/get_child_content_list', methods=['POST'])
def get_child_content_list():
    """ 获取子片段"""
    logger.info("--------------------------获取子片段---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    chunk_id = data.get("chunk_id")
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},请求的kb_name为:{kb_name},file_name:{file_name},chunk_id:{chunk_id}")
        content_result = es_ops.get_child_contents(index_name, kb_name, chunk_id)
        result = {
            "code": 0,
            "message": "success",
            "data": content_result
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name}, chunk_id:{chunk_id},子分段查询的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"获取子分段时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},获取子分段的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/update_child_chunk', methods=['POST'])
def update_child_chunk():
    logger.info("--------------------------更新知识库子段数据---------------------------\n")
    data = request.get_json()
    userId = data.get("userId")
    index_name = INDEX_NAME_PREFIX + userId
    display_kb_name = data.get("kb_name")  # 显示的名字
    embedding_model_id = es_ops.get_uk_kb_emb_model_id(userId, display_kb_name)
    snippet_index_name = SNIPPET_INDEX_NAME_PREFIX + userId.replace('-', '_')
    chunk_id = data.get("chunk_id")
    child_chunk = data.get("child_chunk")
    chunk_current_num = data.get("chunk_current_num")
    try:
        child_content = child_chunk["child_content"]
        child_chunk_current_num = child_chunk["child_chunk_current_num"]
        index_update_data = {
            "embedding_content": child_content,
        }
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(f"用户:{userId},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name}, chunk_id: {chunk_id}, "
                    f"chunk_current_num: {chunk_current_num}, child_chunk: {child_chunk}")

        res = es_ops.get_embs([child_content], embedding_model_id=embedding_model_id)
        dense_vector_dim = len(res["result"][0]["dense_vec"]) if res["result"] else 1024
        field_name = f"q_{dense_vector_dim}_content_vector"
        if dense_vector_dim == 1024:
            # 兼容老索引，避免创建两个1024 dim的向量字段
            content_vector_exist, mapping_properties = es_ops.is_field_exist(index_name, "content_vector")
            if content_vector_exist:
                logger.info(f"es 索引 {index_name} 字段 {field_name} 存在，回退到默认字段 content_vector")
                field_name = "content_vector"

        index_update_data[field_name] = res["result"][0]["dense_vec"]
        snippet_index_update_data = {
            "snippet": child_content,
        }

        # cc index的content id == chunk id
        index_update_actions = es_ops.get_index_update_content_actions(index_name, kb_name, chunk_id, chunk_current_num,
                                                                       child_chunk_current_num, index_update_data)

        snippet_index_update_actions = es_ops.get_index_update_content_actions(snippet_index_name, kb_name, chunk_id,
                                                                               chunk_current_num, child_chunk_current_num, snippet_index_update_data)

        update_actions = {
            index_name: index_update_actions,
            snippet_index_name: snippet_index_update_actions
        }
        result = es_ops.update_index_data(update_actions)
        json_arr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},更新知识库元数据的接口返回结果为：{json_arr}")
        return json_arr
    except Exception as e:
        logger.info(f"更新知识库元数据时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        json_arr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{display_kb_name},更新知识库元数据的接口返回结果为：{json_arr}")
        return json_arr


@app.route('/rag/kn/update_file_metas', methods=['POST'])
def update_file_metas():
    logger.info("--------------------------更新知识库元数据---------------------------\n")
    data = request.get_json()
    userId = data.get("userId")
    index_name = INDEX_NAME_PREFIX + userId
    display_kb_name = data.get("kb_name")  # 显示的名字
    update_datas = data.get("update_datas")
    file_index_name = 'file_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(f"用户:{userId},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name}, update_datas: {update_datas}")

        # 兼容老版本，没有file index的需要创建
        es_ops.create_index_if_not_exists(file_index_name, mappings=es_mapping.mappings)
        result = meta_ops.update_file_metas(userId, kb_name, update_datas)
        json_arr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},更新知识库元数据的接口返回结果为：{json_arr}")
        return json_arr
    except Exception as e:
        logger.info(f"更新知识库元数据时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        json_arr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{display_kb_name},更新知识库元数据的接口返回结果为：{json_arr}")
        return json_arr

@app.route('/rag/kn/batch_delete_chunks', methods=['POST'])
def batch_delete_chunks():
    logger.info("--------------------------根据fileName和chunk_ids删除分段---------------------------\n")
    data = request.get_json()
    userId = data.get("userId")
    index_name = INDEX_NAME_PREFIX + userId
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    chunk_ids = data.get("chunk_ids")
    content_index_name = 'content_control_' + index_name
    snippet_index_name = SNIPPET_INDEX_NAME_PREFIX + userId.replace('-', '_')
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name},file_name:{file_name}, chunk_ids: {chunk_ids}")

        es_result = es_ops.delete_chunks_by_content_ids(index_name, kb_name, chunk_ids)
        es_cc_result = es_ops.delete_chunks_by_content_ids(content_index_name, kb_name, chunk_ids)  # 主控表也需要删除
        es_snippet_result = es_ops.delete_chunks_by_content_ids(snippet_index_name, kb_name, chunk_ids)
        if es_result["success"] and es_cc_result["success"] and es_snippet_result["success"]:
            logger.info(f"用户{index_name},对应的知识库{kb_name}, chunks: {chunk_ids}记录分段删除成功")
            result = {
                "code": 0,
                "message": "success",
                "data": {
                    "success_count": es_cc_result["deleted"]
                }
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(
                f"当前用户:{userId},知识库:{kb_name},chunks:{chunk_ids}, 分段删除的接口返回结果为：{jsonarr},{es_result},{es_cc_result},{es_snippet_result}")
            return jsonarr
        else:
            logger.info(
                f"当前用户:{userId},知识库:{kb_name},chunks:{chunk_ids}, 分段删除时发生错误：{es_result},{es_cc_result},{es_snippet_result}")
            result = {
                "code": 1,
                "message": es_result.get("error", "") + es_cc_result.get("error", "") + es_snippet_result.get("error", ""),
                "data": {
                    "success_count": 0
                }
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"当前用户:{userId},知识库:{kb_name},chunks:{chunk_ids}, 分段删除的接口返回结果为：{jsonarr}")
            return jsonarr

    except Exception as e:
        logger.info(f"用户{index_name},对应的知识库:{kb_name},chunks:{chunk_ids}, 分段删除时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e),
            "data": {
                "success_count": 0
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},chunks:{chunk_ids}, 分段删除的接口返回结果为：{jsonarr}")
        return jsonarr

@app.route('/rag/kn/delete_child_chunks', methods=['POST'])
def delete_child_chunks():
    logger.info("--------------------------删除子分段---------------------------\n")
    data = request.get_json()
    userId = data.get("userId")
    index_name = INDEX_NAME_PREFIX + userId
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    chunk_id = data.get("chunk_id")
    chunk_current_num = data.get("chunk_current_num")
    child_chunk_current_nums = data.get("child_chunk_current_nums")
    snippet_index_name = SNIPPET_INDEX_NAME_PREFIX + userId.replace('-', '_')

    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name},file_name:{file_name}, "
            f"chunk_id: {chunk_id}, chunk_current_num: {chunk_current_num}, "
            f"child_chunk_current_nums: {child_chunk_current_nums}")

        es_result = es_ops.delete_child_chunks(index_name, kb_name, chunk_id, chunk_current_num, child_chunk_current_nums)
        es_snippet_result = es_ops.delete_child_chunks(snippet_index_name, kb_name, chunk_id, chunk_current_num, child_chunk_current_nums)
        if es_result["success"] and es_snippet_result["success"]:
            logger.info(f"用户{index_name},对应的知识库{kb_name}, chunk: {chunk_id}, "
                        f"child_chunk_current_nums: {child_chunk_current_nums} 记录子分段删除成功")
            result = {
                "code": 0,
                "message": "success"
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(
                f"当前用户:{userId},知识库:{kb_name},chunk:{chunk_id}, child_chunk_current_nums: {child_chunk_current_nums} "
                f"子分段删除的接口返回结果为：{jsonarr},{es_result},{es_snippet_result}")
            return jsonarr
        else:
            logger.info(
                f"当前用户:{userId},知识库:{kb_name},chunk:{chunk_id}, child_chunk_current_nums: {child_chunk_current_nums} "
                f"子分段删除时发生错误：{es_result},{es_snippet_result}")
            result = {
                "code": 1,
                "message": es_result.get("error", "") + es_snippet_result.get("error", ""),
            }
            jsonarr = json.dumps(result, ensure_ascii=False)
            logger.info(f"当前用户:{userId},知识库:{kb_name},chunk:{chunk_id}, "
                        f"child_chunk_current_nums: {child_chunk_current_nums} 子分段删除的接口返回结果为：{jsonarr}")
            return jsonarr

    except Exception as e:
        logger.info(f"用户{index_name},对应的知识库:{kb_name},chunk:{chunk_id}, "
                    f"child_chunk_current_nums: {child_chunk_current_nums} 子分段删除时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e),
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},chunk:{chunk_id}, "
                    f"child_chunk_current_nums: {child_chunk_current_nums} 子分段删除的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/update_chunk_labels', methods=['POST'])
def update_chunk_labels():
    logger.info("--------------------------根据fileName和chunk_id更新标签---------------------------\n")
    data = request.get_json()
    userId = data.get("userId")
    index_name = INDEX_NAME_PREFIX + userId
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    chunk_id = data.get("chunk_id")
    labels = data.get("labels")
    content_index_name = 'content_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name},file_name:{file_name}, chunk_id: {chunk_id}, labels: {labels}")

        index_actions = {
            content_index_name: es_ops.get_cc_index_update_label_actions(content_index_name, kb_name, file_name, labels, chunk_id=chunk_id)
        }
        result = es_ops.update_chunk_labels(index_actions)
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},根据fileName和chunk_id更新知识库chunk 标签的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"根据fileName和chunk_id更新知识库chunk 标签时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},根据fileName和chunk_id更新知识库chunk 标签的接口返回结果为：{jsonarr}")
        return jsonarr

@app.route('/rag/kn/get_content_by_ids', methods=['POST'])
def get_content_by_ids():
    """ 根据content_id获取知识库文件片段 """
    logger.info("--------------------------根据content_id获取知识库文件片段信息---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    content_ids = data.get("content_ids")
    kb_id = data.get("kb_id")
    content_index_name = 'content_control_' + index_name
    try:
        if not kb_id:  # 如果没有传入 kb_id,则从映射表中获取
            kb_id = es_ops.get_uk_kb_id(userId, display_kb_name)
        logger.info(
            f"用户:{userId},请求的kb_name为:{kb_id},content_ids:{content_ids}")
        contents = es_ops.get_cc_contents(content_index_name, kb_id, content_ids)
        for item in contents:  # 将 kb_id 转换为 kb_name
            item["kb_name"] = display_kb_name
        result = {
            "code": 0,
            "message": "success",
            "data": {
                "contents": contents
            }
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_id},content_ids:{content_ids}, 根据content_ids获取片段信息的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"根据content_ids获取分段信息时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_id},content_ids:{content_ids}, 根据content_ids获取片段信息的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/update_content_status', methods=['POST'])
def update_content_status():
    """ 根据content_id更新知识库文件片段状态 """
    logger.info("--------------------------根据content_id更新知识库文件片段状态---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    content_id = data.get("content_id")
    status = data.get("status")
    on_off_switch = data.get("on_off_switch", -1)
    content_index_name = 'content_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{userId},请求的kb_name为:{kb_name},file_name:{file_name},content_id:{content_id},status:{status},on_off_switch:{on_off_switch}")
        result = es_ops.update_cc_content_status(content_index_name, kb_name, file_name, content_id, status,
                                                 on_off_switch)
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},content_id:{content_id},search_after:{status},on_off_switch:{on_off_switch}根据content_id更新知识库文件片段状态的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"根据content_id更新知识库文件片段状态时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},file_name:{file_name},根据content_id更新知识库文件片段状态的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/get_useful_content_status', methods=['POST'])
def get_content_status():
    """ 获取文本分块状态用于进行检索后过滤。 """
    logger.info("--------------------------获取文本分块状态用于进行检索后过滤---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    display_kb_name = data.get("kb_name")  # 显示的名字
    content_id_list = data.get("content_id_list")
    content_index_name = 'content_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(userId, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(f"用户:{userId},请求的kb_name为:{kb_name},content_id_list:{content_id_list}")
        useful_content_id_list = es_ops.get_cc_content_status(content_index_name, kb_name, content_id_list)
        result = {'code': 0, 'message': 'success', 'data': {'useful_content_id_list': useful_content_id_list}}
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},content_id_list:{content_id_list},获取文本分块状态用于进行检索后过滤的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"获取文本分块状态用于进行检索后过滤时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{userId},知识库:{kb_name},content_id_list:{content_id_list},获取文本分块状态用于进行检索后过滤的接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/get_kb_id', methods=['POST'])
def get_kb_id():
    """ 获取某个知识库映射的 kb_id接口 """
    logger.info("--------------------------获取知识库映射的 kb_id接口---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    kb_name = data.get("kb_name")
    logger.info(f"用户:{userId},请求的kb_name为:{kb_name}")
    try:
        kb_id = es_ops.get_uk_kb_id(userId, kb_name)
        result = {'code': 0, 'message': 'success', 'data': {'kb_id': kb_id}}
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},获取知识库映射的 kb_id接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"获取知识库映射的 kb_id接口发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},知识库:{kb_name},获取知识库映射的 kb_id接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/rag/kn/update_kb_name', methods=['POST'])
def update_uk_kb_name():
    """ 更新 uk映射表 知识库名接口 """
    logger.info("--------------------------更新 uk映射表 知识库名接口---------------------------\n")
    data = request.get_json()
    index_name = INDEX_NAME_PREFIX + data.get('userId')
    userId = data.get("userId")
    old_kb_name = data.get("old_kb_name")
    new_kb_name = data.get("new_kb_name")
    logger.info(f"用户:{userId},请求的ole_kb_name为:{old_kb_name},请求的new_kb_name为:{new_kb_name}")
    try:
        result = es_ops.update_uk_kb_name(userId, old_kb_name, new_kb_name)
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},{old_kb_name},{new_kb_name},更新uk映射表知识库名接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"更新uk映射表知识库名接口发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(f"当前用户:{userId},{old_kb_name},{new_kb_name},更新uk映射表知识库名接口返回结果为：{jsonarr}")
        return jsonarr


@app.route('/es/test', methods=['GET', 'POST'])
def test():
    data = request.get_json()
    logger.info(f"request: /es/test , data:{data}")
    cluster_health = es_ops.get_health()
    response = json.dumps({'code': 200, 'msg': 'Success', 'result': cluster_health}, indent=4, ensure_ascii=False)
    return Response(response, mimetype='application/json', status=200)


# ***************** 老的 ES snippet API servers **********************

@app.route('/api/v1/rag/es/bulk_add', methods=['POST'])
def snippet_bulk_add():
    logger.info("request: /api/v1/rag/es/bulk_add")
    data = request.get_json()
    # logger.info('bulk_add request_params: '+ json.dumps(data, indent=4,ensure_ascii=False))

    # index_name = data.get('index_name') 之前拼接好的，弃用
    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    kb_name = data.get('kb_name')
    kb_id = data.get('kb_id')
    doc_list = data.get('doc_list')
    logger.info(f"request: bulk_add_data len:{len(doc_list)}")
    try:
        # ========= 往里面传入的 kb_name是真正指代的 kb_id =======
        if not kb_id:  # 如果没有传入 kb_id,则从映射表中获取
            kb_id = es_ops.get_uk_kb_id(data.get('user_id'), data.get('kb_name'))
        es_ops.create_index_if_not_exists(index_name, mappings=es_mapping.snippet_mappings)
        result = es_ops.snippet_bulk_add_index_data(index_name, kb_id, doc_list)
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': result}, indent=4, ensure_ascii=False)
        logger.info("bulk_add response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("bulk_add response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/bulk_add end")


@app.route('/api/v1/rag/es/add_file', methods=['POST'])
def add_file():
    logger.info("--------------------------新增文件---------------------------\n")
    data = request.get_json()
    user_id = data.get("user_id")
    index_name = INDEX_NAME_PREFIX + user_id
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    file_meta = data.get("file_meta")
    file_index_name = 'file_control_' + index_name

    try:
        kb_name = es_ops.get_uk_kb_id(user_id, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{user_id},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name},file_name:{file_name}, file_meta: {file_meta}")

        es_ops.create_index_if_not_exists(file_index_name, mappings=es_mapping.file_mappings)
        result = es_ops.add_file(file_index_name, kb_name, file_name, file_meta)
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{user_id},知识库:{kb_name},file_name:{file_name},新增文件返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"新增文件时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{user_id},知识库:{kb_name},file_name:{file_name},新增文件返回结果为：{jsonarr}")
        return jsonarr

@app.route('/api/v1/rag/es/allocate_chunks', methods=['POST'])
def allocate_chunks():
    logger.info("--------------------------新增分段时分配chunk---------------------------\n")
    data = request.get_json()
    user_id = data.get("user_id")
    index_name = INDEX_NAME_PREFIX + user_id
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    count = data.get("count")
    content_index_name = 'content_control_' + index_name
    file_index_name = 'file_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(user_id, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{user_id},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name},file_name:{file_name}, insert chunk count: {count}")

        es_ops.create_index_if_not_exists(content_index_name, mappings=es_mapping.cc_mappings)
        es_ops.create_index_if_not_exists(file_index_name, mappings=es_mapping.file_mappings)
        result = es_ops.allocate_chunk_nums(file_index_name, content_index_name, kb_name, file_name, count)
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{user_id},知识库:{kb_name},file_name:{file_name},新增分段分配chunk的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"新增分段分配chunk时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{user_id},知识库:{kb_name},file_name:{file_name},新增分段分配chunk返回结果为：{jsonarr}")
        return jsonarr


@app.route('/api/v1/rag/es/allocate_child_chunks', methods=['POST'])
def allocate_child_chunks():
    logger.info("--------------------------新增子分段时分配chunk---------------------------\n")
    data = request.get_json()
    user_id = data.get("user_id")
    index_name = INDEX_NAME_PREFIX + user_id
    display_kb_name = data.get("kb_name")  # 显示的名字
    file_name = data.get("file_name")
    chunk_id = data.get("chunk_id")
    count = data.get("count")
    content_index_name = 'content_control_' + index_name
    try:
        kb_name = es_ops.get_uk_kb_id(user_id, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字
        logger.info(
            f"用户:{user_id},display_kb_name: {display_kb_name},请求的kb_name为:{kb_name},file_name:{file_name}, insert chunk count: {count}")

        es_ops.create_index_if_not_exists(content_index_name, mappings=es_mapping.cc_mappings)
        result = es_ops.allocate_child_chunk_nums(content_index_name, kb_name, file_name, chunk_id, count)
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{user_id},知识库:{kb_name},file_name:{file_name},新增子分段分配chunk的接口返回结果为：{jsonarr}")
        return jsonarr
    except Exception as e:
        logger.info(f"新增子分段分配chunk时发生错误：{e}")
        result = {
            "code": 1,
            "message": str(e)
        }
        jsonarr = json.dumps(result, ensure_ascii=False)
        logger.info(
            f"当前用户:{user_id},知识库:{kb_name},file_name:{file_name},新增子分段分配chunk返回结果为：{jsonarr}")
        return jsonarr


@app.route('/api/v1/rag/es/search', methods=['POST'])
def snippet_search():
    logger.info("request: /api/v1/rag/es/search")
    data = request.get_json()
    logger.info('search request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))

    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    content_index_name = 'content_control_' + INDEX_NAME_PREFIX + user_id
    kb_name = data.get('kb_name')
    query = data.get('query')
    top_k = int(data.get('top_k', 10))
    min_score = float(data.get('min_score', 0.0))
    search_by = data.get('search_by', "snippet")
    filter_file_name_list = data.get("filter_file_name_list", [])
    metadata_filtering_conditions = data.get("metadata_filtering_conditions", [])
    kb_id_2_kb_name = {}
    try:
        # ========= 往里面传入的 kb_name是真正指代的 kb_id =======
        kb_id = es_ops.get_uk_kb_id(data.get('user_id'), data.get('kb_name'))
        kb_id_2_kb_name[kb_id] = kb_name

        final_conditions = []
        for condition in metadata_filtering_conditions:
            if condition["filtering_kb_name"] == kb_name:
                condition["filtering_kb_name"] = kb_id
                final_conditions.append(deepcopy(condition))

        meta_filter_file_name_list = []
        if final_conditions:
            meta_filter_file_name_list = meta_ops.search_with_doc_meta_filter(content_index_name, final_conditions)
            logger.info(f"用户请求的query为:{query}, filter_file_name_list: {filter_file_name_list}, meta_filter_file_name_list: {meta_filter_file_name_list}")
            if len(meta_filter_file_name_list) == 0:
                result = {
                    "search_list": [],
                    "scores": []
                }
                response = json.dumps({'code': 200, 'msg': 'Success', 'result': result}, indent=4, ensure_ascii=False)
                logger.info("search response: %s", response)
                return Response(response, mimetype='application/json', status=200)

        if meta_filter_file_name_list:
            filter_file_name_list = filter_file_name_list + meta_filter_file_name_list

        result = es_ops.search_data_text_recall(index_name, kb_id, query, top_k, min_score, search_by,
                                                filter_file_name_list=filter_file_name_list)
        search_list = result["search_list"]
        for item in search_list:  # 将 kb_id 转换为 kb_name
            item["kb_name"] = kb_id_2_kb_name[item["kb_name"]]
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': result}, indent=4, ensure_ascii=False)
        logger.info("search response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("search response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/search end")



@app.route('/api/v1/rag/es/keyword_search', methods=['POST'])
def keyword_search():
    logger.info("request: /api/v1/rag/es/keyword_search")
    data = request.get_json()
    logger.info('search request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))
    user_id = data.get('user_id')
    index_name = INDEX_NAME_PREFIX + data.get('user_id')
    content_index_name = 'content_control_' + index_name
    display_kb_name = data.get('kb_name')
    keywords = data.get('keywords')
    top_k = int(data.get('top_k', 10))
    min_score = float(data.get('min_score', 0.0))
    search_by = data.get('search_by', "labels")
    filter_file_name_list = data.get('filter_file_name_list', [])
    metadata_filtering_conditions = data.get('metadata_filtering_conditions', [])
    try:
        kb_id = es_ops.get_uk_kb_id(user_id, display_kb_name)  # 从映射表中获取 kb_id ，这是真正的名字

        final_conditions = []
        for condition in metadata_filtering_conditions:
            if condition["filtering_kb_name"] == display_kb_name:
                condition["filtering_kb_name"] = kb_id
                final_conditions.append(deepcopy(condition))

        meta_filter_file_name_list = []
        if final_conditions:
            meta_filter_file_name_list = meta_ops.search_with_doc_meta_filter(content_index_name, final_conditions)
            logger.info(
                f"filter_file_name_list: {filter_file_name_list}, meta_filter_file_name_list: {meta_filter_file_name_list}")
            if len(meta_filter_file_name_list) == 0:
                result = {
                    "search_list": [],
                    "scores": []
                }
                response = json.dumps({'code': 200, 'msg': 'Success', 'result': result}, indent=4, ensure_ascii=False)
                logger.info("search response: %s", response)
                return Response(response, mimetype='application/json', status=200)

        if meta_filter_file_name_list:
            filter_file_name_list = filter_file_name_list + meta_filter_file_name_list

        result = es_ops.search_data_keyword_recall(content_index_name, kb_id, keywords, top_k, min_score, search_by,
                                                   filter_file_name_list=filter_file_name_list)
        search_list = result["search_list"]
        for item in search_list:  # 将 kb_id 转换为 kb_name
            item["kb_name"] = display_kb_name
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': result}, indent=4, ensure_ascii=False)
        logger.info("search response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("search response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/keyword_search end")


@app.route('/api/v1/rag/es/rescore', methods=['POST'])
def snippet_rescore():
    logger.info("request: /api/v1/rag/es/rescore")
    data = request.get_json()
    logger.info('rescore request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))

    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    query = data.get('query')
    weights = data.get('weights')
    search_by = data.get('search_by', "snippet")
    search_list = data.get('search_list', [])
    display_kb_names = data.get("kb_names")
    embedding_model_id = es_ops.get_uk_kb_emb_model_id(user_id, display_kb_names[0])
    logger.info(f"用户:{index_name},请求查询的kb_names为:{display_kb_names},embedding_model_id:{embedding_model_id}")
    kb_id_2_kb_name = {}
    try:
        for kb_name in display_kb_names:
            kb_id = es_ops.get_uk_kb_id(data.get('user_id'), kb_name)
            kb_id_2_kb_name[kb_id] = kb_name
        result = es_ops.rescore_bm25_score(index_name, query, search_by, search_list)
        search_list = result["search_list"]
        bm25_scores = result["scores"]
        cosine_scores = es_ops.calculate_cosine(query, search_list, embedding_model_id)
        logger.info(f"rescore bm25_scores: {bm25_scores}, cosine_scores: {cosine_scores}")

        def normalize_to_01(scores):
            if len(scores) == 1:
                return [1.0]  # 单个分数归一化为1
            min_score = min(scores)
            max_score = max(scores)
            if min_score == max_score:
                return [1.0 for _ in scores]  # 所有分数相同，统一设为1
            return [(score - min_score) / (max_score - min_score) for score in scores]

        bm25_normalized = normalize_to_01(bm25_scores)
        cosine_normalized = normalize_to_01(cosine_scores)

        final_search_list = []
        for item, text_score, vector_score in zip(search_list, bm25_normalized, cosine_normalized):
            score = weights["vector_weight"] * vector_score + weights["text_weight"] * text_score
            item["score"] = score
            item["kb_name"] = kb_id_2_kb_name[item["kb_name"]]
            final_search_list.append(item)

        final_search_list.sort(key=lambda x: x["score"], reverse=True)
        final_results = {
            "search_list": final_search_list,
            "scores": [item["score"] for item in final_search_list]
        }
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': final_results}, indent=4, ensure_ascii=False)
        logger.info("rescore response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("rescore response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/rescore end")


@app.route('/api/v1/rag/es/search_text_title_list', methods=['POST'])
def search_title_list():
    logger.info("request: /api/v1/rag/es/search_text_title_list")
    data = request.get_json()
    logger.info('search request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))

    # index_name = data.get('index_name') 之前拼接好的，弃用
    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    kb_name = data.get('kb_name')
    query = data.get('query')
    top_k = int(data.get('top_k', 10))
    min_score = float(data.get('min_score', 0.0))
    kb_id_2_kb_name = {}
    try:
        # ========= 往里面传入的 kb_name是真正指代的 kb_id =======
        kb_id = es_ops.get_uk_kb_id(data.get('user_id'), data.get('kb_name'))
        kb_id_2_kb_name[kb_id] = kb_name
        result = es_ops.search_text_title_list(index_name, kb_id, query, top_k, min_score)
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': result}, indent=4, ensure_ascii=False)
        logger.info("search response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("search response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/search_text_title_list end")


@app.route('/api/v1/rag/es/fetch_all', methods=['POST'])
def snippet_fetch_all():
    logger.info("request: /api/v1/rag/es/fetch_all")
    data = request.get_json()
    logger.info('fetch_all request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))

    # index_name = data.get('index_name') 之前拼接好的，弃用
    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    kb_name = data.get('kb_name')
    try:
        documents = es_ops.fetch_all_documents(index_name)
        documents_list = list(documents)
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': documents_list}, indent=4, ensure_ascii=False)
        logger.info("fetch_all response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("fetch_all response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/fetch_all end")


@app.route('/api/v1/rag/es/delete_doc', methods=['POST'])
def snippet_delete_doc_by_kbname_title():
    logger.info("request: /api/v1/rag/es/delete_doc")
    data = request.get_json()
    logger.info('delete_doc_by_title request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))

    # index_name = data.get('index_name') 之前拼接好的，弃用
    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    kb_name = data.get('kb_name')
    title = data.get('title')
    try:
        # ========= 往里面传入的 kb_name是真正指代的 kb_id =======
        kb_id = es_ops.get_uk_kb_id(data.get('user_id'), data.get('kb_name'))
        status = es_ops.delete_data_by_kbname_title(index_name, kb_id, title)
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': status}, indent=4, ensure_ascii=False)
        logger.info("delete_doc_by_title response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("delete_doc_by_title response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/delete_doc end")


@app.route('/api/v1/rag/es/delete_index', methods=['POST'])
def snippet_delete_index_kb_name():
    logger.info("request: /api/v1/rag/es/delete_index")
    data = request.get_json()
    logger.info('delete_index request_params: ' + json.dumps(data, indent=4, ensure_ascii=False))

    # index_name = data.get('index_name') 之前拼接好的，弃用
    user_id = data.get('user_id')
    user_id = user_id.replace('-', '_')
    index_name = SNIPPET_INDEX_NAME_PREFIX + user_id
    kb_name = data.get('kb_name')
    try:
        # ========= 往里面传入的 kb_name是真正指代的 kb_id =======
        kb_id = es_ops.get_uk_kb_id(data.get('user_id'), data.get('kb_name'))
        status = es_ops.delete_data_by_kbname(index_name, kb_id)
        response = json.dumps({'code': 200, 'msg': 'Success', 'result': status}, indent=4, ensure_ascii=False)
        logger.info("delete_index response: %s", response)
        return Response(response, mimetype='application/json', status=200)
    except Exception as e:
        response = json.dumps({'code': 400, 'msg': str(e), 'result': None}, ensure_ascii=False)
        logger.info("delete_index response: %s", response)
        return Response(response, mimetype='application/json', status=400)
    finally:
        logger.info("request: /api/v1/rag/es/delete_index end")


# ********************* 重启服务后，检查uk映射表索引的mappping，进行一些处理 *********************
# es_ops.check_status() # 不使用此更新。

if __name__ == '__main__':
    app.run()  # debug=True
