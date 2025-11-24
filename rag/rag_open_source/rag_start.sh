#!/bin/bash

for port in 8681 10891 15000 8613 20041; do
ps -ef | grep $port | grep -v grep | awk '{print $2}' | xargs kill -9
done

#-------- run_app ---------
echo "正在启动run应用..."

# 定义日志文件名称
RUN_BASE_LOG_FILE="run_logs_"
RUN_START_LOG_FILE="run_console.log"
RUN_NUM_WORKERS=5
RUN_TIMEOUT=600

# 循环启动FastAPI应用
for PORT in {8681..8681}; do
    # 启动应用
    echo "正在启动FlaskAPI应用，端口号为$PORT..."
    cd ./build/dist/rag_core/ || exit
    LOG_FILE=$RUN_BASE_LOG_FILE$PORT nohup ./run_app --port $PORT --workers $RUN_NUM_WORKERS --timeout $RUN_TIMEOUT >>$RUN_START_LOG_FILE 2>&1 &
    echo "应用启动成功，日志文件为./build/dist/rag_core/logs/$RUN_BASE_LOG_FILE$PORT.log。"
    cd ../../../
done
#-------- sse_app ---------
echo "正在启动sse应用..."

# 定义日志文件名称
SSE_BASE_LOG_FILE="kb_sse_"
SSE_START_LOG_FILE="kb_sse_console.log"

# 循环启动FastAPI应用
for PORT in {10891..10891}; do
    # 启动应用
    echo "正在启动FastAPI应用，端口号为$PORT..."
    cd ./build/dist/rag_core/ || exit
    LOG_FILE=$SSE_BASE_LOG_FILE$PORT nohup ./sse_app --port $PORT >>$SSE_START_LOG_FILE 2>&1 &
    echo "应用启动成功，日志文件为./logs/$SSE_BASE_LOG_FILE$PORT.log。"
    cd ../../../
done

#-------- asyn_doc_status_init ---------
nohup ./build/dist/rag_core/asyn_doc_status_init > logs/init_asyn.out 2>&1
sleep 1

#-------- async_add_file ---------
#kill现有进程
ps -ef | grep '[a]syn_add_file' | grep -v grep | awk '{print $2}' | xargs kill -9
# 定义日志文件名称
ASYNC_ADD_FILE_BASE_LOG_FILE="asyn_add_"
ASYNC_ADD_FILE_START_LOG_FILE="kb_sse_console.log"

cd ./build/dist/rag_core/ || exit
#循环5次，启动asyn_add
for ADDID in $(seq -f "%03g" 1 2)
do
    LOG_FILE=$ASYNC_ADD_FILE_BASE_LOG_FILE$ADDID nohup ./asyn_add_file >>$ASYNC_ADD_FILE_START_LOG_FILE 2>&1 &
    echo "应用启动成功，日志文件为./logs/$ASYNC_ADD_FILE_BASE_LOG_FILE$ADDID.log。"
done
cd ../../../

#-------- guarding ---------
#kill现有进程
ps -ef | grep '[g]uarding_asyn_add_app' | grep -v grep | awk '{print $2}' | xargs kill -9
# 定义日志文件名称
nohup ./build/dist/rag_core/guarding_asyn_add_app > logs/guarding_asyn_add_process.out 2>&1 &
echo "守护异步解析进程监控已启动"

#-------- url single ---------
# 定义日志文件名称
URL_SINGLE_BASE_LOG_FILE="url_single_logs_"
URL_SINGLE_START_LOG_FILE="url_single_console.log"

# 循环启动FastAPI应用
for PORT in {8613..8613}; do
    # 启动应用
    echo "正在启动url single应用，端口号为$PORT..."
    cd ./build/dist/rag_core/url_parser || exit
    LOG_FILE=$URL_SINGLE_BASE_LOG_FILE$PORT nohup ./url_single_app --port $PORT >>$URL_SINGLE_START_LOG_FILE 2>&1 &
    echo "应用url single启动成功，日志文件为./logs/$URL_SINGLE_BASE_LOG_FILE$PORT.log。"
    cd ../../../../
done

#-------- minio ---------
# 定义日志文件名称
MINIO_BASE_LOG_FILE="minio_logs_"
MINIO_START_LOG_FILE="minio_console.log"

# 循环启动FastAPI应用
for PORT in {15000..15000}; do
    # 启动应用
    echo "正在启动minio FlaskAPI应用，端口号为$PORT..."
    cd ./build/dist/minio_project || exit
    LOG_FILE=$MINIO_BASE_LOG_FILE$PORT nohup ./minio_app --port $PORT >>$MINIO_START_LOG_FILE 2>&1 &
    echo "应用启动成功，日志文件为./logs/$MINIO_BASE_LOG_FILE$PORT.log。"
    cd ../../../
done

#-------- es server ---------
ES_NUM_WORKERS=5
# 定义日志文件名称
ES_BASE_LOG_FILE="es_logs_"
ES_START_LOG_FILE="es_console.log"

# 循环启动FastAPI应用
for PORT in {20041..20041}; do
    # 启动应用
    echo "启动 es_app 程序...到${PORT}端口"
    cd ./build/dist/rag_es_server_unify/ || exit
    LOG_FILE=$ES_BASE_LOG_FILE$PORT nohup ./es_app --port $PORT --workers $ES_NUM_WORKERS >>$ES_START_LOG_FILE 2>&1 &
    echo "应用启动成功，日志文件为./logs/$ES_BASE_LOG_FILE$PORT.log。"
    cd ../../../
done

#-------- 阻塞启动防止docker 自动重启 ---------
sleep 99999d
