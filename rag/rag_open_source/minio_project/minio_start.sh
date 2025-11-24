#!/bin/bash  
  
# 定义日志文件的路径，这里使用当前目录下的flask_app.log  
LOGFILE="minio.log"  
  
# 使用nohup在后台运行Flask应用，并将输出重定向到日志文件  
# 注意：将下面的/path/to/your/app.py替换为你的Flask应用脚本的实际路径  
source activate
conda activate rag-new
nohup python3 /model_extend/minio_project/minio_new.py > "$LOGFILE" 2>&1 &  
  
echo "Flask app started in the background."
