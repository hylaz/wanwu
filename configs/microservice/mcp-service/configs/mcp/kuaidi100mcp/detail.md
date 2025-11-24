快递100 MCP Server是快递100推出的国内首个兼容MCP协议的物流信息服务平台。快递100旗下百递云·API开放平台的核心服务现已全面支持MCP协议。开发者通过简单配置即可快速接入快递查询、运费预估、智能时效预估（含全程与在途模式）等核心功能。其AI Agent不仅显著降低了开发过程中物流数据服务调用的门槛，提高了开发效率，还增强了对各大行业的物流数据赋能，助力其创新与发展。

核心功能
----

1.   **快递查询**
    *   快递100API开放平台，提供查询全球3000+物流公司轨迹查询的能力。
    *   输入：物流单号
    *   输出：物流轨迹信息包含物流时间节点和轨迹详情

2.   **快递价格预估**
    *   通过快递100大数据分析，根据快递公司、收寄件地址和重量来预估快递公司运费价格。
    *   输入：收件地址、寄件地址、快递公司名称、包裹重量
    *   输出：预估快递寄件价格

3.   **智能时效预估（全程模式）**
    *   发货前，根据收寄件地址预测不同快递公司的预计送达时间。
    *   输入：快递公司编码、收件地址、寄件地址
    *   输出：预计送达时间

4.   **智能时效预估（在途模式）**
    *   利用快递100自研AI时效预估模型，预测快递的送达时间。
    *   输入：下单时间、物流轨迹信息、收件地址、寄件地址
    *   输出：预计送达时间

快递100 MCP Server的关键特性
---------------------

1.   **快递物流数据质量国内第一**
    *   数据源自百递云开放平台，支持全球3000+快递公司物流信息数据查询，专注行业15年，有250万+企业客户选择。

2.   **AI+Data+MCP API重新定义**
    *   AI Agent显著降低了开发过程中物流数据服务调用的门槛，提高了开发效率，增强了对各大行业的物流数据赋能，助力其创新与发展。快递100 MCP Server支持通过SSE方式接入任意支持MCP协议的平台。

3.   **AI+Data快递图谱智能跃迁，查询产品焕新**
    *   结合快递100AI大模型能力、亿级规模数据清洗与快递物流知识图谱技术沉淀，提供精准的快递预计到达时间与预计途径路线，实现快递查询从“快递到了哪里”向“快递何时可到达”升维。

4.   **一次配置，自动迭代**
    *   用户配置完成后无需反复操作，百递云开放平台会持续对服务进行更新迭代。

快递100 MCP Server的使用案例
---------------------

*   查询订单包裹的实时物流信息
*   根据快递运费、快递时效两大维度，对比不同快递公司物流方案，智能比选合适的快递服务
*   查询快递预计送达时间，结合物流轨迹及路线节点信息动态更新，预计到达时间越临近目的地越精准，方便用户提前规划收货安排

常见问题解答
------

**Q：使用快递100 MCP Server是否需要付费**

 A：用户需在[快递100API开放平台](https://api.kuaidi100.com/extend/register?code=d1660fe0390d4084b4f27b19d2feee02) 注册获取API Key，平台为每个用户提供单独的免费调用额度；如果后续超额，可在平台进行充值操作，如有疑问请联系：

 联系邮箱：api@kuaidi100.com

 联系电话：0755-86719032

快递100 MCP Server (Python)
-------------------------

通过`uv`安装`python`，最低版本要求为3.11

```
uv python install 3.11
```

### 一、在线获取依赖并使用（推荐）

通过`uvx`命令一步获取kuaidi100_mcp并使用

```
{
  "mcpServers": {
    "kuaidi100": {
      "command": "uvx",
      "args": [
        "kuaidi100-mcp"
      ],
      "env": {
        "KUAIDI100_API_KEY": "<YOUR_API_KEY>"
      }
    }
  }
}
```

### 二、下载至本地配置本地项目

通过`uv`创建一个项目

```
uv init mcp_server_kuaidi100
```

将`api_mcp.py`拷贝到该目录下，通过如下命令测试mcp server是否正常运行

```
uv run --with mcp[cli] mcp run {YOUR_PATH}/mcp_server_kuaidi100/api_mcp.py
# 如果是mac，需要加转义符
uv run --with mcp\[cli\] mcp run {YOUR_PATH}/mcp_server_kuaidi100/api_mcp.py
```

如果没有报错则MCP Server启动成功

### 获取快递100 API KEY

登录快递100获取 [https://api.kuaidi100.com/extend/register?code=d1660fe0390d4084b4f27b19d2feee02](https://api.kuaidi100.com/extend/register?code=d1660fe0390d4084b4f27b19d2feee02) （注意不要泄露授权key，以防被他人盗用！！！）

### 在支持MCP的客户端中使用

在MCP Server配置文件中添加如下内容后保存

```
{
  "mcpServers": {
    "kuaidi100": {
      "command": "uv",
      "args": [
        "run",
        "--with",
        "mcp[cli]",
        "mcp",
        "run",
        "{YOUR_PATH}/mcp_server_kuaidi100/api_mcp.py"
      ],
      "env": {
        "KUAIDI100_API_KEY": "<YOUR_API_KEY>"
      }
    }
  }
}
```

### 测试

#### 物流轨迹查询：

![Image 1: trae_test_queryTrace.png](https://file.kuaidi100.com/downloadfile/DTjS9PHPonJXikObm8OTcEA3OnuWBw0livDDJc73jYGMQmcwqfJpKhTzSVA-UwVX9LJZE3Nnnw7iLRgmekijRw)

#### 快递预估时效：

![Image 2: trae_test_estimateTime.png](https://file.kuaidi100.com/downloadfile/NL6vRCRVQkmvdavX19DISKf8uCvrj3q5NkSNl0ALv8GOOUufxrYRTRxoZJ20_uF-MGURmZRcKxS5XfAaz9t39Q)

#### 快递预估价格

![Image 3: trae_test_estimatePrice.png](https://file.kuaidi100.com/downloadfile/mPv7xFAUbsY5yFbaQZn7Z0ihtIU781pksXTTj-L2wwVgZ3dH-OSvqEdm3IaJzimTF_xIWbtHD6OFP8w2i35xsQ)

### Tips

如需获取账号信息（如 key、customer、secret），或免费试用100单，请访问[API开放平台](https://api.kuaidi100.com/home)进行注册