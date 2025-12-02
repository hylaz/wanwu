import time

from flask import Flask, jsonify, make_response

general_err_code = 200000


class BizError(Exception):
    """业务异常，默认返回 code=200000"""

    def __init__(self, msg, code=general_err_code):
        self.code = code
        self.msg = msg
        super().__init__(msg)


def response_ok(data=None, code=0, msg="success"):
    return make_response(jsonify({"code": code, "msg": msg, "data": data or {}}), 200)


def response_err(err: Exception, http_status=400, data=None):
    if isinstance(err, BizError):
        code = err.code
        msg = err.msg
    else:
        http_status = http_status or 500
        code = general_err_code
        msg = str(err)
    return make_response(
        jsonify({"code": code, "msg": msg, "data": data or {}}), http_status
    )


def response_err_msg(code, msg, http_status=400, data=None):
    return make_response(
        jsonify({"code": code, "msg": msg, "data": data or {}}), http_status
    )


def register_error_handlers(app: Flask):

    @app.errorhandler(BizError)
    def handle_biz_error(e):
        return response_err_msg(code=e.code, msg=e.msg, http_status=400)

    @app.errorhandler(Exception)
    def handle_exception(e):
        return response_err_msg(code=general_err_code, msg=str(e), http_status=500)
