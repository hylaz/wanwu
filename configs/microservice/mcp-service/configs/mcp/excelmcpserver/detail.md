Excel MCP 服务器
----------------

基于模型上下文协议（MCP）的Excel文件处理服务器，提供读取、写入和分析Excel文件的功能。

功能特点
--------

*   📖 读取Excel文件

    *   获取工作表列表
    *   读取特定工作表的数据
    *   读取所有工作表的数据

*   ✍️ 写入Excel文件

    *   创建新的Excel文件
    *   向特定工作表写入数据
    *   支持多工作表

*   🔍 分析Excel结构

    *   分析工作表结构
    *   将结构导出到新文件

*   💾 缓存管理

    *   自动缓存文件内容
    *   定时清理缓存
    *   手动清除缓存

*   📝 日志管理

    *   自动记录操作日志
    *   定期清理日志


安装
------------

### 通过Smithery安装

要通过[Smithery](https://smithery.ai/server/@zhiweixu/excel-mcp-server)自动为Claude Desktop安装excel-mcp-server：

```
npx -y @smithery/cli install @zhiweixu/excel-mcp-server --client claude
```

### 手动安装

通过NPM安装 可以通过在MCP服务器配置中添加以下配置来自动生成并安装excel-mcp-server。

Windows平台：

```
{
  "mcpServers": {
    "excel": {
        "command": "cmd",
        "args": ["/c", "npx", "--yes", "@zhiweixu/excel-mcp-server"],
        "env": {
            "LOG_PATH": "[set an accessible absolute path]"
        }
    }
}
```

其他平台：

```
{
  "mcpServers": {
    "excel": {
        "command": "npx",
        "args": ["--yes", "@zhiweixu/excel-mcp-server"],
        "env": {
            "LOG_PATH": "[set an accessible absolute path]"
        }
    }
}
```

注意：LOG_PATH是可选的。如果不设置，日志将存储在应用程序根目录下的'logs'文件夹中。

API工具
---------

### 结构工具
1.   analyzeExcelStructure
     *   功能：获取Excel文件结构，包括以JSON格式显示的工作表列表和列标题
     *   参数：
         *   fileAbsolutePath: Excel文件的绝对路径
         *   headerRows: 标题行数（默认值：1）

2.   exportExcelStructure
     *   功能：将Excel文件结构（工作表及标题）导出到一个新的Excel模板文件
     *   参数：
         *   sourceFilePath: 源Excel文件路径
         *   targetFilePath: 目标Excel文件路径
         *   headerRows: 标题行数（默认值：1）

### 读取工具
1.   readSheetNames
     *   功能：从Excel文件中获取所有工作表名称
     *   参数：
         *   fileAbsolutePath: Excel文件的绝对路径

2.   readDataBySheetName
     *   功能：从Excel文件中的特定工作表获取数据 
     *   参数：
         *   fileAbsolutePath: Excel文件的绝对路径
         *   sheetName: 要读取的工作表名称
         *   headerRow: 标题行号（默认值：1）
         *   dataStartRow: 数据起始行号（默认值：2）

3.   readSheetData
     *   功能：从Excel文件的所有工作表获取数据
     *   参数：
         *   fileAbsolutePath: Excel文件的绝对路径
         *   headerRow: 标题行号（默认值：1）
         *   dataStartRow: 数据起始行号（默认值：2）

### 写入工具
1.   writeDataBySheetName
     *   功能：将数据写入 Excel 文件中的特定工作表（如果工作表存在则覆盖）
     *   参数：
         *   fileAbsolutePath: Excel 文件的绝对路径
         *   sheetName: 要写入的工作表名称
         *   data: 要写入的数据数组

2.   writeSheetData
      *   功能：使用提供的数据创建一个新的 Excel 文件
      *   参数：
          *   fileAbsolutePath: 新 Excel 文件的绝对路径
          *   data: 包含多个工作表数据的对象

### 缓存工具
1.   clearFileCache
      *   功能：清除指定 Excel 文件的缓存数据
      *   参数：
          *   fileAbsolutePath: 需要从缓存中清除的 Excel 文件的绝对路径


配置
-------------

*   缓存配置
    *   缓存过期时间：1小时
    *   缓存清理间隔：4小时


*   日志配置
    *   日志保留天数：7天
    *   清理间隔：24小时


依赖
------------

*   @modelcontextprotocol/sdk: ^1.7.0
*   xlsx: ^0.18.5
*   typescript: ^5.8.2


开发依赖
------------------------

*   @types/node: ^22.13.10
*   nodemon: ^3.1.9
*   ts-node: ^10.9.2


许可证
------------------------

本项目采用 MIT 许可证。这意味着您可以：

*   将软件用于商业或非商业目的

*   修改源代码

*   分发原始或修改后的代码 要求：

*   保留原始版权声明

*   不得因软件使用向作者提出任何责任主张 有关详细的许可证信息，请参阅 LICENSE 文件。