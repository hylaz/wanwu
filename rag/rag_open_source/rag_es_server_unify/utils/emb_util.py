import json
import time
import numpy as np

from openai import OpenAI
from model.model_manager import get_model_configure

from log.logger import logger

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
    rate_limit_backoff = [10, 20, 40, 60]  # 限流退避
    other_error_max_retries = 2  # 其他错误最多重试2次
    other_error_wait = 0.5  # 每次0.5s

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


def calculate_cosine(query, contents, embedding_model_id="") -> list[float]:
    query_vector_scores = []
    query_vector = get_embs([query], embedding_model_id=embedding_model_id)["result"][0]["dense_vec"]
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

