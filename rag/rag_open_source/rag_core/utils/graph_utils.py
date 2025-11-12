import os
import json
import time

import requests
import pandas as pd
from datetime import datetime

from utils import milvus_utils
from utils import redis_utils
from utils import timing

from logging_config import setup_logging
from settings import GRAPH_SERVER_URL

logger_name = 'rag_graph_utils'
app_name = os.getenv("LOG_FILE")
logger = setup_logging(app_name, logger_name)
logger.info(logger_name + '---------LOG_FILE：' + repr(app_name))


def parse_excel_to_schema_json(file_path):
    """
    解析 Excel 文件中的 '类目表' 和 '类目属性表'，输出指定 JSON 结构
    """
    schema = {}
    try:
        # 使用 pd.read_excel 自动推断引擎（支持 .xls 和 .xlsx）
        df_category = pd.read_excel(file_path, sheet_name='类目表')
        df_attribute = pd.read_excel(file_path, sheet_name='类目属性表')
        # 清理列名：去除空格和换行
        df_category.columns = df_category.columns.str.strip()
        df_attribute.columns = df_attribute.columns.str.strip()

        # === 解析 类目表 ===
        category_list = []
        for _, row in df_category.iterrows():
            item = {
                "类名": str(row["类名"]).strip() if pd.notna(row["类名"]) else "",
                "类描述": str(row["类描述"]).strip() if pd.notna(row["类描述"]) else ""
            }
            category_list.append(item)

        # === 解析 类目属性表 ===
        attribute_list = []

        for _, row in df_attribute.iterrows():
            class_name = str(row["类名"]).strip() if pd.notna(row["类名"]) else ""
            attr_name = str(row["属性名"]).strip() if pd.notna(row["属性名"]) else ""

            # 修复说明字段
            key = (class_name, attr_name)

            desc = str(row["属性说明"]).strip() if pd.notna(row["属性说明"]) else ""

            # 处理别名字段（支持多个别名用 | 分隔）
            alias = row["别名(多别名以|隔开)"]
            if pd.isna(alias) or str(alias).strip() == "" or str(alias).lower() == "nan":
                alias_str = ""
            else:
                alias_str = str(alias).strip()

            value_type = str(row["值类型"]).strip() if pd.notna(row["值类型"]) else ""

            attribute_list.append({
                "类名": class_name,
                "属性名": attr_name,
                "属性说明": desc,
                "属性别名(多别名以|隔开)": alias_str,
                "值类型": value_type
            })

        # 构建最终 JSON 结构
        schema = {
            "schema定义": {
                "类目表": category_list,
                "类目属性表": attribute_list
            }
        }
        logger.info("schema:%s" % json.dumps(schema, ensure_ascii=False))
    except Exception as e:
        import traceback
        logger.error(traceback.format_exc())
        logger.error(f"无法读取Excel文件或工作表不存在: {e}")
        raise Exception("parse graph schema 发生异常：" + str(e))
    return schema


@timing.timing_decorator(logger, include_args=False)
def get_extrac_graph_data(user_id, kb_name, chunks, file_name, schema=None):
    """获取知识图谱数据"""
    try:
        start_time = datetime.now()
        headers = {
            "Content-Type": "application/json",
        }
        data = {
            "user_id": user_id,
            "kb_name": kb_name,
            "chunks": chunks,
            "schema": schema,
            "file_name":file_name
        }
        # 将JSON数据转换好格式
        json_data = json.dumps(data)
        extract_graph_url = GRAPH_SERVER_URL + "/extrac_graph_data"
        response = requests.post(extract_graph_url, headers=headers, data=json_data, timeout=600)
        if response.status_code == 200:
            result_data = json.loads(response.text)
            finish_time1 = datetime.now()
            time_difference1 = finish_time1 - start_time
            logger.info(f"extrac_graph_data -{extract_graph_url}: 请求成功 耗时：{time_difference1}")
            return result_data
        else:
            # 如果不是200，则抛出一个自定义异常
            raise Exception(f"{extract_graph_url} 请求失败，错误信息：" + response.text)
    except Exception as e:
        raise Exception("get_extrac_graph_data 发生异常：" + str(e))


@timing.timing_decorator(logger, include_args=False)
def generate_community_reports(user_id, kb_name, file_name="", graph_data=[]):
    """获取知识图谱社区报告"""
    try:
        start_time = datetime.now()
        headers = {
            "Content-Type": "application/json",
        }
        data = {
            "user_id": user_id,
            "graph_data": graph_data,
            "kb_name": kb_name,
            "file_name": file_name
        }
        # 将JSON数据转换好格式
        json_data = json.dumps(data)
        community_url = GRAPH_SERVER_URL + "/generate_community_reports"
        response = requests.post(community_url, headers=headers, data=json_data, timeout=600)
        if response.status_code == 200:
            result_data = json.loads(response.text)
            finish_time1 = datetime.now()
            time_difference1 = finish_time1 - start_time
            logger.info(f"generate_community_reports -{community_url}: 请求成功 耗时：{time_difference1}")
            return result_data
        else:
            # 如果不是200，则抛出一个自定义异常
            raise Exception(f"{community_url} 请求失败，错误信息：" + response.text)
    except Exception as e:
        raise Exception("generate_community_reports 发生异常：" + str(e))


@timing.timing_decorator(logger, include_args=False)
def delete_file_from_graph(user_id, kb_name, file_name):
    """知识图谱删除文件"""
    try:
        start_time = datetime.now()
        headers = {
            "Content-Type": "application/json",
        }
        data = {
            "user_id": user_id,
            "kb_name": kb_name,
            "file_name": file_name
        }
        # 将JSON数据转换好格式
        json_data = json.dumps(data)
        delete_file_url = GRAPH_SERVER_URL + "/delete_file"
        response = requests.post(delete_file_url, headers=headers, data=json_data, timeout=600)
        if response.status_code == 200:
            result_data = json.loads(response.text)
            finish_time1 = datetime.now()
            time_difference1 = finish_time1 - start_time
            logger.info(f"graph delete_file -{delete_file_url}: 请求成功 耗时：{time_difference1}")
            return result_data
        else:
            # 如果不是200，则抛出一个自定义异常
            raise Exception(f"{delete_file_url} 请求失败，错误信息：" + response.text)
    except Exception as e:
        raise Exception("graph delete_file 发生异常：" + str(e))


@timing.timing_decorator(logger, include_args=False)
def delete_kb_graph(user_id, kb_name):
    """知识图谱删除"""
    try:
        start_time = datetime.now()
        headers = {
            "Content-Type": "application/json",
        }
        data = {
            "user_id": user_id,
            "kb_name": kb_name,
        }
        # 将JSON数据转换好格式
        json_data = json.dumps(data)
        delete_kb_url = GRAPH_SERVER_URL + "/delete_kb"
        response = requests.post(delete_kb_url, headers=headers, data=json_data, timeout=600)
        if response.status_code == 200:
            result_data = json.loads(response.text)
            finish_time1 = datetime.now()
            time_difference1 = finish_time1 - start_time
            logger.info(f"graph delete_kb_graph -{delete_kb_url}: 请求成功 耗时：{time_difference1}")
            return result_data
        else:
            # 如果不是200，则抛出一个自定义异常
            raise Exception(f"{delete_kb_url} 请求失败，错误信息：" + response.text)
    except Exception as e:
        raise Exception("graph delete_kb 发生异常：" + str(e))


def get_graph_vocabulary_set(kb_ids: list):
    """处理获取知识图谱实体词表数据"""
    graph_redis_client = redis_utils.get_redis_connection()
    kb_graph_vocabulary_list = []
    for kb_id in kb_ids:
        kb_graph_vocabulary_set = redis_utils.query_graph_vocabulary_set(graph_redis_client, kb_id)
        res_vocabulary_list = [v.split("|||schema_type:")[0] for v in kb_graph_vocabulary_set]
        res_vocabulary_type_list = [v.split("|||schema_type:")[1] for v in kb_graph_vocabulary_set]
        kb_graph_vocabulary_list.append((res_vocabulary_list, res_vocabulary_type_list))
    # 处理返回结果
    return kb_graph_vocabulary_list


def get_all_extrac_graph_chunks(user_id, kb_name, file_name, kb_id=""):
    """
    获取用户知识库中对应文件的所有chunk
    """

    try:
        time.sleep(10)  # 先等待 10s
        retry_num = 0
        all_chunks = []
        chunk_total_num = 1
        page_size = 100
        search_after = 0
        complete_flag = True
        while len(all_chunks) < chunk_total_num and complete_flag:
            response_info = milvus_utils.get_milvus_file_content_list(user_id, kb_name, file_name, page_size,
                                                                      search_after, kb_id=kb_id)
            temp_content_list = response_info["data"]["content_list"]
            if not temp_content_list:  # 取不到则重试或置完成
                retry_num += 1
                time.sleep(10)
                if retry_num >= 5 or len(all_chunks) >= chunk_total_num:
                    complete_flag = False
            else:
                if chunk_total_num == 1:
                    chunk_total_num = temp_content_list[0]["meta_data"]["chunk_total_num"]
                for doc in temp_content_list:  # 整理格式添加好
                    all_chunks.append({
                        "title": doc["file_name"],
                        "snippet": doc["content"],
                        "source_type": "RAG_KB",
                        "meta_data": doc["meta_data"]
                    })
                search_after += 100
        # ======== 取完了直接返回 =========
        return all_chunks
    except Exception as e:
        logger.info(f"get_all_extrac_graph_chunks Error: {e}")
        return []