import logging
import posixpath
from urllib.parse import urljoin

from flask import jsonify, request

from callback.services import minio as minio_service
from configs.config import config

from . import callback_bp


@callback_bp.route("/upload", methods=["POST"])
def upload_file():
    """
    上传文件到 MinIO 存储
    ---
    tags:
      - minio
    requestBody:
      required: true
      content:
        multipart/form-data:
          schema:
            type: object
            required:
              - file
            properties:
              file:
                format: binary
                description: 需要上传的文件
              bucket_name:
                type: string
                description: 目标存储桶名称
              file_name:
                type: string
                description: (可选) 重写保存的文件名
    responses:
      200:
        description: 上传成功
        content:
          application/json:
            schema:
              type: object
              properties:
                download_link:
                  type: string
                  description: 文件的下载链接
                  example: "http://base-url/callback/filename.jpg"
      500:
        description: 服务器内部错误
        content:
          application/json:
            schema:
              type: object
              properties:
                error:
                  type: string
                  example: "Failed to upload file to Minio."
    """
    try:
        uploaded_file = request.files["file"]
        original_filename = uploaded_file.filename
        overwrite_filename = request.form.get("file_name", None)
        bucket_name = request.form.get("bucket_name", None)

        object_path = minio_service.upload_file_to_minio(
            uploaded_file, original_filename, overwrite_filename, bucket_name
        )
        if object_path:
            download_link = posixpath.join(
                config.callback_cfg["URL"]["MINIO_DOWNLOAD"], object_path
            )
            return jsonify({"download_link": download_link})
        else:
            return jsonify({"error": "Failed to upload file to Minio."}), 500
    except Exception as e:
        return jsonify({"error": str(e)}), 500
