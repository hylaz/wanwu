#!/bin/bash

#启动conda环境,在执行脚本前，执行此命令
#conda activate rag-new

# 整个脚本开始时间
TOTAL_START=$(date +%s)

# 更新包列表（可选，但推荐）
apt-get update

# 安装binutils，-y选项表示自动同意
apt-get install -y binutils

# 检查依赖
START_DEP=$(date +%s)
if ! command -v pyinstaller &> /dev/null
then
    pip install pyinstaller
fi

#下载pymilvus和pymilvus.model
pip install pymilvus
pip install pymilvus.model
END_DEP=$(date +%s)
echo "依赖安装耗时: $((END_DEP - START_DEP)) 秒"

# 确保logs目录存在
mkdir -p logs

#移除已打包目录
rm -r ./build/

#config 文件copy
START_CONFIG=$(date +%s)
mkdir -p ./build/dist/rag_core/configs
mkdir -p ./build/dist/rag_es_server_unify/config
cp -r ./rag_core/configs/config.ini ./build/dist/rag_core/configs/
cp -r ./rag_es_server_unify/config/config.ini ./build/dist/rag_es_server_unify/config/
END_CONFIG=$(date +%s)
echo "配置文件复制耗时: $((END_CONFIG - START_CONFIG)) 秒"

#-------- run_app ---------
# 使用PyInstaller打包应用
echo "正在打包run应用..."
START_RUN=$(date +%s)
pyinstaller --name run_app \
            --distpath=./build/dist/rag_core \
            --onefile \
            --add-data "./rag_core/configs:configs" \
            --add-data "../root/miniconda3/envs/rag-new/lib/python3.10/site-packages/pymilvus/model/sparse/bm25/lang.yaml:pymilvus/model/sparse/bm25" \
            --hidden-import="gunicorn" \
            --hidden-import="gunicorn.glogging" \
            --hidden-import="gunicorn.app" \
            --hidden-import="gunicorn.app.base" \
            --hidden-import="gunicorn.app.wsgiapp" \
            --hidden-import="gunicorn.config" \
            --hidden-import="gunicorn.workers" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.ggevent" \
            --hidden-import="gunicorn.workers.gthread" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.geventlet" \
            --hidden-import="pymilvus" \
            --hidden-import="pymilvus.model" \
            --hidden-import "pymilvus.model.sparse.bm25.tokenizers" \
            --hidden-import="langchain_community" \
            --hidden-import="langchain_community.document_loaders" \
            --hidden-import="langchain_community.document_loaders.text" \
            --hidden-import="langchain_community.document_loaders.unstructured" \
            --hidden-import="langchain_community.document_loaders.csv_loader" \
            --hidden-import="tiktoken_ext.openai_public" \
            --hidden-import="tiktoken_ext" \
            ./rag_core/run_entrypoint.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_RUN=$(date +%s)
    echo "run_app打包成功！耗时: $((END_RUN - START_RUN)) 秒"
    echo "可执行文件位于 ./build/dist/rag_core"
else
    echo "run_app 打包失败，请检查错误信息"
    exit 1
fi

#-------- run_sse_app ---------

# 使用PyInstaller打包应用
echo "正在打包run_sse_app应用..."
START_SSE=$(date +%s)
pyinstaller --name sse_app \
            --distpath=./build/dist/rag_core \
            --onefile \
            --add-data "./rag_core/configs:configs" \
            --add-data "../root/miniconda3/envs/rag-new/lib/python3.10/site-packages/pymilvus/model/sparse/bm25/lang.yaml:pymilvus/model/sparse/bm25" \
            --hidden-import="uvicorn.logging" \
            --hidden-import="uvicorn.loops" \
            --hidden-import="uvicorn.loops.auto" \
            --hidden-import="uvicorn.protocols" \
            --hidden-import="uvicorn.protocols.http" \
            --hidden-import="uvicorn.protocols.http.auto" \
            --hidden-import="uvicorn.protocols.websockets" \
            --hidden-import="uvicorn.protocols.websockets.auto" \
            --hidden-import="uvicorn.workers" \
            --hidden-import="uvicorn.loops.auto" \
            --hidden-import="uvicorn.protocols.http.auto" \
            --hidden-import="uvicorn.lifespan" \
            --hidden-import="uvicorn.lifespan.on" \
            --hidden-import="pymilvus" \
            --hidden-import="pymilvus.model" \
            --hidden-import "pymilvus.model.sparse.bm25.tokenizers" \
            --hidden-import="langchain_community" \
            --hidden-import="langchain_community.document_loaders" \
            --hidden-import="langchain_community.document_loaders.text" \
            --hidden-import="langchain_community.document_loaders.unstructured" \
            --hidden-import="langchain_community.document_loaders.csv_loader" \
            --hidden-import="tiktoken_ext.openai_public" \
            --hidden-import="tiktoken_ext" \
            --hidden-import="gunicorn" \
            --hidden-import="gunicorn.glogging" \
            --hidden-import="gunicorn.app" \
            --hidden-import="gunicorn.app.base" \
            --hidden-import="gunicorn.app.wsgiapp" \
            --hidden-import="gunicorn.config" \
            --hidden-import="gunicorn.workers" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.ggevent" \
            --hidden-import="gunicorn.workers.gthread" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.geventlet" \
            ./rag_core/sse_entrypoint.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_SSE=$(date +%s)
    echo "run_sse_app打包成功！耗时: $((END_SSE - START_SSE)) 秒"
    echo "可执行文件位于./build/dist/rag_core"
else
    echo "run_sse_app打包失败，请检查错误信息"
    exit 1
fi

#-------- async_doc_status_init ---------
echo "正在打包初始化脚本 asyn_doc_status_init..."
START_INIT=$(date +%s)
pyinstaller --name asyn_doc_status_init \
            --distpath=./build/dist/rag_core \
            --onefile \
            --add-data "./rag_core/configs/config.ini:configs" \
            ./rag_core/asyn_doc_status_init.py
# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_INIT=$(date +%s)
    echo "async_doc_status_init打包成功！耗时: $((END_INIT - START_INIT)) 秒"
    echo "可执行文件位于 ./build/dist/rag_core 目录"
else
    echo "async_doc_status_init打包失败，请检查错误信息"
    exit 1
fi
#-------- async_add_file ---------

echo "正在打包主应用 asyn_add_file..."
START_ADD=$(date +%s)
pyinstaller --name asyn_add_file \
            --distpath=./build/dist/rag_core \
            --onefile \
            --add-data "./rag_core/configs/config.ini:configs" \
            --add-data "./rag_core/utils:utils" \
            --add-data "./rag_core/logging_config.py:." \
            --add-data "../root/miniconda3/envs/rag-new/lib/python3.10/site-packages/pymilvus/model/sparse/bm25/lang.yaml:pymilvus/model/sparse/bm25" \
            --hidden-import="nltk" \
            --hidden-import="utils.milvus_utils" \
            --hidden-import="utils.minio_utils" \
            --hidden-import="utils.es_utils" \
            --hidden-import="utils.file_utils" \
            --hidden-import="utils.mq_rel_utils" \
            --hidden-import="utils.knowledge_base_utils" \
            --hidden-import="pymilvus" \
            --hidden-import="pymilvus.model" \
            --hidden-import "pymilvus.model.sparse.bm25.tokenizers" \
            --hidden-import="langchain_community" \
            --hidden-import="langchain_community.document_loaders" \
            --hidden-import="langchain_community.document_loaders.text" \
            --hidden-import="langchain_community.document_loaders.unstructured" \
            --hidden-import="langchain_community.document_loaders.csv_loader" \
            --hidden-import="tiktoken_ext.openai_public" \
            --hidden-import="tiktoken_ext" \
            ./rag_core/asyn_add_file.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_ADD=$(date +%s)
    echo "async_add_file打包成功！耗时: $((END_ADD - START_ADD)) 秒"
    echo "可执行文件位于 ./build/dist/rag_core 目录"
else
    echo "async_add_file打包失败，请检查错误信息"
    exit 1
fi

#-------- guarding ---------
echo "正在打包guarding守护进程监控工具..."
START_GUARD=$(date +%s)
pyinstaller --name guarding_asyn_add_app \
           --distpath=./build/dist/rag_core \
           --add-data "./rag_core/asyn_add_file.sh:." \
           --onefile \
           ./rag_core/guarding_file_asyn_add_process.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_GUARD=$(date +%s)
    echo "guarding打包成功！耗时: $((END_GUARD - START_GUARD)) 秒"
    echo "可执行文件位于 ./build/dist/rag_core 目录"
else
   echo "guarding打包失败，请检查错误信息"
   exit 1
fi


#-------- url single ---------
# 使用PyInstaller打包应用
echo "正在打包url single应用..."
START_URL=$(date +%s)
pyinstaller --name url_single_app \
            --distpath=./build/dist/rag_core/url_parser \
            --onefile \
            --add-data "./rag_core/configs:configs" \
            --hidden-import="gunicorn" \
            --hidden-import="gunicorn.glogging" \
            --hidden-import="gunicorn.app" \
            --hidden-import="gunicorn.app.base" \
            --hidden-import="gunicorn.app.wsgiapp" \
            --hidden-import="gunicorn.config" \
            --hidden-import="gunicorn.workers" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.ggevent" \
            --hidden-import="gunicorn.workers.gthread" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.geventlet" \
            ./rag_core/url_entrypoint.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_URL=$(date +%s)
    echo "url single打包成功！耗时: $((END_URL - START_URL)) 秒"
    echo "可执行文件位于 ./build/dist/rag_core/url_parser"
else
    echo "url single打包失败，请检查错误信息"
    exit 1
fi


#-------- minio ---------
# 使用PyInstaller打包应用
echo "正在打包minio应用..."
START_MINIO=$(date +%s)
pyinstaller --name minio_app \
            --distpath=./build/dist/minio_project \
            --onefile \
            --hidden-import="gunicorn" \
            --hidden-import="gunicorn.glogging" \
            --hidden-import="gunicorn.app" \
            --hidden-import="gunicorn.app.base" \
            --hidden-import="gunicorn.app.wsgiapp" \
            --hidden-import="gunicorn.config" \
            --hidden-import="gunicorn.workers" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.ggevent" \
            --hidden-import="gunicorn.workers.gthread" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.geventlet" \
            ./minio_project/minio_entrypoint.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_MINIO=$(date +%s)
    echo "minio 打包成功！耗时: $((END_MINIO - START_MINIO)) 秒"
    echo "可执行文件位于 ./build/dist/minio_project"
else
    echo "打包失败，请检查错误信息"
    exit 1
fi

#-------- es ---------
# 使用PyInstaller打包应用
echo "正在打包es应用..."
START_ES=$(date +%s)
pyinstaller --name es_app \
            --onefile \
            --distpath=./build/dist/rag_es_server_unify \
            --add-data "./rag_es_server_unify/config/config.ini:config" \
            --hidden-import="gunicorn" \
            --hidden-import="gunicorn.glogging" \
            --hidden-import="gunicorn.app" \
            --hidden-import="gunicorn.app.base" \
            --hidden-import="gunicorn.app.wsgiapp" \
            --hidden-import="gunicorn.config" \
            --hidden-import="gunicorn.workers" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.ggevent" \
            --hidden-import="gunicorn.workers.gthread" \
            --hidden-import="gunicorn.workers.sync" \
            --hidden-import="gunicorn.workers.geventlet" \
            ./rag_es_server_unify/es_entrypoint.py

# 检查打包是否成功
if [ $? -eq 0 ]; then
    END_ES=$(date +%s)
    echo "es打包成功！耗时: $((END_ES - START_ES)) 秒"
    echo "可执行文件位于 ./build/dist/rag_es_server_unify"
else
    echo "es打包失败，请检查错误信息"
    exit 1
fi

echo "清理打包过程文件"
cd ./build
# 启用extglob扩展（支持!匹配）
shopt -s extglob

# 删除除dist外的所有文件和目录
echo "正在清理..."
rm -rf !(dist)

# 关闭extglob
shopt -u extglob
cd ../

echo "清理完成！"

# 计算总耗时
TOTAL_END=$(date +%s)
TOTAL_TIME=$((TOTAL_END - TOTAL_START))
echo "=============================================="
echo "所有打包任务完成！"
echo "总耗时: $TOTAL_TIME 秒"
echo "=============================================="