Surf MCP 服务器
---------------

为冲浪爱好者和网络冲浪者设计的MCP服务器。


视频演示
----------

[https://github.com/user-attachments/assets/0a4453e2-66df-4bf5-8366-8538cda366ed](https://github.com/user-attachments/assets/0a4453e2-66df-4bf5-8366-8538cda366ed)

特性
--------

*   使用纬度和经度获取任何地点的潮汐信息
*   支持特定日期的潮汐查询
*   详细的潮汐数据，包括高潮/低潮时间和站点信息
*   自动时区处理（UTC）

前提条件
-------------

*   Python 3.x
*   Storm Glass API 密钥

获取您的Storm Glass API密钥
--------------------------------

1.   访问 [Storm Glass](https://stormglass.io/)
2.   点击“免费试用”或“登录”创建账户
3.   注册后，您将收到API密钥

关于API使用限制的说明：

*   免费层级：每天10次请求
*   可选付费计划：
    *   小型：每天500次请求（€19/月）
    *   中型：每天5000次请求（€49/月）
    *   大型：每天25,000次请求（€129/月）
    *   企业级：可定制计划

根据您的使用需求选择合适的计划。免费层级适合测试和个人使用。

安装
------------

1.   克隆仓库：

```
git clone https://github.com/ravinahp/surf-mcp.git
cd surf-mcp
```

2.   使用uv安装依赖项：

```
uv sync
```

注意：我们使用 `uv` 而不是pip，因为项目使用 `pyproject.toml` 进行依赖管理。

配置为MCP服务器
-----------------------

要将此工具添加为MCP服务器，需要修改您的Claude桌面配置文件。该配置包括您的Storm Glass API密钥，因此无需单独设置。

配置文件的位置取决于您的操作系统：

*   MacOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
*   Windows: `%APPDATA%/Claude/claude_desktop_config.json`

向您的JSON文件中添加以下配置：

```
{
    "surf-mcp": {
        "command": "uv",
        "args": [
            "--directory",
            "/Users/YOUR_USERNAME/Code/surf-mcp",
            "run",
            "surf-mcp"
        ],
        "env": {
            "STORMGLASS_API_KEY": "your_api_key_here"
        }
    }
}
```

⚠️ 重要提示：

1.   将 `YOUR_USERNAME` 替换为您的实际系统用户名
2.   将 `your_api_key_here` 替换为您实际的Storm Glass API密钥
3.   确保目录路径与本地安装相匹配

部署
----------

### 构建

准备包的步骤：

1.   同步依赖并更新锁文件:

```
uv sync
```

2.   构建包：

```
uv build
```

这将在 `dist/` 目录下创建分发文件。

### 调试

由于MCP服务器通过标准输入输出运行，调试可能会比较困难。为了获得最佳调试体验，我们强烈建议使用MCP检查器。

您可以使用以下命令启动MCP检查器：

```
npx @modelcontextprotocol/inspector uv --directory /path/to/surf-mcp run surf-mcp
```

启动后，检查器会显示一个URL，您可以在浏览器中访问该URL开始调试。

检查器提供：

*   实时请求/响应监控
*   输入/输出验证
*   错误跟踪
*   性能指标

使用
-----

服务提供了一个FastMCP工具来获取潮汐信息：

```
@mcp.tool()
async def get_tides(latitude: float, longitude: float, date: str) -> str:
    """Get tide information for a specific location and date."""
```

### 参数：

*   `latitude`: 浮点数，表示地点的纬度
*   `longitude`: 浮点数，表示地点的纬度
*   `date`: 浮点数，表示地点的纬度

### 示例响应：

```
Tide Times:
Time: 2024-01-20T00:30:00+00:00 (UTC)
Type: HIGH tide
Height: 1.52m

Time: 2024-01-20T06:45:00+00:00 (UTC)
Type: LOW tide
Height: 0.25m

Station Information:
Name: Sample Station
Distance: 20.5km from requested location
```

用例
---------

### 示例 #1: 寻找最佳冲浪时间

您可以使用此工具来确定您最喜欢的海滩及最近站点的最佳冲浪时间。通常，最佳的冲浪条件是在涨潮期间，大约在高潮前2小时。


注意：不同的海滩可能基于其特定的地理环境和浪点类型有不同的最佳潮汐条件。该工具还提供了站点距离信息，这应该与潮汐信息一同考虑。（即，站点距离越长意味着不准确性越高 - 您也可以在提示时向 Claude 询问这一点）。

错误处理
--------------

该服务包括针对以下情况的强大错误处理：

*   API 请求失败
*   无效坐标
*   缺失或无效的 API 密钥
*   网络超时