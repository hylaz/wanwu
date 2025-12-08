from flask import request

from callback.services import hello as hello_service
from utils.log import logger
from utils.response import BizError, response_err, response_ok

from . import callback_bp


@callback_bp.route("/hello", methods=["GET"])
def get_hello():
    """
    API 示例
    ---
    tags:
      - hello
    parameters:
      - name: username
        description: 用户名
        in: query
        required: true
        schema:
          type: string
    responses:
      200:
        description: 返回信息
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 状态码
                  example: 0
                msg:
                  type: string
                  description: 描述信息
                  example: "success"
                data:
                  type: string
                  example: Hello, {username}!
      400:
        description: 业务逻辑错误
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 错误码
                  example: 200000
                msg:
                  type: string
                  description: 错误描述信息
      500:
        description: 服务内部错误
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 错误码
                  example: 200000
                msg:
                  type: string
                  description: 错误描述信息
    """
    username = request.args.get("username")
    if not username:
        raise BizError("username is None")
    return response_ok(hello_service.get_message(username))


@callback_bp.route("/hello", methods=["POST"])
def post_hello():
    """
    API 示例
    ---
    tags:
      - hello
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              username:
                type: string
                description: 用户名
            required:
              - username
    responses:
      200:
        description: 返回信息
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 状态码
                  example: 0
                msg:
                  type: string
                  description: 描述信息
                  example: "success"
                data:
                  type: string
                  example: Hello, {username}!
      400:
        description: 业务逻辑错误
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 错误码
                  example: 200000
                msg:
                  type: string
                  description: 错误描述信息
      500:
        description: 服务内部错误
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 错误码
                  example: 200000
                msg:
                  type: string
                  description: 错误描述信息
    """
    data = request.get_json()
    if data is None:
        raise BizError("request body is None")
    username = data["username"]
    if not username:
        raise BizError("username is None")
    return response_ok(hello_service.get_message(username))


@callback_bp.route("/hello", methods=["PUT"])
def put_hello():
    """
    API 示例
    ---
    tags:
      - hello
    requestBody:
      content:
        multipart/form-data:
          schema:
            type: object
            properties:
              username:
                type: string
                description: 用户名
            required:
              - username
    responses:
      200:
        description: 返回信息
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 状态码
                  example: 0
                msg:
                  type: string
                  description: 描述信息
                  example: "success"
                data:
                  type: string
                  example: Hello, {username}!
      400:
        description: 业务逻辑错误
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 错误码
                  example: 200000
                msg:
                  type: string
                  description: 错误描述信息
      500:
        description: 服务内部错误
        content:
          application/json:
            schema:
              type: object
              properties:
                code:
                  type: integer
                  description: 错误码
                  example: 200000
                msg:
                  type: string
                  description: 错误描述信息
    """
    username = request.form["username"]
    if not username:
        raise BizError("username is None")
    return response_ok(hello_service.get_message(username))
