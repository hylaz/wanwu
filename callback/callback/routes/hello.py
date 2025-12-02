from flask import Blueprint, request

from callback.services.hello_demo import get_message
from utils.log import logger
from utils.response import BizError, response_err, response_ok

hello_demo = Blueprint("hello", __name__)


@hello_demo.route("/hello", methods=["GET"])
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
        schema:
          type: object
          properties:
            msg:
              type: string
              example: Hello, {username}!
    """
    username = request.args.get("username")
    if not username:
        raise BizError("username is None")
    return response_ok(get_message(username))


@hello_demo.route("/hello", methods=["POST"])
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
        schema:
          type: object
          properties:
            msg:
              type: string
              example: Hello, {username}!
    """
    data = request.get_json()
    if data is None:
        raise BizError("request body is None")
    username = data["username"]
    if not username:
        raise BizError("username is None")
    return response_ok(get_message(username))


@hello_demo.route("/hello", methods=["PUT"])
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
        schema:
          type: object
          properties:
            msg:
              type: string
              example: Hello, {username}!
    """
    username = request.form["username"]
    if not username:
        raise BizError("username is None")
    return response_ok(get_message(username))
