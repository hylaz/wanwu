
### RAG-README.md

gitlab地址：
https://gitlab.ai-yuanjing.cn/model_extend_platform/langchain_rag

###  langchain_rag 项目简介


**特别版本说明：固定RAG2.1分支为通用版本-部署应用于云端MAAS,N卡和华为910B私有化部署，3090一体机部署**

### RAG2.1分支 部署N卡和华为910B私有化使用docker，3090一体机裸机部署，以下是一些说明

#### N卡和华为910B私有化部署使用方法（依赖docker）
使用docker,固定docker_name:model_ext

进入镜像：
```
sudo docker exec -it --user root model_ext /bin/bash 
sh /model_extend/start.sh

```

说明，私有化部署启动容器后，会调用start.sh 进行服务启动


#### 私有化部署RAG服务及使用到的端口说明

> http://172.17.0.1:10891 RAG-问答流式输出接口-端口
> 
> http://172.17.0.1:8681 RAG-主控接口-端口
> 
> http://172.17.0.1:20041 RAG-ES数据库接口-端口
> 
> http://172.17.0.1:49021 RAG-BGE embedding接口-端口  
> 
> http://172.17.0.1:49031 RAG-BCE rerank接口-端口 
> 
> http://172.17.0.1:15000 RAG-minio上传接口-端口 



#### 3090一体机部署使用方法（裸机部署）

```
sh /model_extend/start.sh
```
说明，3090一体机，会调用start.sh 进行服务启动

#### 3090一体机RAG服务及使用到的端口说明

> http://172.17.0.1:10891 RAG-问答流式输出接口-端口
> 
> http://172.17.0.1:8681 RAG-主控接口-端口
> 
> http://172.17.0.1:20041 RAG-ES数据库接口-端口
> 
> http://172.17.0.1:49021 RAG-BGE embedding, rerank接口-端口
> 
> http://172.17.0.1:15000 RAG-minio上传接口-端口



## 版本更新日志（从下往上更新，最上方是最新的版本及更新操作说明）


### 2 更新发布 master，主分支新增知识飞轮升级及知识库改名接口服务 2024-11-11
知识飞轮升级，只适用于MAAS云端版本。私有化版本不适用

- RAG 主控及知识库功能更新
  - 1.【feature】: 知识库文件删除接口和知识库删除接口的执行流程增加知识飞轮数据删除流程
  - 2.【feature】: 知识库新增知识库改名接口，新增查询知识库名对应 kb_id 的接口

- RAG 问答流式接口接口功能更新
  - 2.【feature】:知识增强问答流式增加可选参数 `data_flywheel`:是否启用数据飞轮机制回答

### 1 更新发布 RAG-2.1, 分支新增功能及版本进行适配商飞等升级 2024-11-11
该分支为通用固定功能版本-可部署应用于云端MAAS,N卡和华为910B私有化部署，3090一体机部署
- 代码架构调整更新
  - 1.【feature】: 代码结构调整，以配置文件 `/configs/config.ini` 为不同部署环境配置。
  - 2.【bugfix】:LLM上下文长度限制问题，配置文件增加配置参数`CONTEXT_LENGTH`,私有化部署默认为6024
  
- RAG 主控及知识库功能更新
  - 1.【feature】:支持私有化部署的前端平台使用显示文件切分的知识片段功能。
  - 2.【feature】:私有化部署的版本配置文件增加 REPLACE_MINIO_IP_API,可从这接口获取minio下载文件服务的IP地址和端口


- RAG 文件解析功能更新
  - 1.【feature】:excel文件及.doc文件增加 解析知识片段的 meta_data 里增加 download_link。
  - 2.【bugfix】:修改kafka消费者，按顺序消费消息(处理完当前消息后，再获取下一条消息),当前设置4个分片，4个消费者


- RAG 问答流式接口接口功能更新
  - 1.【feature】:知识增强问答流式增加可选参数 `auto_citation`:是否开启自动引文
  - 2.【feature】:知识增强问答流式增加可选参数 `return_score_`:是否返回search_list的相关性得分
  - 3.【feature】:知识增强问答流式支持图文混合输出(以markdown格式输出)
  - 4.【feature】:当有配置REPLACE_MINIO_IP_API时，返回的minio下载链接自动替换IP地址和端口





