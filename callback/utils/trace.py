import json
import time

from flask import Flask, g, request

from utils.log import logger


def register_tracing(app: Flask):
    """封装路由追踪"""

    @app.before_request
    def start_trace():
        # 保存请求开始时间
        g.start_time = time.time()

        # 尝试获取请求体
        try:
            if request.is_json:
                req_body = request.get_json(silent=True)
            elif request.form:
                req_body = request.form.to_dict()
            elif request.data:
                req_body = request.get_data(as_text=True)
            else:
                req_body = None
        except Exception:
            req_body = "<无法解析请求体>"

        # 记录请求基本信息（此时不记录完整日志，等响应后再统一输出）
        g.request_log = {
            "method": request.method,
            "full_path": request.full_path,
            "header": dict(request.headers),
            "body": req_body,
        }

    @app.after_request
    def end_trace(response):
        request = g.get("request_log", {})
        if "/apidocs" in request.get("full_path", ""):
            return response  # 跳过 apidocs 的日志记录

        # 耗时ms
        cost = round((time.time() - g.get("start_time", time.time())) * 1000, 2)

        # 获取原始响应体（仅适用于非流式响应）
        try:
            if response.is_streamed:
                resp_body = "<流式响应，暂无记录>"
            else:
                resp_body = response.get_data(as_text=True)
        except Exception:
            resp_body = "<无法读取响应体>"

        log_msg = f"{cost}ms | {response.status_code} | {request["method"]} | {request["full_path"]} | {json.dumps(request["body"], ensure_ascii=False)} | {resp_body.rstrip('\n')}"
        if response.status_code < 400:
            logger.info(log_msg)
        else:
            logger.error(log_msg)
        return response
