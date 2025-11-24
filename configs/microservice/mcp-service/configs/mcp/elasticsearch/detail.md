
# Elasticsearch 工具集

本项目包含三个独立的Elasticsearch工具服务，分别提供不同的功能集。

## 1. ES_Calc_Tools (端口: 9031)

基于ES计算的MCP工具，支持数据查询、统计和导出功能。

### 主要功能

- **ES_Analytics类**:
  - 连接Elasticsearch集群
  - 支持多种查询类型(match/gte/range)
  - 数据导出到Excel/CSV
  - 聚合统计功能(sum/avg/max/min/count)

### 工具接口

```python
@mcp.tool()
def mcp_search(index='local_test', properties='is_activate', query_value=True, query_type=''):
    """执行ES查询"""

@mcp.tool()
def mcp_q2w(index_name='local_test', query={}, output_file="report.xlsx"):
    """导出查询结果到Excel"""

@mcp.tool()
def mcp_q2w(index="local_test", filter_conditions={"is_active": True},
           field_name="salary", operation="sum"):
    """执行聚合统计"""
```

## 2. ES_file_Tools (端口: 9032)

基于ES的文件管理工具集，支持文件数据导入到Elasticsearch。

### 主要功能

- **ES_DataImporter类**:
  - 检查ES连接状态
  - 创建索引(支持自定义mapping)
  - 自动转换数据类型(日期/布尔值等)
  - 分块批量导入(CSV/Excel)
  - 日志记录

### 工具接口

```python
@mcp.tool()
def create_index(index, SAMPLE_MAPPING, hosts='http://127.0.0.1:9200',
               ES_USER='elastic', ES_PASSWORD='infini_rag_flow'):
    """新建索引"""

@mcp.tool()
def load_file(index, file_path, id_field, hosts='http://127.0.0.1:9200',
             ES_USER='elastic', ES_PASSWORD='infini_rag_flow', chunk_size=1000):
    """导入文件数据到ES"""
```

## 3. ES_Control_Tools (端口: 9033)

基于ES的管理工具集，提供增删改查等基础操作。

### 主要功能

- **Text_Database类**:
  - 文档CRUD操作(增删改查)
  - 批量操作
  - 全文搜索(支持高亮)
  - 语义搜索(向量相似度)
  - 数据导入导出(CSV)
  - 索引管理

### 工具接口

```python
@mcp.tool()
def mcp_add_document(index: str, document: Dict[str, Any]) -> Dict[str, Any]:
    """添加单个文档"""

@mcp.tool()
def mcp_update(doc_id, update_fields, index="my_index"):
    """更新文档"""

@mcp.tool()
def mcp_search(query_vector, top_k: int = 5):
    """语义搜索"""

@mcp.tool()
def mcp_del(doc_id):
    """删除文档"""
```

## 环境配置

所有服务都支持以下环境变量配置:
- `ES_HOST`: Elasticsearch主机地址(默认: 10.103.66.86)
- `ES_PORT`: Elasticsearch端口(默认: 9200)
- `ES_USER`: ES用户名(默认: elastic)
- `ES_PASSWORD`: ES密码(默认: infini_rag_flow)

## 启动方式

每个服务都可以通过以下方式启动:
```bash
python <脚本名>.py --hosts <ES地址> --index <索引名> [其他参数]
```

默认端口:
- ES_Calc_Tools: 9031
- ES_file_Tools: 9032
- ES_Control_Tools: 9033
