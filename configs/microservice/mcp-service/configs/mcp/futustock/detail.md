# Futu Stock MCP 服务器

基于模型上下文协议(Model Context Protocol, MCP)的富途开放API功能访问服务器。

## 特性

- 完全兼容 MCP 2.0 协议标准
- 全面覆盖富途 API 功能
- 支持实时数据订阅
- 市场数据访问
- 衍生品信息查询
- 账户查询功能
- 基于资源的数据访问
- 交互式分析提示

## 环境要求

- Python 3.10+
- 富途开放API SDK
- 模型上下文协议SDK
- OpenD (富途网关程序) - [官方文档](https://openapi.futunn.com/futu-api-doc/intro/intro.html)
  - 运行于本地电脑或云端服务器
  - 负责中转协议请求到富途后台
  - 支持 Windows、MacOS、CentOS、Ubuntu
  - **必须先安装并运行 OpenD 才能使用本服务器**
- uv (推荐)

## 市场支持

本服务器支持以下市场的数据访问（需要相应的行情权限）：

### 香港市场
- 股票、ETFs、窝轮、牛熊、界内证
- 期权
- 期货
- 指数
- 板块

### 美国市场
- 股票、ETFs（纽交所、美交所、纳斯达克）
- 期权
- 期货
- 板块

### A股市场
- 股票、ETFs
- 指数
- 板块

更多市场数据权限说明请参考[富途OpenAPI文档](https://openapi.futunn.com/futu-api-doc/intro/authority.html)

## 安装步骤

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/futu-stock-mcp-server.git
cd futu-stock-mcp-server
```

2. 安装 uv：
```bash
# macOS/Linux
curl -LsSf https://astral.sh/uv/install.sh | sh

# Windows (PowerShell)
powershell -c "irm https://astral.sh/uv/install.ps1 | iex"
```

3. 创建并激活虚拟环境：
```bash
# 创建虚拟环境
uv venv

# 激活虚拟环境
# macOS/Linux:
source .venv/bin/activate
# Windows:
.venv\Scripts\activate
```

4. 安装依赖：
```bash
# 以可编辑模式安装
uv pip install -e .
```

5. 复制并配置环境文件：
```bash
cp .env.example .env
```

编辑 `.env` 文件，设置服务器配置：
```
HOST=0.0.0.0
PORT=8000
FUTU_HOST=127.0.0.1
FUTU_PORT=11111
```

## 开发指南

### 依赖管理

在 `pyproject.toml` 中添加新依赖：
```toml
[project]
dependencies = [
    # ... 现有依赖 ...
    "new-package>=1.0.0",
]
```

然后更新环境：
```bash
uv pip install -e .
```

### 代码风格

本项目使用 Ruff 进行代码检查和格式化。配置在 `pyproject.toml` 中：

```toml
[tool.ruff]
line-length = 100
target-version = "py38"

[tool.ruff.lint]
select = ["E", "F", "I", "N", "W", "B", "UP"]
```

运行代码检查：
```bash
uv pip install ruff
ruff check .
```

运行代码格式化：
```bash
ruff format .
```

## 使用方法

1. 启动服务器：
```bash
python -m futu_stock_mcp_server.server
```

2. 使用 MCP 客户端连接服务器：
```python
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    server_params = StdioServerParameters(
        command="python",
        args=["src/server.py"]
    )
    
    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # 初始化连接
            await session.initialize()
            
            # 列出可用工具
            tools = await session.list_tools()
            
            # 调用工具
            result = await session.call_tool(
                "get_stock_quote",
                arguments={"symbols": ["HK.00700"]}
            )
            
            # 访问资源
            content, mime_type = await session.read_resource(
                "market://HK.00700"
            )
            
            # 获取提示
            prompt = await session.get_prompt(
                "market_analysis",
                arguments={"symbol": "HK.00700"}
            )

if __name__ == "__main__":
    import asyncio
    asyncio.run(main())
```

## 可用 MCP 功能

### 市场数据功能

#### get_stock_quote
获取股票报价数据。
```python
symbols = ["HK.00700", "US.AAPL", "SH.600519"]
result = await session.call_tool("get_stock_quote", {"symbols": symbols})
```
返回价格、成交量、成交额等报价数据。

#### get_market_snapshot
获取市场快照数据。
```python
symbols = ["HK.00700", "US.AAPL", "SH.600519"]
result = await session.call_tool("get_market_snapshot", {"symbols": symbols})
```
返回包括价格、成交量、买卖盘价格等完整市场数据。

#### get_cur_kline
获取当前K线数据。
```python
result = await session.call_tool("get_cur_kline", {
    "symbol": "HK.00700",
    "ktype": "K_1M",  # K_1M, K_5M, K_15M, K_30M, K_60M, K_DAY, K_WEEK, K_MON
    "count": 100
})
```

#### get_history_kline
获取历史K线数据。
```python
result = await session.call_tool("get_history_kline", {
    "symbol": "HK.00700",
    "ktype": "K_DAY",
    "start": "2024-01-01",
    "end": "2024-03-31"
})
```

#### get_rt_data
获取实时交易数据。
```python
result = await session.call_tool("get_rt_data", {"symbol": "HK.00700"})
```

#### get_ticker
获取逐笔成交数据。
```python
result = await session.call_tool("get_ticker", {"symbol": "HK.00700"})
```

#### get_order_book
获取买卖盘数据。
```python
result = await session.call_tool("get_order_book", {"symbol": "HK.00700"})
```

#### get_broker_queue
获取经纪队列数据。
```python
result = await session.call_tool("get_broker_queue", {"symbol": "HK.00700"})
```

### 订阅功能

#### subscribe
订阅实时数据。
```python
result = await session.call_tool("subscribe", {
    "symbols": ["HK.00700", "US.AAPL"],
    "sub_types": ["QUOTE", "TICKER", "K_1M"]
})
```
订阅类型：
- "QUOTE": 基本报价
- "ORDER_BOOK": 买卖盘
- "TICKER": 逐笔成交
- "RT_DATA": 实时数据
- "BROKER": 经纪队列
- "K_1M" 到 "K_MON": K线数据

#### unsubscribe
取消订阅实时数据。
```python
result = await session.call_tool("unsubscribe", {
    "symbols": ["HK.00700", "US.AAPL"],
    "sub_types": ["QUOTE", "TICKER"]
})
```

### 期权功能

#### get_option_chain
获取期权链数据。
```python
result = await session.call_tool("get_option_chain", {
    "symbol": "HK.00700",
    "start": "2024-04-01",
    "end": "2024-06-30"
})
```

#### get_option_expiration_date
获取期权到期日。
```python
result = await session.call_tool("get_option_expiration_date", {
    "symbol": "HK.00700"
})
```

#### get_option_condor
获取期权康多策略数据。
```python
result = await session.call_tool("get_option_condor", {
    "symbol": "HK.00700",
    "expiry": "2024-06-30",
    "strike_price": 350.0
})
```

#### get_option_butterfly
获取期权蝶式策略数据。
```python
result = await session.call_tool("get_option_butterfly", {
    "symbol": "HK.00700",
    "expiry": "2024-06-30",
    "strike_price": 350.0
})
```

### 账户功能

#### get_account_list
获取账户列表。
```python
result = await session.call_tool("get_account_list", {"random_string": "dummy"})
```

#### get_funds
获取账户资金信息。
```python
result = await session.call_tool("get_funds", {"random_string": "dummy"})
```

#### get_positions
获取持仓信息。
```python
result = await session.call_tool("get_positions", {"random_string": "dummy"})
```

#### get_max_power
获取最大交易能力。
```python
result = await session.call_tool("get_max_power", {"random_string": "dummy"})
```

#### get_margin_ratio
获取股票保证金比率。
```python
result = await session.call_tool("get_margin_ratio", {"symbol": "HK.00700"})
```

### 市场信息功能

#### get_market_state
获取市场状态。
```python
result = await session.call_tool("get_market_state", {"market": "HK"})
```
可用市场："HK", "US", "SH", "SZ"

#### get_security_info
获取证券信息。
```python
result = await session.call_tool("get_security_info", {
    "market": "HK",
    "code": "00700"
})
```

#### get_security_list
获取市场证券列表。
```python
result = await session.call_tool("get_security_list", {"market": "HK"})
```

#### get_stock_filter
基于条件筛选股票。
```python
result = await session.call_tool("get_stock_filter", {
    "market": "HK.Motherboard",
    "base_filters": [{
        "field_name": 1,  # 价格
        "filter_min": 10.0,
        "filter_max": 50.0,
        "sort_dir": 1  # 升序
    }],
    "page": 1,
    "page_size": 50
})
```

### 时间功能

#### get_current_time
获取服务器当前时间。
```python
result = await session.call_tool("get_current_time", {"random_string": "dummy"})
```
返回时间戳、格式化日期时间、日期、时间和时区。

## 资源

### 市场数据
- `market://{symbol}`: 获取股票市场数据
- `kline://{symbol}/{ktype}`: 获取K线数据

## 提示功能

### 分析
- `market_analysis`: 创建市场分析提示
- `option_strategy`: 创建期权策略分析提示

## 错误处理

服务器遵循 MCP 2.0 错误响应格式：

```json
{
    "jsonrpc": "2.0",
    "id": "request_id",
    "error": {
        "code": -32000,
        "message": "错误信息",
        "data": null
    }
}
```

## 安全性

- 服务器使用安全的WebSocket连接
- 所有API调用通过富途OpenAPI进行认证
- 使用环境变量存储敏感配置信息

## 开发指南

### 添加新工具

使用 `@mcp.tool()` 装饰器添加新工具：

```python
@mcp.tool()
async def new_tool(param1: str, param2: int) -> Dict[str, Any]:
    """工具描述"""
    # 实现代码
    return result
```

### 添加新资源

使用 `@mcp.resource()` 装饰器添加新资源：

```python
@mcp.resource("resource://{param1}/{param2}")
async def new_resource(param1: str, param2: str) -> Dict[str, Any]:
    """资源描述"""
    # 实现代码
    return result
```

### 添加新提示

使用 `@mcp.prompt()` 装饰器添加新提示：

```python
@mcp.prompt()
async def new_prompt(param1: str) -> str:
    """提示描述"""
    return f"包含{param1}的提示模板"
```

## 许可证

MIT 许可证 